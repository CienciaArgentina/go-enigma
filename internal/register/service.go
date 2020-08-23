package register

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/CienciaArgentina/go-enigma/internal/domain"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/encryption"
	"github.com/jmoiron/sqlx"
)

var errCannotDelete = errors.New("El usuario que se intenta borrar no existe o no se puede alcanzar")

const (
	// User - Sign up.

	// General.
	errCantCreateUser = "No es posible crear esta cuenta ya que hay errores en los campos"

	// Email regex.
	errInvalidEmailFormat = "El email no respeta el formato de email (ejemplo: ejemplo@dominio.com)"

	// Email already exists.
	errEmailAlreadyExists            = "La dirección de correo electrónica ya se encuentra registrada"
	errEmailAlreadyExistsInternalErr = "Ocurrió un error al intentar validar si el email existe"

	// Username already exists.
	errUserAlreadyExists            = "Este nombre de usuario ya se encuentra registrado"
	errUserAlreadyExistsInternalErr = "Ocurrió un error al intentar validar si el usuario existe"

	// Username characters.
	errInvalidUsernameCode        = "invalid_username"
	errUsernameCotainsIlegalChars = "El nombre de usuario posee caracteres no permitidos (Sólo letras, números y los caracteres `.` `-` `_`)"

	// Password.
	errInvalidPasswordCode = "invalid_password"
	errPwContainsSpace     = "La contraseña no puede poseer espacios"

	// Password characters.
	errPwDoesNotContainsUppercase     = "La contraseña debe contener al menos un caracter en mayúscula"
	errPwDoesNotContainsLowercase     = "La contraseña debe contener al menos un caracter en minúscula"
	errPwDoesNotContainsNonAlphaChars = "La contraseña debe poseer al menos 1 caracter (permitidos: ~!@#$%^&*()-+=?/<>|{}_:;.,)"
	errPwDoesNotContainsADigit        = "La contraseña debe poseer al menos 1 dígito"

	// Password hash error.
	errPasswordHash     = "Se generó un problema al encriptar la contraseña"
	errPasswordHashCode = "password_hash_failed"

	// Add register.
	errInvalidRegisterCode = "invalid_register"
	errAddingUser          = "Ocurrió un error al intentar agregar el usuario"

	// Add register email in register.
	errAddingUserEmail = "Ocurrió un error al intentar agregar el email del usuario"

	errGenerateVerificationToken = "Ocurrió un error al generar el token de verificación"
	errGenerateSecurityToken     = "Ocurrió un error al generar el security token"

	errTokenGeneration = "failed_token_generation"
)

type registerService struct {
	cfg             *config.EnigmaConfig
	db              *sqlx.DB
	registerOptions *config.RegisterOptions
	repository      RegisterRepository
}

func NewService(c *config.EnigmaConfig, db *sqlx.DB, r RegisterRepository) RegisterService {
	return &registerService{
		cfg:             c,
		db:              db,
		registerOptions: initRegisterOptions(),
		repository:      r,
	}
}

func initRegisterOptions() *config.RegisterOptions {
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
	var err error
	if cansignup, err := u.UserCanSignUp(usr); !cansignup {
		return 0, err
	}

	tx := u.db.MustBegin()

	verificationToken, err := encryption.GenerateVerificationToken(usr.Email, u.registerOptions.UserOptions.EmailVerificationExpiryDuration, u.cfg)
	if err != nil {
		return 0, apierror.NewInternalServerApiError(errGenerateVerificationToken, err, errTokenGeneration)
	}

	user := &domain.User{
		Username:           usr.Username,
		NormalizedUsername: strings.ToUpper(usr.Username),
		VerificationToken:  verificationToken,
	}

	user.SecurityToken.String, err = encryption.GenerateSecurityToken(usr.Password, u.cfg)
	if err != nil {
		return 0, apierror.NewInternalServerApiError(errGenerateSecurityToken, err, errTokenGeneration)
	}

	user.PasswordHash, err = encryption.GenerateEncodedHash(usr.Password, u.cfg)
	if err != nil {
		return 0, apierror.NewInternalServerApiError(errPasswordHash, err, errPasswordHashCode)
	}

	userId, err := u.repository.AddUser(tx, user)
	if err != nil {
		return 0, apierror.NewInternalServerApiError(errAddingUser, err, errInvalidRegisterCode)
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
		return 0, apierror.NewInternalServerApiError(errAddingUserEmail, err, errInvalidRegisterCode)
	}

	// TODO: Send verification email

	tx.Commit()
	return userId, nil
}

func (u *registerService) UserCanSignUp(usr *domain.UserSignupDTO) (bool, apierror.ApiError) {
	errs := apierror.NewWithStatus(http.StatusBadRequest).WithMessage(errCantCreateUser)

	// Check that every field is correct
	if usr.Username == "" {
		return false, apierror.NewBadRequestApiError(domain.ErrEmptyUsername)
	}

	if usr.Password == "" {
		return false, apierror.NewBadRequestApiError(domain.ErrEmptyPassword)
	}

	if usr.Email == "" {
		return false, apierror.NewBadRequestApiError(domain.ErrEmptyEmail)
	}

	validEmail, err := regexp.Match("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,"+
		"61}[a-zA-Z0-9])?)*$", []byte(usr.Email))
	if !validEmail || err != nil {
		return false, apierror.NewBadRequestApiError(errInvalidEmailFormat)
	}

	if u.registerOptions.UserOptions.RequireUniqueEmail {
		exists, err := u.repository.CheckEmailExists(usr.Email)
		if exists {
			return false, apierror.NewBadRequestApiError(errEmailAlreadyExists)
		} else if err != nil && err != sql.ErrNoRows {
			return false, apierror.NewInternalServerApiError(errEmailAlreadyExistsInternalErr, err, domain.ErrInternalCode)
		}
	}

	usrexists, err := u.repository.CheckUsernameExists(usr.Username)
	if usrexists {
		return false, apierror.NewBadRequestApiError(errUserAlreadyExists)
	} else if err != nil && err != sql.ErrNoRows {
		return false, apierror.NewInternalServerApiError(errUserAlreadyExistsInternalErr, err, domain.ErrInternalCode)
	}

	usernameMatch, _ := regexp.Match(u.registerOptions.UserOptions.AllowedCharacters, []byte(usr.Username))
	if usernameMatch {
		errs.AddError(errUsernameCotainsIlegalChars, errInvalidUsernameCode)
	}

	if strings.Contains(usr.Password, " ") {
		errs.AddError(errPwContainsSpace, errInvalidPasswordCode)
	}

	if len(usr.Password) < u.registerOptions.PasswordOptions.RequiredLength {
		errs.AddError(fmt.Sprintf("El campo de contraseña tiene menos de %d caracteres", u.registerOptions.PasswordOptions.RequiredLength), errInvalidPasswordCode)
	}

	if u.registerOptions.PasswordOptions.RequireUppercase {
		match, _ := regexp.Match(".*[A-Z].*", []byte(usr.Password))
		if !match {
			errs.AddError(errPwDoesNotContainsUppercase, errInvalidPasswordCode)
		}
	}

	if u.registerOptions.PasswordOptions.RequireLowercase {
		match, _ := regexp.Match(".*[a-z].*", []byte(usr.Password))
		if !match {
			errs.AddError(errPwDoesNotContainsLowercase, errInvalidPasswordCode)
		}
	}

	// List of avalaible chars: ~!@#$%^&*()-+=?/<>|{}_:;.,
	if u.registerOptions.PasswordOptions.RequireNonAlphanumeric {
		match, _ := regexp.Match(".*[~!@#$%^&*()-+=?/<>|{}_:;.,].*", []byte(usr.Password))
		if !match {
			errs.AddError(errPwDoesNotContainsNonAlphaChars, errInvalidPasswordCode)
		}
	}

	if u.registerOptions.PasswordOptions.RequireDigit {
		match, _ := regexp.Match(".*\\d.*", []byte(usr.Password))
		if !match {
			errs.AddError(errPwDoesNotContainsADigit, errInvalidPasswordCode)
		}
	}

	if len(errs.Errors()) > 0 {
		return false, errs
	}

	return true, nil
}
