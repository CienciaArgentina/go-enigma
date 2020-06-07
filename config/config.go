package config

import (
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

const (
	GoEnvironment = "GO_ENVIRONMENT"
	Scope         = "SCOPE"

	Production  = "production"
	Test        = "test"
	Development = "development"
)

var (
	Config *Configuration

	errNotEvenDefaultConfiguration = fmt.Errorf("No es posible generar configuración default ya que el scope es %s", os.Getenv(Scope))

	ErrInvalidHash         = errors.New("El hash no usa el encoding correcto")
	ErrIncompatibleVersion = errors.New("Versión de argon2 incompatible")

	ErrEmptyUserId                     = errors.New("El userId no puede estar vacío")
	ErrEmailAlreadyRegistered          = errors.New("Este email ya se encuentra registrado en nuestra base de datos")
	ErrUsernameAlreadyRegistered       = errors.New("Este nombre de usuario ya se encuentra registrado")


	ErrEmailSendServiceNotWorking      = errors.New("Por alguna razón el servicio de envío de emails falló")





	ErrEmptySearch                     = errors.New("La búsqueda no arrojó ningún resultado")
)

const (
	// Request
	ErrInvalidBody     = "El cuerpo del mensaje que intentás enviar no es válido"
	ErrInvalidBodyCode = "invalid_body"

	// Empty
	ErrEmptyField                      = "Hay algún campo vacío y no puede estarlo"
	ErrEmptyFieldCode = "empty_field"
	ErrEmptyUsername            = "El nombre de usuario no puede estar vacío"
	ErrEmptyPassword            = "La contraseña no puede estar vacía"
	ErrEmptyEmail               = "El email no puede estar vacío"
	ErrEmptyEmailCode = "empty_email"
	ErrEmptyFieldUserCodeSignup = "invalid_user_signup"
	ErrEmptyFieldUserCodeLogin  = "invalid_user_login"

	// General
	ErrUnexpectedError                 = "Ocurrió un error en el sistema, por favor, ponete en contacto con sistemas"
)

type Configuration struct {
	AppName string `yaml:"appname"`
	Database      `yaml:"database"`
	Server        `yaml:"server"`
	Keys          `yaml:"keys"`
	Microservices `yaml:"microservices"`
	ArgonParams   `yaml:"argonparams"`
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
	Scheme string `yaml:"scheme"`
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
	Memory      uint32 `yaml:"memory"`
	Iterations  uint32 `yaml:"iterations"`
	Parallelism uint8  `yaml:"parallelism"`
	SaltLength  uint32 `yaml:"salt_length"`
	KeyLength   uint32 `yaml:"key_length"`
}

type RegisterOptions struct {
	UserOptions struct {
		// Set the allowed characters in username - Use a regex
		AllowedCharacters string
		// Email should not be registered on the database
		RequireUniqueEmail bool
		// How long the email verification token lasts
		EmailVerificationExpiryDuration time.Duration
	}
	PasswordOptions struct {
		// Password minimun required length
		RequiredLength int
		// Is it needed to have any non alphanumeric character in the password? (!*/$%&...)
		RequireNonAlphanumeric bool
		// Is it needed at least one lowercase character?
		RequireLowercase bool
		// Is it needed at least one uppercase character?
		RequireUppercase bool
		// Is it needed at least one digit? (123456...)
		RequireDigit bool
		// How many unique chars do the password need?
		RequiredUniqueChars int
	}
}

type LoginOptions struct {
	LockoutOptions struct {
		LockoutTimeDuration time.Duration
		MaxFailedAttempts   int
	}
	SignInOptions struct {
		RequireConfirmedEmail bool
	}
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
			ArgonParams: ArgonParams{
				Memory:      65536,
				Iterations:  2,
				Parallelism: 1,
				SaltLength:  32,
				KeyLength:   32,
			},
		}
	}

	panic(errNotEvenDefaultConfiguration)

	//pwd, err := os.Getwd()
	//files, _ := ioutil.ReadDir(pwd)
	//var sb strings.Builder
	//for _, f := range files {
	//	sb.WriteString(fmt.Sprintf("- %s \n", f.Name()))
	//}
	//panic(fmt.Sprintf("LS:  %s \n | WD: %s", sb.String(), pwd))
}

func New() *Configuration {
	config := &Configuration{}
	scope := os.Getenv(Scope)
	if scope == "" {
		scope = Development
	}

	data, err := os.Open(fmt.Sprintf("./config/config.%s.yml", scope))
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
