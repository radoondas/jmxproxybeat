// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Jmxbeat JmxbeatConfig
}

type JmxbeatConfig struct {
	Period         string               `yaml:"period"`
	URLs           []string             `yaml:"urls"`
	Authentication AuthenticationConfig `yaml:"authentication"`
	Beans          []BeanConfig         `yaml:"beans"`
}

type AuthenticationConfig struct {
	Username string
	Password string
}

type BeanConfig struct {
	Name       string   `yaml:"name"`
	Attributes []string `yaml:"attributes"`
	Keys       []string `yaml:"keys"`
}
