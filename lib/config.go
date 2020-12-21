package lib

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

type ServerConfig struct {
	Name    string  `yaml:"name"`
	URL     string  `yaml:"url"`
	Regex   *string `yaml:"regex"`
	JSONKey *string `yaml:"json_key"`

	parsedRegex *regexp.Regexp
}

func (s *ServerConfig) FetchAmount() (int, error) {
	if s.Regex != nil {
		if s.parsedRegex == nil {
			s.parsedRegex = regexp.MustCompile(strings.TrimSpace(*s.Regex))
		}

		return FetchPageRegex(s.URL, s.parsedRegex)
	}

	if s.JSONKey != nil {
		m := map[string]interface{}{}

		if err := FetchPageJson(s.URL, &m); err != nil {
			return 0, errors.Wrap(err, "unable to fetch json page")
		}

		if v, ok := m[*s.JSONKey]; ok {
			if v, ok := v.(float64); ok {
				return int(v), nil
			}

			return 0, errors.Errorf("unable to get player number of field type %T (%v)", v, v)
		}
		return 0, errors.Errorf("unable to find json key %s", *s.JSONKey)
	}

	return 0, errors.Errorf("unsupported fetch? no regex or json key found")
}

func (s ServerConfig) GetName() string { return s.Name }

var _ PrivateServer = (*ServerConfig)(nil)

type Config struct {
	Servers []ServerConfig `yaml:"servers"`
	ElasticSearch ElasticSearch `yaml:"elastic_search"`
	Interval string
}

func NewConfig(file string) *Config {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(errors.Wrap(err, "unable to read yaml"))
	}

	ret := &Config{}

	if err := yaml.Unmarshal(b, ret); err != nil {
		panic(errors.Wrap(err, "unable to unmarshal yaml"))
	}

	return ret
}
