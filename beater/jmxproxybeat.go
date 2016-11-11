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
	hosts  []*url.URL
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

	bt.hosts = make([]*url.URL, len(bt.config.Hosts))
	for i := 0; i < len(bt.config.Hosts); i++ {
		h, err := url.Parse(bt.config.Hosts[i])
		if err != nil {
			logp.Err("Invalid JMX hosts: %v", err)
			return err
		}
		bt.hosts[i] = h
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
	for _, u := range bt.hosts {

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
