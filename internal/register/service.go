package register

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/config"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/CienciaArgentina/go-enigma/internal/encryption"
	"github.com/jmoiron/sqlx"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	errCannotDelete = errors.New("El usuario que se intenta borrar no existe o no se puede alcanzar")
)

const (
	// User - Sign up

	// General
	ErrCantCreateUser = "No es posible crear esta cuenta ya que hay errores en los campos"

	// Email regex
	ErrInvalidEmailFormat     = "El email no respeta el formato de email (ejemplo: ejemplo@dominio.com)"
	ErrInvalidEmailFormatCode = "invalid_email"

	// Email already exists
	ErrEmailAlreadyExists                = "La dirección de correo electrónica ya se encuentra registrada"
	ErrEmailAlreadyExistsCode            = "duplicate_email"
	ErrEmailAlreadyExistsInternalErr     = "Ocurrió un error al intentar validar si el email existe"
	ErrEmailAlreadyExistsInternalErrCode = "internal_error"

	// Username already exists
	ErrUserAlreadyExists                = "Este nombre de usuario ya se encuentra registrado"
	ErrUserAlreadyExistsCode            = "duplicate_user"
	ErrUserAlreadyExistsInternalErr     = "Ocurrió un error al intentar validar si el usuario existe"
	ErrUserAlreadyExistsInternalErrCode = "internal_error"

	// Username characters
	ErrInvalidUsernameCode        = "invalid_username"
	ErrUsernameCotainsIlegalChars = "El nombre de usuario posee caracteres no permitidos (Sólo letras, números y los caracteres `.` `-` `_`)"

	// Password
	ErrInvalidPasswordCode = "invalid_password"
	ErrPwContainsSpace     = "La contraseña no puede poseer espacios"

	// Password characters
	ErrPwDoesNotContainsUppercase     = "La contraseña debe contener al menos un caracter en mayúscula"
	ErrPwDoesNotContainsLowercase     = "La contraseña debe contener al menos un caracter en minúscula"
	ErrPwDoesNotContainsNonAlphaChars = "La contraseña debe poseer al menos 1 caracter (permitidos: ~!@#$%^&*()-+=?/<>|{}_:;.,)"
	ErrPwDoesNotContainsADigit        = "La contraseña debe poseer al menos 1 dígito"

	// Password hash error
	ErrPasswordHash     = "Se generó un problema al encriptar la contraseña"
	ErrPasswordHashCode = "password_hash_failed"

	// Add register
	ErrInvalidRegisterCode = "invalid_register"
	ErrAddingUser          = "Ocurrió un error al intentar agregar el usuario"

	// Add register email in register
	ErrAddingUserEmail = "Ocurrió un error al intentar agregar el email del usuario"
)

type registerService struct {
	cfg             *config.Configuration
	db              *sqlx.DB
	registerOptions *config.RegisterOptions
	repository      RegisterRepository
}

func NewService(c *config.Configuration, db *sqlx.DB, ro *config.RegisterOptions, r RegisterRepository) RegisterService {
	if ro == nil {
		ro = defaultRegisterOptions()
	}
	return &registerService{
		cfg:             c,
		db:              db,
		registerOptions: ro,
		repository:      r,
	}
}

// These are the default default options
func defaultRegisterOptions() *config.RegisterOptions {
	o := &config.RegisterOptions{}

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

func (u *registerService) CreateUser(usr *domain.UserSignupDTO) (int64, apierror.ApiError) {
	if cansignup, err := u.UserCanSignUp(usr); !cansignup {
		return 0, err
	}

	tx := u.db.MustBegin()

	user := &domain.User{
		Username:           usr.Username,
		NormalizedUsername: strings.ToUpper(usr.Username),
		VerificationToken:  encryption.GenerateVerificationToken(usr.Email, u.registerOptions.UserOptions.EmailVerificationExpiryDuration, u.cfg),
	}

	user.SecurityToken.String = encryption.GenerateSecurityToken(usr.Password, u.cfg)

	var err error
	user.PasswordHash, err = encryption.GenerateEncodedHash(usr.Password, u.cfg)
	if err != nil {
		return 0, apierror.New(http.StatusInternalServerError, ErrPasswordHash, apierror.NewErrorCause(ErrPasswordHash, ErrPasswordHashCode))
	}

	userId, err := u.repository.AddUser(tx, user)
	if err != nil {
		return 0, apierror.New(http.StatusInternalServerError, ErrAddingUser, apierror.NewErrorCause(err.Error(), ErrInvalidRegisterCode))
	}

	email := &domain.UserEmail{
		UserId:          userId,
		Email:           usr.Email,
		NormalizedEmail: strings.ToUpper(usr.Email),
		VerfiedEmail:    false,
	}

	_, err = u.repository.AddUserEmail(tx, email)
	if err != nil {
		tx.Rollback()
		return 0, apierror.New(http.StatusInternalServerError, ErrAddingUserEmail, apierror.NewErrorCause(err.Error(), ErrInvalidRegisterCode))
	}

	// TODO: Send verification email

	tx.Commit()
	return userId, nil
}

func (u *registerService) UserCanSignUp(usr *domain.UserSignupDTO) (bool, apierror.ApiError) {
	errs := apierror.NewWithStatus(http.StatusBadRequest).WithMessage(ErrCantCreateUser)

	// Check that every field is correct
	if usr.Username == "" {
		return false, apierror.New(http.StatusBadRequest, config.ErrEmptyUsername, apierror.NewErrorCause(config.ErrEmptyUsername, config.ErrEmptyFieldUserCodeSignup))
	}

	if usr.Password == "" {
		return false, apierror.New(http.StatusBadRequest, config.ErrEmptyPassword, apierror.NewErrorCause(config.ErrEmptyPassword, config.ErrEmptyFieldUserCodeSignup))
	}

	if usr.Email == "" {
		return false, apierror.New(http.StatusBadRequest, config.ErrEmptyEmail, apierror.NewErrorCause(config.ErrEmptyEmail, config.ErrEmptyFieldUserCodeSignup))
	}

	validEmail, err := regexp.Match("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,"+
		"61}[a-zA-Z0-9])?)*$", []byte(usr.Email))
	if !validEmail || err != nil {
		return false, apierror.New(http.StatusBadRequest, ErrInvalidEmailFormat, apierror.NewErrorCause(ErrInvalidEmailFormat, ErrInvalidEmailFormatCode))
	}

	if u.registerOptions.UserOptions.RequireUniqueEmail {
		exists, err := u.repository.CheckEmailExists(usr.Email)
		if exists {
			return false, apierror.New(http.StatusBadRequest, ErrEmailAlreadyExists, apierror.NewErrorCause(ErrEmailAlreadyExists, ErrEmailAlreadyExistsCode))
		} else if err != nil && err != sql.ErrNoRows {
			return false, apierror.New(http.StatusInternalServerError, ErrEmailAlreadyExistsInternalErr, apierror.NewErrorCause(err.Error(), ErrEmailAlreadyExistsInternalErrCode))
		}
	}

	usrexists, err := u.repository.CheckUsernameExists(usr.Username)
	if usrexists {
		return false, apierror.New(http.StatusBadRequest, ErrUserAlreadyExists, apierror.NewErrorCause(ErrUserAlreadyExists, ErrUserAlreadyExistsCode))
	} else if err != nil && err != sql.ErrNoRows {
		return false, apierror.New(http.StatusInternalServerError, ErrUserAlreadyExistsInternalErr, apierror.NewErrorCause(err.Error(), ErrUserAlreadyExistsInternalErrCode))
	}

	usernameMatch, _ := regexp.Match(u.registerOptions.UserOptions.AllowedCharacters, []byte(usr.Username))
	if usernameMatch {
		errs.AddError(ErrUsernameCotainsIlegalChars, ErrInvalidUsernameCode)
	}

	if strings.Contains(usr.Password, " ") {
		errs.AddError(ErrPwContainsSpace, ErrInvalidPasswordCode)
	}

	if len(usr.Password) < u.registerOptions.PasswordOptions.RequiredLength {
		errs.AddError(fmt.Sprintf("El campo de contraseña tiene menos de %d caracteres", u.registerOptions.PasswordOptions.RequiredLength), ErrInvalidPasswordCode)
	}

	if u.registerOptions.PasswordOptions.RequireUppercase {
		match, _ := regexp.Match(".*[A-Z].*", []byte(usr.Password))
		if !match {
			errs.AddError(ErrPwDoesNotContainsUppercase, ErrInvalidPasswordCode)
		}
	}

	if u.registerOptions.PasswordOptions.RequireLowercase {
		match, _ := regexp.Match(".*[a-z].*", []byte(usr.Password))
		if !match {
			errs.AddError(ErrPwDoesNotContainsLowercase, ErrInvalidPasswordCode)
		}
	}

	// List of avalaible chars: ~!@#$%^&*()-+=?/<>|{}_:;.,
	if u.registerOptions.PasswordOptions.RequireNonAlphanumeric {
		match, _ := regexp.Match(".*[~!@#$%^&*()-+=?/<>|{}_:;.,].*", []byte(usr.Password))
		if !match {
			errs.AddError(ErrPwDoesNotContainsNonAlphaChars, ErrInvalidPasswordCode)
		}
	}

	if u.registerOptions.PasswordOptions.RequireDigit {
		match, _ := regexp.Match(".*\\d.*", []byte(usr.Password))
		if !match {
			errs.AddError(ErrPwDoesNotContainsADigit, ErrInvalidPasswordCode)
		}
	}

	if len(errs.ErrError) > 0 {
		return false, errs
	}

	return true, nil
}
