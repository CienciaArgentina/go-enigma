package register

import (
	"fmt"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/encryption"

	"regexp"
	"strings"
	"time"
)

// The following errors will be returned to the user
var ()

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

type Service interface {
	SignUp(u *UserSignUp) (int64, []error)
}

type registerService struct {
	repository      Repository
	registerOptions *RegisterOptions
	config          *config.Configuration
}

func NewService(r Repository, ro *RegisterOptions, cfg *config.Configuration) Service {
	if ro == (&RegisterOptions{}) || ro == nil {
		ro = defaultRegisterOptions()
	}
	return &registerService{repository: r, registerOptions: ro, config: cfg}
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

func (rs *registerService) SignUp(u *UserSignUp) (int64, []error) {
	// User sign up form verifications
	var errs []error
	if ok, errs := rs.userSignUpDtoCanRegister(u); !ok {
		return 0, errs
	}

	user := &User{
		Username:           u.Username,
		NormalizedUsername: strings.ToUpper(u.Username),
		DateCreated:        time.Now(),
		VerificationToken:  encryption.GenerateVerificationToken(u.Email, rs.registerOptions.UserOptions.EmailVerificationExpiryDuration, rs.config),
		SecurityToken:      encryption.GenerateSecurityToken(u.Password, rs.config),
	}

	var err error
	user.PasswordHash, err = encryption.GenerateEncodedHash(u.Password)
	if err != nil {
		errs = append(errs, err)
	}

	userId, err := rs.repository.AddUser(user)
	if err != nil {
		errs = append(errs, err)
		return 0, errs
	}

	email := &UserEmail{
		UserId:          userId,
		Email:           u.Email,
		NormalizedEmail: strings.ToUpper(u.Email),
		VerfiedEmail:    false,
		DateCreated:     time.Now(),
	}

	_, err = rs.repository.AddEmail(email)
	if err != nil {
		errs = append(errs, err)
		return 0, errs
	}

	return userId, nil
}

func (rs *registerService) userSignUpDtoCanRegister(u *UserSignUp) (bool, []error) {
	var errs []error
	// Check that every field is correct
	if u.Username == "" {
		return false, append(errs, config.ErrEmptyUsername)
	}

	if u.Password == "" {
		return false, append(errs, config.ErrEmptyPassword)
	}

	if u.Email == "" {
		return false, append(errs, config.ErrEmptyEmail)
	}

	validEmail, err := regexp.Match("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
		[]byte(u.Email))
	if err != nil {
		return false, append(errs, err)
	}
	if !validEmail {
		return false, append(errs, config.ErrInvalidEmail)
	}

	if rs.registerOptions.UserOptions.RequireUniqueEmail {
		exists, err := rs.repository.VerifyIfUserExists(u.Email, u.Username)
		if err != nil {
			return false, append(errs, err)
		}

		if exists {
			return false, append(errs, config.ErrEmailAlreadyRegistered)
		}
	}

	usernameMatch, _ := regexp.Match(rs.registerOptions.UserOptions.AllowedCharacters, []byte(u.Username))
	if usernameMatch {
		errs = append(errs, config.ErrUsernameCotainsIlegalChars)
	}

	if strings.Contains(u.Password, " ") {
		errs = append(errs, config.ErrPwContainsSpace)
	}

	// Password checks
	if len(u.Password) < rs.registerOptions.PasswordOptions.RequiredLength {
		errs = append(errs, fmt.Errorf("El campo de contraseña tiene menos de %d caracteres", rs.registerOptions.PasswordOptions.RequiredLength))
	}

	if rs.registerOptions.PasswordOptions.RequireUppercase {
		match, _ := regexp.Match(".*[A-Z].*", []byte(u.Password))
		if !match {
			errs = append(errs, config.ErrPwDoesNotContainsUppercase)
		}
	}

	if rs.registerOptions.PasswordOptions.RequireLowercase {
		match, _ := regexp.Match(".*[a-z].*", []byte(u.Password))
		if !match {
			errs = append(errs, config.ErrPwDoesNotContainsLowercase)
		}
	}

	// List of avalaible chars: ~!@#$%^&*()-+=?/<>|{}_:;.,
	if rs.registerOptions.PasswordOptions.RequireNonAlphanumeric {
		match, _ := regexp.Match(".*[~!@#$%^&*()-+=?/<>|{}_:;.,].*", []byte(u.Password))
		if !match {
			errs = append(errs, config.ErrPwDoesNotContainsNonAlphaChars)
		}
	}

	if rs.registerOptions.PasswordOptions.RequireDigit {
		match, _ := regexp.Match(".*\\d.*", []byte(u.Password))
		if !match {
			errs = append(errs, config.ErrPwDoesNotContainsADigit)
		}
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, nil
}
