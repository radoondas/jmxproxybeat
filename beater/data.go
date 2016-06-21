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
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

const (
	managerJmxproxy = "/manager/jmxproxy/"
)

func (bt *Jmxproxybeat) GetJMX(u url.URL) error {
	for i := 0; i < len(bt.Beans); i++ {
		for j := 0; j < len(bt.Beans[i].Attributes); j++ {
			if len(bt.Beans[i].Attributes[j].Keys) > 0 {
				for k := 0; k < len(bt.Beans[i].Attributes[j].Keys); k++ {

					err := bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j].Name, bt.Beans[i].Attributes[j].Keys[k], bt.CAFile)
					if err != nil {
						logp.Err("Error requesting JMX: %v", err)
					}
				}
			} else {
				if len(bt.Beans[i].Keys) > 0 {
					for k := 0; k < len(bt.Beans[i].Keys); k++ {

						err := bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j].Name, bt.Beans[i].Keys[k], bt.CAFile)
						if err != nil {
							logp.Err("Error requesting JMX: %v", err)
						}
					}

				} else {

					err := bt.GetJMXObject(u, bt.Beans[i].Name, bt.Beans[i].Attributes[j].Name, "", bt.CAFile)
					if err != nil {
						logp.Err("Error requesting JMX: %v", err)
					}
				}
			}
		}
	}
	return nil
}

func (bt *Jmxproxybeat) GetJMXObject(u url.URL, name, attribute, key string, CAFile string) error {

	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
    var ParsedUrl *url.URL

    if CAFile != "" {
		// Load our trusted certificate path
		pemData, err := ioutil.ReadFile(CAFile)
		if err != nil {
			panic(err)
		}
		ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
		if !ok {
		    logp.Err("Unable to load CA file")
			panic("Couldn't load PEM data")
		}
    }

	//client := &http.Client{}
	client := &http.Client{Transport: transport}

	ParsedUrl, err := url.Parse(u.String())
    if err != nil {
		logp.Err("Unable to parse URL String")
		panic(err)
    }

    ParsedUrl.Path += managerJmxproxy
    parameters := url.Values{}

	parameters.Add("get", name)

	//var jmxObject, 
    var jmxAttribute string
	if key != "" {
		//jmxObject = name + attributeURI + attribute + keyURI + key
		parameters.Add("att", attribute)
		parameters.Add("key", key)
		jmxAttribute = attribute + "." + key
	} else {
		//jmxObject = name + attributeURI + attribute
		parameters.Add("att", attribute)
		jmxAttribute = attribute
	}


	ParsedUrl.RawQuery = parameters.Encode()

	logp.Debug(selector, "Requesting JMX: %s", ParsedUrl.String())  

	req, err := http.NewRequest("GET", ParsedUrl.String(), nil)

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
