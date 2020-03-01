package register

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/CienciaArgentina/go-enigma/config"
	"golang.org/x/crypto/argon2"
	"regexp"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	GenerateVerificationToken(email string) string
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
		VerificationToken:  rs.GenerateVerificationToken(u.Email),
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      email,
		"expiryDate": time.Now().Add(rs.registerOptions.UserOptions.EmailVerificationExpiryDuration).Unix(),
		"timestamp":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(rs.config.Keys.PasswordHashingKey))
	if err != nil {
		// TODO: log this
	}
	return tokenString
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
		exists, err := rs.repository.VerifyIfEmailExists(u.Email)
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

func generateFromPassword(password string, p *config.ArgonParams) (string, error) {
	// Generate a cryptographically secure random salt.
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

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

// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// For guidance and an outline process for choosing appropriate parameters see https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04#section-4.
func GenerateEncodedHash(pw string) (string, error) {
	p := &config.ArgonParams{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: 4,
		SaltLength:  32,
		KeyLength:   32,
	}

	encodedHash, err := generateFromPassword(pw, p)
	if err != nil {
		return "", err
	}

	return encodedHash, nil
}
