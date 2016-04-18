package beater

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

const MANAGER_JMXPROXY = "/manager/jmxproxy/?get="

func (bt *Jmxproxybeat) GetJMX(u url.URL) error {
	for i := 0; i < len(bt.Beans); i++ {
		if len(bt.Beans[i].Keys) > 0 {
			for j := 0; j < len(bt.Beans[i].Attributes); j++ {
				for k := 0; k < len(bt.Beans[i].Keys); k++ {
					//str := bt.Beans[i].Name + "&att=" + bt.Beans[i].Attributes[j] + "&key=" + bt.Beans[i].Keys[k]
					//logp.Info("Req: %s", str)
					bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j], bt.Beans[i].Keys[k])
				}
			}
		} else {
			for j := 0; j < len(bt.Beans[i].Attributes); j++ {
				//str := bt.Beans[i].Name + "&att=" + bt.Beans[i].Attributes[j]
				//logp.Info("Req: %s", str)
				bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j], "")
			}
		}

	}
	return nil
}

func (bt *Jmxproxybeat) GetJMXObject(u url.URL, name, attribute, key string) error {
	client := &http.Client{}
	var jmxObject, jmxAttribute string
	if key != "" {
		jmxObject = name + "&att=" + attribute + "&key=" + key
		jmxAttribute = attribute + "." + key
	} else {
		jmxObject = name + "&att=" + attribute
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
		return fmt.Errorf("HTTP%s", res.Status)
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Scan()

	//logp.Info("Response body: %v", scanner.Text())
	jmxValue, err := GetJMXValue(scanner.Text())
	//TODO: error handling

	event := common.MapStr{
		"@timestamp": common.Time(time.Now()),
		"type":       "jmx",
		"bean": common.MapStr{
			"Name":      name,
			"Attribute": jmxAttribute,
			"Value":     jmxValue,
			"hostname":   u.Host,
		},
	}
	bt.events.PublishEvent(event)
	logp.Info("Event: %+v", event)

	return nil
}

func GetJMXValue(responseBody string) (float64, error) {
	var re *regexp.Regexp
	var respValue float64

	//TODO: This requires lots of tuning!!
	re = regexp.MustCompile("= (\\d+)$")
	if matches := re.FindStringSubmatch(responseBody); matches != nil {
		respV, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0.0, err
		}
		respValue = respV
	}
	return respValue, nil
}
