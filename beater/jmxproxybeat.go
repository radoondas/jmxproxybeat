package beater

import (
	"fmt"
	"net/url"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/radoondas/jmxproxybeat/config"
)

const (
	selector = "jmxproxybeat"
)

type Jmxproxybeat struct {
	config config.Config
	done   chan struct{}
	client beat.Client
	hosts  []*url.URL
	auth   bool
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	bt := &Jmxproxybeat{
		done:   make(chan struct{}),
		config: c,
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
			logp.NewLogger(selector).Error("Invalid JMX hosts: %v", err)
			return err
		}
		bt.hosts[i] = h
	}

	if bt.config.SSL.CAfile == "" {
		logp.NewLogger(selector).Info("CAFile IS NOT set.")
	}

	//Disable authentication when no username or password is set
	if bt.config.Authentication.Username == "" || bt.config.Authentication.Password == "" {
		logp.NewLogger(selector).Info("One of username or password IS NOT set.")
		bt.auth = false
	}

	return nil
}

func (bt *Jmxproxybeat) Run(b *beat.Beat) error {
	logp.NewLogger(selector).Info("Jmxproxybeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

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
					logp.NewLogger(selector).Error("Error while getting JMX: %v", err)
				}
			}
		GotoFinish:
		}(u)
	}

	<-bt.done
	return nil
}

func (bt *Jmxproxybeat) Stop() {
	logp.NewLogger(selector).Info("Stopping Jmxproxybeat")
	if bt.done != nil {
		bt.client.Close()
		close(bt.done)
	}
}
