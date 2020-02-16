package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	GoEnvironment = "GO_ENVIRONMENT"
	Scope         = "SCOPE"

	Production  = "prod"
	Test        = "test"
	Development = "dev"
)

var (
	errNotEvenDefaultConfiguration = fmt.Errorf("No es posible generar configuración default ya que el scope es %s", os.Getenv(Scope))
)

type Configuration struct {
	Database `yaml:database`
	Server   `yaml:server`
}

type Database struct {
	Username string `yaml:db_username`
	Password string `env:db_password`
	Hostname string `env:db_hostname`
	Port     string `yaml:db_port`
}

type Server struct {
	Port string `yaml:server_port`
}

func DefaultConfiguration() *Configuration {
	// Even though it's kind of difficult to get to this point, I made this function so I'm sure that I'm always connected to a development scope
	if os.Getenv(GoEnvironment) != Production && os.Getenv(Scope) != Production {
		return &Configuration{
			Database: Database{
				Username: "cienciaArgentinaDev",
				Password: "cienciaArgentina",
				Hostname: "localhost",
				Port:     "3306",
			},
			Server: Server{
				Port: ":8080",
			},
		}
	}

	panic(errNotEvenDefaultConfiguration)
}

func New() *Configuration {
	config := &Configuration{}
	scope := os.Getenv(Scope)
	if scope == "" {
		scope = Development
	}

	a, _ := os.Getwd()
	fmt.Print(a)
	data, err := os.Open(fmt.Sprintf("config.%s.yml", scope))
	if err != nil {
		return DefaultConfiguration()
	}

	defer data.Close()

	decoder := yaml.NewDecoder(data)
	if err := decoder.Decode(config); err != nil {
		panic(err)
	}

	err = envconfig.Process("db", config)
	if err != nil {
		return DefaultConfiguration()
	}

	return config
}

func ConsolePrintMessageByCienciaArgentina(msg string) {
	statusColor := "\033[30;45m"
	resetColor := "\033[0m"
	fmt.Print(fmt.Sprintf("%s %v %s - %s", statusColor, "Ciencia Argentina", resetColor, msg))
}
