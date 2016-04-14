package beater

import (
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/radoondas/jmxbeat/config"
)

const selector = "jmxbeat"

type Jmxbeat struct {
	beatConfig *config.Config
	done       chan struct{}
	period     time.Duration
	urls       []*url.URL
	auth       bool
	username   string
	password   string
	Beans      []Bean
	events     publisher.Client
}

type Bean struct {
	Name       string
	Attributes []string
	Keys       []string
}

// Creates beater
func New() *Jmxbeat {
	return &Jmxbeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Jmxbeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	return nil
}

func (bt *Jmxbeat) Setup(b *beat.Beat) error {

	bt.events = b.Events

	// Setting default period if not set
	if bt.beatConfig.Jmxbeat.Period == "" {
		bt.beatConfig.Jmxbeat.Period = "1s"
	}

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Jmxbeat.Period)
	if err != nil {
		return err
	}

	//define default URL if none provided
	var urlConfig []string
	if bt.beatConfig.Jmxbeat.URLs != nil {
		urlConfig = bt.beatConfig.Jmxbeat.URLs
	} else {
		urlConfig = []string{"http://127.0.0.1:8888"}
	}

	bt.urls = make([]*url.URL, len(urlConfig))
	for i := 0; i < len(urlConfig); i++ {
		u, err := url.Parse(urlConfig[i])
		if err != nil {
			logp.Err("Invalid JMX url: %v", err)
			return err
		}
		bt.urls[i] = u
	}

	if bt.beatConfig.Jmxbeat.Authentication.Username == "" || bt.beatConfig.Jmxbeat.Authentication.Password == "" {
		logp.Err("Username or password IS NOT set.")
		bt.auth = false
	} else {
		bt.username = bt.beatConfig.Jmxbeat.Authentication.Username
		bt.password = bt.beatConfig.Jmxbeat.Authentication.Password
		bt.auth = true
		logp.Info("Username and password IS set.")
	}

	bt.Beans = make([]Bean, len(bt.beatConfig.Jmxbeat.Beans))
	if bt.beatConfig.Jmxbeat.Beans == nil {
		logp.Err("No beans are configured set.")
		//TODO: default values (HeapMemory)?
	} else {
		for i := 0; i < len(bt.beatConfig.Jmxbeat.Beans); i++ {
			bt.Beans[i].Name = bt.beatConfig.Jmxbeat.Beans[i].Name
			bt.Beans[i].Attributes = bt.beatConfig.Jmxbeat.Beans[i].Attributes
			bt.Beans[i].Keys = bt.beatConfig.Jmxbeat.Beans[i].Keys

			logp.Debug(selector, "Bean name: %s", bt.beatConfig.Jmxbeat.Beans[i].Name)
			for j := 0; j < len(bt.beatConfig.Jmxbeat.Beans[i].Attributes); j++ {
				logp.Debug(selector, "\tBean attribute: %s", bt.beatConfig.Jmxbeat.Beans[i].Attributes[j])
			}
			for k := 0; k < len(bt.beatConfig.Jmxbeat.Beans[i].Keys); k++ {
				logp.Debug(selector, "\t\tAttribute key: %s", bt.beatConfig.Jmxbeat.Beans[i].Keys[k])
			}
		}
	}

	return nil
}

func (bt *Jmxbeat) Run(b *beat.Beat) error {
	logp.Info("jmxbeat is running! Hit CTRL-C to stop it.")

	//for each url
	for _, u := range bt.urls {
		go func(u *url.URL) {
			ticker := time.NewTicker(bt.period)
			defer ticker.Stop()

			for {
				select {
				case <-bt.done:
					goto GotoFinish
				case <-ticker.C:
				}

				err := bt.GetJMX(*u)
				if err != nil {
					logp.Err("Error while getttig JMX: %v", err)
				}
			}
		GotoFinish:
		}(u)
	}

	<-bt.done
	return nil
}

func (bt *Jmxbeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Jmxbeat) Stop() {
	close(bt.done)
}
