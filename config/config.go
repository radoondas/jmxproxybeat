// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	SSL            SSL            `config:"ssl"`
	URLs           []string       `config:"urls"`
	Authentication Authentication `config:"authentication"`
	Beans          []Bean         `config:"beans"`
	Period         time.Duration  `config:"period"`
}

type SSL struct {
	CAfile string
}

type Authentication struct {
	Username string
	Password string
}

type Bean struct {
	Name       string      `config:"name"`
	Attributes []Attribute `config:"attributes"`
	Keys       []string    `config:"keys"`
}

type Attribute struct {
	Name string   `config:"name"`
	Keys []string `config:"keys"`
}

var (
	DefaultConfig = Config{
		Period: 10 * time.Second,
		URLs:   []string{"http://127.0.0.1:8080"},
		Authentication: Authentication{
			Username: "",
			Password: "",
		},
		SSL: SSL{
			CAfile: "",
		},
		Beans: []Bean{
			{
				Name: "java.lang:type=Memory",
				Keys: []string{"committed", "init", "max", "used"},
				Attributes: []Attribute{
					{
						Name: "HeapMemoryUsage",
						Keys: []string{},
					}, {
						Name: "NonHeapMemoryUsage",
						Keys: []string{},
					},
				},
			},
		},
	}
)
