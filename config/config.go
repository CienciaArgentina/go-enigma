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
	Keys     `yaml:keys`
}

type Database struct {
	Username string `yaml:"db_username"`
	Password string `env:"db_password"`
	Hostname string `env:"db_hostname"`
	Port     string `yaml:"db_port"`
	Database string `env:"db_name"`
}

type Server struct {
	Port string `yaml:"server_port"`
}

type Keys struct {
	PasswordHashingKey string `env:"key_passwordHashing"`
}

type ArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
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
				Database: "cienciaargentinaauthdev",
			},
			Server: Server{
				Port: ":8080",
			},
			Keys: Keys{
				// This is just for a development scope
				PasswordHashingKey: "98616F779CAA278695ADAF88BF4C1",
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
