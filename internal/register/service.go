package register

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// The following errors will be returned to the user
var (
	errEmptyUser = errors.New("El usuario no puede estar vacío")
	errEmptyUsername = errors.New("El nombre de usuario no puede estar vacío")
	errEmptyEmail = errors.New("El email no puede estar vacío")
)

type RegisterOptions struct {
	UserOptions struct {
		// Email should not be registered on the database
		RequireUniqueEmail bool
		// How long the email verification token lasts
		EmailVerificationExpiryDuration time.Duration
	}
	PasswordOptions struct {
		// Password minimun required length
		RequiredLength         int
		// Is it needed to have any non alphanumeric character in the password? (!*/$%&...)
		RequireNonAlphanumeric bool
		// Is it needed at least one lowercase character?
		RequireLowercase       bool
		// Is it needed at least one uppercase character?
		RequireUppercase       bool
		// Is it needed at least one digit? (123456...)
		RequireDigit           bool
		// How many unique chars do the password need?
		RequiredUniqueChars     int
	}
}

// These are the standard default options
func defaultRegisterOptions() *RegisterOptions {
	o := &RegisterOptions{}

	o.UserOptions.RequireUniqueEmail = true
	o.UserOptions.EmailVerificationExpiryDuration, _ = time.ParseDuration("1d")

	o.PasswordOptions.RequiredLength = 8
	o.PasswordOptions.RequireLowercase = true
	o.PasswordOptions.RequireUppercase = true
	o.PasswordOptions.RequireDigit = true
	o.PasswordOptions.RequireNonAlphanumeric = true
	o.PasswordOptions.RequiredUniqueChars = 1

	return o
}

type RegisterRepository interface {
	SignUp(u *User) []error
	GenerateVerificationToken(email string) string
}

type Service interface {
	SignUp(u *User) []error
	GenerateVerificationToken(email string) string
}

type registerService struct {
	repository RegisterRepository
	config *RegisterOptions
}

func New(r RegisterRepository, ro *RegisterOptions) Service {
	if ro == (&RegisterOptions{}) || ro == nil {
		ro = defaultRegisterOptions()
	}
	return &registerService{repository: r, config: ro}
}

func (rs *registerService) SignUp(u *User) []error {
	var errs []error
	if u == (&User{}) {
		return append(errs, errEmptyUser)
	}

	if u.Username == "" {
		return append(errs, errEmptyUsername)
	}

	// This will normalize the username to uppercase just in case its needed
	u.NormalizedUsername = strings.ToUpper(u.Username)

	u.DateCreated = time.Now()

	return nil
}

func (rs *registerService) GenerateVerificationToken(email string) string {
	// TODO: Move the sign key to cfg file
	var sign = []byte("clave")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"expiryDate": time.Now().Add(rs.config.UserOptions.EmailVerificationExpiryDuration).Unix(),
		"timestamp": time.Now().Unix(),
	})

	tokenString, _ := token.SignedString(sign)

	return tokenString
}
