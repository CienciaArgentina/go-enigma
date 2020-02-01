package conf

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	SCOPE = "SCOPE"
	GOENVIRONMENT = "GO_ENVIRONMENT"
	Production = "production"
	Testing = "testing"
	Local = "local"

)

type Configuration struct {
	Database struct {
		User string `yaml:user`
		Password string `envconfig:DB_PASSWORD`
		Host string `envconfig:DB_HOST`
		Database string `yaml:dbname`
	} `yaml:database`
}

func New() (*Configuration, error) {
	config := &Configuration{}
	err := envconfig.Process("db", &config)
	if err != nil {
		return nil, err
	}

	scope := os.Getenv(SCOPE)
	if scope == "" {
		scope = Local
	}

	data, err := os.Open(fmt.Sprintf("config.%s.yml", scope))
	if err != nil {
		return nil, err
	}

	defer data.Close()

	decoder := yaml.NewDecoder(data)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
