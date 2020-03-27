package config

import (
	"errors"
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
	Config *Configuration

	errNotEvenDefaultConfiguration     = fmt.Errorf("No es posible generar configuración default ya que el scope es %s", os.Getenv(Scope))
	ErrInvalidLogin                    = errors.New("El usuario o la contraseña especificados no existe")
	ErrInvalidHash                     = errors.New("El hash no usa el encoding correcto")
	ErrIncompatibleVersion             = errors.New("Versión de argon2 incompatible")
	ErrThroughLogin                    = errors.New("Ocurrió un error al momento de loguear")
	ErrEmailNotVerified                = errors.New("Tu dirección de email no fue verificada aún")
	ErrEmptyUsername                   = errors.New("El nombre de usuario no puede estar vacío")
	ErrEmptyUserId                     = errors.New("El userId no puede estar vacío")
	ErrEmptyEmail                      = errors.New("El email no puede estar vacío")
	ErrEmptyPassword                   = errors.New("El campo de contraseña no puede estar vacío")
	ErrPwDoesNotContainsUppercase      = errors.New("La contraseña debe contener al menos un caracter en mayúscula")
	ErrPwDoesNotContainsLowercase      = errors.New("La contraseña debe contener al menos un caracter en minúscula")
	ErrPwContainsSpace                 = errors.New("La contraseña no puede poseer el caracter de espacio")
	ErrPwDoesNotContainsNonAlphaChars  = errors.New("La contraseña debe poseer al menos 1 caracter (permitidos: ~!@#$%^&*()-+=?/<>|{}_:;.,)")
	ErrPwDoesNotContainsADigit         = errors.New("La contraseña debe poseer al menos 1 dígito")
	ErrUsernameCotainsIlegalChars      = errors.New("El nombre de usuario posee caracteres no permitidos (Sólo letras, números y los caracteres `.` `-` `_`)")
	ErrEmailAlreadyRegistered          = errors.New("Este email ya se encuentra registrado en nuestra base de datos")
	ErrUsernameAlreadyRegistered = errors.New("Este nombre de usuario ya se encuentra registrado")
	ErrInvalidEmail                    = errors.New("El email no respeta el formato de email (ejemplo: ejemplo@dominio.com)")
	ErrUnexpectedError                 = errors.New("Ocurrió un error en el sistema")
	ErrEmailAlreadyVerified            = errors.New("El mail ya se encuentra confirmado")
	ErrEmailSendServiceNotWorking      = errors.New("Por alguna razón el servicio de envío de emails falló")
	ErrEmailValidationFailed           = errors.New("La validación del email falló por algún campo vacío")
	ErrEmptyField                      = errors.New("Hay algún campo vacío y no puede estarlo")
	ErrValidationTokenFailed           = errors.New("La validación del token falló")
	ErrPasswordConfirmationDoesntMatch = errors.New("Los passwords ingresados no son idénticos")
	ErrPasswordTokenIsNotValid         = errors.New("El token para resetear la contraseña no es válido")
	ErrEmptySearch                     = errors.New("La búsqueda no arrojó ningún resultado")
)

type Configuration struct {
	Database      `yaml:"database"`
	Server        `yaml:"server"`
	Keys          `yaml:"keys"`
	Microservices `yaml:"microservices"`
}

type Database struct {
	Username string `yaml:"db_username"`
	Password string `envconfig:"ENV_DB_PASSWORD"`
	Hostname string `envconfig:"ENV_DB_HOSTNAME"`
	Port     string `yaml:"db_port"`
	Database string `envconfig:"ENV_DB_NAME"`
}

type Server struct {
	Port string `yaml:"server_port"`
}

type Keys struct {
	PasswordHashingKey string `env:"ENV_KEY_PASSWORDHASHING"`
}

type Microservices struct {
	BaseUrl        string `yaml:"base_url"`
	UsersEndpoints struct {
		BaseResource          string `yaml:"base_resource"`
		Login                 string `yaml:"login"`
		SignUp                string `yaml:"sign_up"`
		SendConfirmationEmail string `yaml:"send_confirmation_email"`
		ConfirmEmail          string `yaml:"confirm_email"`
		SendPasswordReset     string `yaml:"send_password_reset"`
	} `yaml:"user_endpoints"`
	EmailSenderAddr      string `yaml:"email_sender_addr"`
	EmailSenderEndpoints struct {
		SendEmail string `yaml:"email_sender_send_email"`
	} `yaml:"email_sender_endpoints"`
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

	data, err := os.Open(fmt.Sprintf("../../config/config.%s.yml", scope))
	if err != nil {
		return DefaultConfiguration()
	}

	defer data.Close()

	decoder := yaml.NewDecoder(data)
	if err := decoder.Decode(config); err != nil {
		panic(err)
	}

	err = envconfig.Process("env_", config)
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
