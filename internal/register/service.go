package register

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// The following errors will be returned to the user
var (
	errEmptyUsername                  = errors.New("El nombre de usuario no puede estar vacío")
	errEmptyEmail                     = errors.New("El email no puede estar vacío")
	errEmptyPassword                  = errors.New("El campo de contraseña no puede estar vacío")
	errPwDoesNotContainsUppercase     = errors.New("La contraseña debe contener al menos un caracter en mayúscula")
	errPwDoesNotContainsLowercase     = errors.New("La contraseña debe contener al menos un caracter en minúscula")
	errPwContainsSpace                = errors.New("La contraseña no puede poseer el caracter de espacio")
	errPwDoesNotContainsNonAlphaChars = errors.New("La contraseña debe poseer al menos 1 caracter (permitidos: ~!@#$%^&*()-+=?/<>|{}_:;.,)")
	errPwDoesNotContainsADigit 		  = errors.New("La contraseña debe poseer al menos 1 dígito")
	errUsernameCotainsIlegalChars	  = errors.New("El nombre de usuario posee caracteres no permitidos (Sólo letras, números y los caracteres `.` `-` `_`)")
)

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

// These are the standard default options
func defaultRegisterOptions() *RegisterOptions {
	o := &RegisterOptions{}

	o.UserOptions.RequireUniqueEmail = true
	o.UserOptions.AllowedCharacters = "[^a-zA-Z0-9\\s._\\-/]"
	o.UserOptions.EmailVerificationExpiryDuration, _ = time.ParseDuration("1d")

	o.PasswordOptions.RequiredLength = 8
	o.PasswordOptions.RequireLowercase = true
	o.PasswordOptions.RequireUppercase = true
	o.PasswordOptions.RequireDigit = true
	o.PasswordOptions.RequireNonAlphanumeric = true
	o.PasswordOptions.RequiredUniqueChars = 1

	return o
}

type Repository interface {
	SignUp(u *UserSignUpDto) []error
	GenerateVerificationToken(email string) string
}

type Service interface {
	SignUp(u *UserSignUpDto) []error
	GenerateVerificationToken(email string) string
}

type registerService struct {
	repository Repository
	config     *RegisterOptions
}

func New(r Repository, ro *RegisterOptions) Service {
	if ro == (&RegisterOptions{}) || ro == nil {
		ro = defaultRegisterOptions()
	}
	return &registerService{repository: r, config: ro}
}

func (rs *registerService) SignUp(u *UserSignUpDto) []error {
	// User sign up form verifications
	if ok, errs := rs.userSignUpDtoCanRegister(u); !ok {
		return errs
	}

	return nil
}

func (rs *registerService) GenerateVerificationToken(email string) string {
	// TODO: Move the sign key to cfg file
	var sign = []byte("clave")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      email,
		"expiryDate": time.Now().Add(rs.config.UserOptions.EmailVerificationExpiryDuration).Unix(),
		"timestamp":  time.Now().Unix(),
	})

	tokenString, _ := token.SignedString(sign)

	return tokenString
}

func (rs *registerService) userSignUpDtoCanRegister(u *UserSignUpDto) (bool, []error) {
	var errs []error
	// Check that every field is correct
	if u.Username == "" {
		return false, append(errs, errEmptyUsername)
	}

	if u.Password == "" {
		return false, append(errs, errEmptyPassword)
	}

	if u.Email == "" {
		return false, append(errs, errEmptyEmail)
	}

	if rs.config.UserOptions.RequireUniqueEmail {
		// TODO: Verify if email already exists
	}

	usernameMatch, _ := regexp.Match(rs.config.UserOptions.AllowedCharacters, []byte(u.Username))
	if usernameMatch {
		errs = append(errs, errUsernameCotainsIlegalChars)
	}

	if strings.Contains(u.Password, " ") {
		errs = append(errs, errPwContainsSpace)
	}

	// Password checks
	if len(u.Password) < rs.config.PasswordOptions.RequiredLength {
		errs = append(errs, fmt.Errorf("El campo de contraseña tiene menos de %s caracteres", rs.config.PasswordOptions.RequiredLength))
	}

	if rs.config.PasswordOptions.RequireUppercase {
		match, _ := regexp.Match(".*[A-Z].*", []byte(u.Password))
		if !match {
			errs = append(errs, errPwDoesNotContainsUppercase)
		}
	}

	if rs.config.PasswordOptions.RequireLowercase {
		match, _ := regexp.Match(".*[a-z].*", []byte(u.Password))
		if !match {
			errs = append(errs, errPwDoesNotContainsLowercase)
		}
	}

	// List of avalaible chars: ~!@#$%^&*()-+=?/<>|{}_:;.,
	if rs.config.PasswordOptions.RequireNonAlphanumeric {
		match, _ := regexp.Match(".*[~!@#$%^&*()-+=?/<>|{}_:;.,].*", []byte(u.Password))
		if !match {
			errs = append(errs, errPwDoesNotContainsNonAlphaChars)
		}
	}

	if rs.config.PasswordOptions.RequireDigit {
		match, _ := regexp.Match(".*\\d.*", []byte(u.Password))
		if !match {
			errs = append(errs, errPwDoesNotContainsADigit)
		}
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, nil
}
