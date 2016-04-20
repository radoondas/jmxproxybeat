// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Jmxproxybeat JmxproxybeatConfig
}

type JmxproxybeatConfig struct {
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
	Attributes []Attribute `yaml:"attributes"`
	Keys       []string `yaml:"keys"`
}

type Attribute struct {
	Name string `yaml:"name"`
	Keys []string `yaml:"keys"`
}