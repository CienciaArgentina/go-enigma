package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/clog"
	"github.com/CienciaArgentina/go-backend-commons/pkg/scope"
)

const (
	envPasswordHashing  = "PASSWORD_HASHING_KEY"
	envArgonMemory      = "ARGON_MEMORY"
	envArgonIterations  = "ARGON_ITERATIONS"
	envArgonParallelism = "ARGON_PARALLELISM"
	envArgonSaltLength  = "ARGON_SALT_LENGTH"
	envArgonKeyLength   = "ARGON_KEY_LENGTH"

	defaultArgonMemory      = 65536
	defaultArgonIterations  = 2
	defaultArgonParallelism = 1
	defaultArgonSaltLength  = 32
	defaultArgonKeyLength   = 32
)

type EnigmaConfig struct {
	Keys            *Keys
	ArgonParams     *ArgonParams
	RegisterOptions *RegisterOptions
	LoginOptions    *LoginOptions
	Microservices
}

type Server struct {
	Port string `yaml:"server_port"`
}

type Keys struct {
	PasswordHashingKey string
}

type Microservices struct {
	Scheme         string `yaml:"scheme"`
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
	Parallelism uint8
	Memory      uint32
	Iterations  uint32
	SaltLength  uint32
	KeyLength   uint32
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

func NewEnigmaConfig() (*EnigmaConfig, error) {
	cfg := &EnigmaConfig{}

	var err error
	cfg.Keys = &Keys{}
	cfg.Keys.PasswordHashingKey, err = cfg.getPasswordHashingKey()
	if err != nil {
		clog.Panic(err.Error(), "get-password-hashing-key", err, nil)
		return nil, err
	}
	cfg.ArgonParams = &ArgonParams{}
	cfg.ArgonParams, err = cfg.getArgonParams()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (e *EnigmaConfig) getPasswordHashingKey() (string, error) {
	hash := os.Getenv(envPasswordHashing)
	if hash == "" {
		msg := "password hashing key is empty"
		err := errors.New(msg)
		return "", err
	}

	return hash, nil
}

func (e *EnigmaConfig) getArgonParams() (*ArgonParams, error) {
	params := &ArgonParams{}

	memory := os.Getenv(envArgonMemory)
	iterations := os.Getenv(envArgonIterations)
	parallelism := os.Getenv(envArgonParallelism)
	saltLen := os.Getenv(envArgonSaltLength)
	keyLen := os.Getenv(envArgonKeyLength)

	if !scope.IsProductiveScope() {
		params.Memory = defaultArgonMemory
		params.Iterations = defaultArgonIterations
		params.Parallelism = defaultArgonParallelism
		params.SaltLength = defaultArgonSaltLength
		params.KeyLength = defaultArgonKeyLength
	} else {
		if memory == "" || iterations == "" || parallelism == "" || saltLen == "" || keyLen == "" {
			msg := "argon params are empty"
			err := errors.New(msg)
			clog.Panic(msg, "get-argon-params", err, map[string]string{"memory": memory, "iterations": iterations, "parallelism": parallelism, "salt": saltLen, "key": keyLen})
			return nil, err
		}

		mem, err := strconv.ParseInt(memory, 10, 32)
		if err != nil {
			clog.Panic("Memory cannot be casted", "get-argon-params", err, nil)
			return nil, err
		}
		params.Memory = uint32(mem)

		it, err := strconv.ParseInt(iterations, 10, 32)
		if err != nil {
			clog.Panic("Iterations cannot be casted", "get-argon-params", err, nil)
			return nil, err
		}
		params.Iterations = uint32(it)

		para, err := strconv.ParseInt(parallelism, 10, 8)
		if err != nil {
			clog.Panic("Parallelism cannot be casted", "get-argon-params", err, nil)
			return nil, err
		}
		params.Parallelism = uint8(para)

		salt, err := strconv.ParseInt(saltLen, 10, 8)
		if err != nil {
			clog.Panic("Salt Length cannot be casted", "get-argon-params", err, nil)
			return nil, err
		}
		params.SaltLength = uint32(salt)

		key, err := strconv.ParseInt(keyLen, 10, 8)
		if err != nil {
			clog.Panic("Salt Length cannot be casted", "get-argon-params", err, nil)
			return nil, err
		}
		params.KeyLength = uint32(key)
	}

	return params, nil
}
