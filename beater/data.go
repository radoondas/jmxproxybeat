package beater

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

const (
	MANAGER_JMXPROXY = "/manager/jmxproxy/?get="
	ATTRIBUTE_URI = "&att="
	KEY_URI = "&key="
)

func (bt *Jmxproxybeat) GetJMX(u url.URL) error {
	for i := 0; i < len(bt.Beans); i++ {
		for j := 0; j < len(bt.Beans[i].Attributes); j++ {
			if len(bt.Beans[i].Attributes[j].Keys) > 0 {
				for k := 0; k < len(bt.Beans[i].Attributes[j].Keys); k++ {
					logp.Debug(selector, "Host: %s, request: %s", u.String(), bt.Beans[i].Name +
					ATTRIBUTE_URI + bt.Beans[i].Attributes[j].Name +
					KEY_URI + bt.Beans[i].Attributes[j].Keys[k])

					err := bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j].Name, bt.Beans[i].Attributes[j].Keys[k])
					if err != nil {
						logp.Err("Error requesting JMX for %s: %v", bt.Beans[i].Name +
						ATTRIBUTE_URI + bt.Beans[i].Attributes[j].Name +
						KEY_URI + bt.Beans[i].Attributes[j].Keys[k], err)
					}
				}
			} else {
				if len(bt.Beans[i].Keys) > 0 {
					for k := 0; k < len(bt.Beans[i].Keys); k++ {
						logp.Debug(selector, "Host: %s, request: %s", u.String(), bt.Beans[i].Name +
						ATTRIBUTE_URI + bt.Beans[i].Attributes[j].Name +
						KEY_URI + bt.Beans[i].Keys[k])

						err := bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j].Name, bt.Beans[i].Keys[k])
						if err != nil {
							logp.Err("Error requesting JMX for %s: %v", bt.Beans[i].Name +
							ATTRIBUTE_URI + bt.Beans[i].Attributes[j].Name +
							KEY_URI + bt.Beans[i].Keys[k], err)
						}
					}

				} else {
					logp.Debug(selector, "Host: %s, request: %s", u.String(), bt.Beans[i].Name +
					ATTRIBUTE_URI + bt.Beans[i].Attributes[j].Name)

					err := bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j].Name, "")
					if err != nil {
						logp.Err("Error requesting JMX for %s: %v", bt.Beans[i].Name +
						ATTRIBUTE_URI + bt.Beans[i].Attributes[j].Name, err)
					}
				}
			}
		}
	}
	return nil
}

func (bt *Jmxproxybeat) GetJMXObject(u url.URL, name, attribute, key string) error {
	client := &http.Client{}
	var jmxObject, jmxAttribute string
	if key != "" {
		jmxObject = name + ATTRIBUTE_URI + attribute + KEY_URI + key
		jmxAttribute = attribute + "." + key
	} else {
		jmxObject = name + ATTRIBUTE_URI + attribute
		jmxAttribute = attribute
	}

	req, err := http.NewRequest("GET", u.String()+MANAGER_JMXPROXY+jmxObject, nil)

	if bt.auth {
		req.SetBasicAuth(bt.username, bt.password)
	}
	res, err := client.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("HTTP %s", res.Status)
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Scan()

	jmxValue, err := GetJMXValue(scanner.Text())
	if err != nil {
		return err
	}

	event := common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "jmx",
		"bean": common.MapStr{
			"name":      name,
			"attribute": jmxAttribute,
			"value":     jmxValue,
			"hostname":  u.Host,
		},
	}
	bt.events.PublishEvent(event)
	logp.Info("Event: %+v", event)

	return nil
}

func GetJMXValue(responseBody string) (float64, error) {
	var re *regexp.Regexp
	var respValue float64

	if strings.HasPrefix(responseBody, "Error") {
		return 0, errors.New(responseBody)
	}

	//TODO: This requires lots of tuning!!
	re = regexp.MustCompile("\\d+(\\.\\d+)?$")
	if matches := re.FindStringSubmatch(responseBody); matches != nil {
		respV, err := strconv.ParseFloat(matches[0], 64)
		//TODO: test for empty string!
		if err != nil {
			return 0.0, err
		}
		respValue = respV
	}
	return respValue, nil
}
