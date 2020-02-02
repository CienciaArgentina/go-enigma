package conf

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	SCOPE         = "SCOPE"
	GOENVIRONMENT = "GO_ENVIRONMENT"
	Production    = "production"
	Testing       = "testing"
	Development   = "development"
)

type Configuration struct {
	Database struct {
		User     string `yaml:user`
		Password string `envconfig:"DB_PASSWORD"`
		Host     string `envconfig:"DB_HOST"`
		Database string `yaml:dbname`
	} `yaml:database`
	Server struct {
		Port int `yaml:port`
	} `yaml:server`
}

func New() *Configuration {
	config := &Configuration{}
	scope := os.Getenv(SCOPE)
	if scope == "" {
		scope = Development
	}

	a, _ := os.Getwd()
	fmt.Print(a)
	data, err := os.Open(fmt.Sprintf("conf/config.%s.yml", scope))
	if err != nil {
		// handle error
	}

	defer data.Close()

	decoder := yaml.NewDecoder(data)
	if err := decoder.Decode(config); err != nil {
		// handle error
	}

	err = envconfig.Process("DB", config)
	if err != nil {
		// handle error
	}

	return config
}

func ConsolePrintMessageByCienciaArgentina(msg string) {
	statusColor := "\033[30;45m"
	resetColor := "\033[0m"
	fmt.Print(fmt.Sprintf("%s %v %s - %s", statusColor, "Ciencia Argentina", resetColor, msg))
}
