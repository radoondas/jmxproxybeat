package beater

import (
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	cfg "github.com/radoondas/jmxproxybeat/config"
)

const (
	selector = "jmxproxybeat"
)

type Jmxproxybeat struct {
	config cfg.Config
	done   chan struct{}
	urls   []*url.URL
	auth   bool
	client publisher.Client
}

// Creates beater
func New(b *beat.Beat, rawCfg *common.Config) (beat.Beater, error) {
	config := cfg.DefaultConfig
	if err := rawCfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Jmxproxybeat{
		done:   make(chan struct{}),
		config: config,
		auth:   true,
	}

	err := bt.init(b)
	if err != nil {
		return nil, err
	}

	return bt, nil
}

/// *** Beater interface methods ***///
func (bt *Jmxproxybeat) init(b *beat.Beat) error {

	bt.urls = make([]*url.URL, len(bt.config.URLs))
	for i := 0; i < len(bt.config.URLs); i++ {
		u, err := url.Parse(bt.config.URLs[i])
		if err != nil {
			logp.Err("Invalid JMX url: %v", err)
			return err
		}
		bt.urls[i] = u
	}

	if bt.config.SSL.CAfile == "" {
		logp.Info("CAFile IS NOT set.")
	}

	//Disable authentication when no username or password is set
	if bt.config.Authentication.Username == "" || bt.config.Authentication.Password == "" {
		logp.Info("One of username or password IS NOT set.")
		bt.auth = false
	}

	return nil
}

func (bt *Jmxproxybeat) Run(b *beat.Beat) error {
	logp.Info("Jmxproxybeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()
	//for each url
	for _, u := range bt.urls {

		go func(u *url.URL) {
			ticker := time.NewTicker(bt.config.Period)
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

func (bt *Jmxproxybeat) Stop() {
	logp.Info("Stopping Jmxproxybeat")
	if bt.done != nil {
		bt.client.Close()
		close(bt.done)
	}
}

// OLD snippet of code
//for i := 0; i < len(bt.config.Beans); i++ {
//		bt.Beans[i].Name = bt.config.Beans[i].Name
//		bt.Beans[i].Attributes = bt.config.Beans[i].Attributes
//		bt.Beans[i].Keys = bt.config.Beans[i].Keys
//
//		logp.Debug(selector, "Bean name: %s", bt.config.Beans[i].Name)
//		for j := 0; j < len(bt.config.Beans[i].Attributes); j++ {
//			logp.Debug(selector, "\tBean attribute: %s", bt.config.Beans[i].Attributes[j].Name)
//			if len(bt.config.Beans[i].Attributes[j].Keys) > 0 {
//				for k := 0; k < len(bt.config.Beans[i].Attributes[j].Keys); k++ {
//					logp.Debug(selector, "\t\tAttribute key: %s", bt.config.Beans[i].Attributes[j].Keys[k])
//				}
//			}
//		}
//		for k := 0; k < len(bt.config.Beans[i].Keys); k++ {
//			logp.Debug(selector, "\tBean key: %s", bt.config.Beans[i].Keys[k])
//		}
//	}
