package register

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
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
	errPwDoesNotContainsADigit        = errors.New("La contraseña debe poseer al menos 1 dígito")
	errUsernameCotainsIlegalChars     = errors.New("El nombre de usuario posee caracteres no permitidos (Sólo letras, números y los caracteres `.` `-` `_`)")
	errEmailAlreadyRegistered         = errors.New("Este email ya se encuentra registrado en nuestra base de datos")
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

type Repository interface {
	AddUser(u *User) (int64, error)
	AddEmail(e *UserEmail) (int64, error)
	VerifyIfEmailExists(email string) (bool, error)
}

type Service interface {
	SignUp(u *UserSignUpDto) (int64, []error)
	GenerateVerificationToken(email string) string
}

type registerService struct {
	repository Repository
	config     *RegisterOptions
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

func New(r Repository, ro *RegisterOptions) Service {
	if ro == (&RegisterOptions{}) || ro == nil {
		ro = defaultRegisterOptions()
	}
	return &registerService{repository: r, config: ro}
}

func (rs *registerService) SignUp(u *UserSignUpDto) (int64, []error) {
	// User sign up form verifications
	var errs []error
	if ok, errs := rs.userSignUpDtoCanRegister(u); !ok {
		return 0, errs
	}

	user := &User{
		Username:            u.Username,
		NormalizedUsername:  strings.ToUpper(u.Username),
		LockoutEnabled:      false,
		FailedLoginAttempts: 0,
		DateCreated:         time.Now(),
		VerificationToken:   rs.GenerateVerificationToken(u.Email),
	}

	var err error
	user.PasswordHash, err = GenerateEncodedHash(u.Password)
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

	// TODO: Verify if email regex match

	if rs.config.UserOptions.RequireUniqueEmail {
		exists, err := rs.repository.VerifyIfEmailExists(u.Email)
		if err != nil {
			return false, append(errs, err)
		}

		if exists {
			return false, append(errs, errEmailAlreadyRegistered)
		}
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
		errs = append(errs, fmt.Errorf("El campo de contraseña tiene menos de %d caracteres", rs.config.PasswordOptions.RequiredLength))
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

func generateFromPassword(password string, p *argonParams) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type argonParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// For guidance and an outline process for choosing appropriate parameters see https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04#section-4.
func GenerateEncodedHash(pw string) (string, error) {
	p := &argonParams{
		memory:      128 * 1024,
		iterations:  4,
		parallelism: 4,
		saltLength:  32,
		keyLength:   32,
	}

	encodedHash, err := generateFromPassword(pw, p)
	if err != nil {
		return "", err
	}

	return encodedHash, nil
}
