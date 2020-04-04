package user

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
)

var (
	errCannotDelete = errors.New("El usuario que se intenta borrar no existe o no se puede alcanzar")
)

type userService struct {
	cfg             *config.Configuration
	db              *sqlx.DB
	registerOptions *config.RegisterOptions
	repository      domain.UserRepository
}

func NewService(c *config.Configuration, db *sqlx.DB) domain.UserService {
	return &userService{
		cfg: c,
		db:  db,
	}
}

func (u *userService) CreateUser(usr *domain.UserDTO) (int64, apierror.ApiError) {
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
		return 0, apierror.New(http.StatusInternalServerError, config.ErrPasswordHash, apierror.NewErrorCause(config.ErrPasswordHash, config.ErrPasswordHashCode))
	}

	userId, err := u.repository.AddUser(tx, user)
	if err != nil {
		return 0, apierror.New(http.StatusInternalServerError, config.ErrAddingUser, apierror.NewErrorCause(err.Error(), config.ErrInvalidRegisterCode))
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
		return 0, apierror.New(http.StatusInternalServerError, config.ErrAddingUserEmail, apierror.NewErrorCause(err.Error(), config.ErrInvalidRegisterCode))
	}

	// TODO: Send verification email

	return userId, nil
}

func (u *userService) UserCanSignUp(usr *domain.UserDTO) (bool, apierror.ApiError) {
	errs := apierror.NewWithStatus(http.StatusBadRequest).WithMessage(config.ErrCantCreateUser)

	// Check that every field is correct
	if usr.Username == "" {
		return false, apierror.New(http.StatusBadRequest, config.ErrEmptyUsername, apierror.NewErrorCause(config.ErrEmptyUsername, config.ErrEmptyFieldUserCode))
	}

	if usr.Password == "" {
		return false, apierror.New(http.StatusBadRequest, config.ErrEmptyPassword, apierror.NewErrorCause(config.ErrEmptyPassword, config.ErrEmptyFieldUserCode))
	}

	if usr.Email == "" {
		return false, apierror.New(http.StatusBadRequest, config.ErrEmptyEmail, apierror.NewErrorCause(config.ErrEmptyEmail, config.ErrEmptyFieldUserCode))
	}

	validEmail, err := regexp.Match("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,"+
		"61}[a-zA-Z0-9])?)*$", []byte(usr.Email))
	if !validEmail || err != nil {
		return false, apierror.New(http.StatusBadRequest, config.ErrInvalidEmailFormat, apierror.NewErrorCause(config.ErrInvalidEmailFormat, config.ErrInvalidEmailFormatCode))
	}

	if u.registerOptions.UserOptions.RequireUniqueEmail {
		exists, err := u.repository.CheckEmailExists(usr.Email)
		if exists {
			return false, apierror.New(http.StatusBadRequest, config.ErrEmailAlreadyExists, apierror.NewErrorCause(config.ErrEmailAlreadyExists, config.ErrEmailAlreadyExistsCode))
		} else if err != nil && err != sql.ErrNoRows {
			return false, apierror.New(http.StatusInternalServerError, config.ErrEmailAlreadyExistsInternalErr, apierror.NewErrorCause(err.Error(), config.ErrEmailAlreadyExistsInternalErrCode))
		}
	}

	usrexists, err := u.repository.CheckUsernameExists(usr.Username)
	if usrexists {
		return false, apierror.New(http.StatusBadRequest, config.ErrUserAlreadyExists, apierror.NewErrorCause(config.ErrUserAlreadyExists, config.ErrUserAlreadyExistsCode))
	} else if err != nil && err != sql.ErrNoRows {
		return false, apierror.New(http.StatusInternalServerError, config.ErrUserAlreadyExistsInternalErr, apierror.NewErrorCause(err.Error(), config.ErrUserAlreadyExistsInternalErrCode))
	}

	usernameMatch, _ := regexp.Match(u.registerOptions.UserOptions.AllowedCharacters, []byte(usr.Username))
	if usernameMatch {
		errs.AddError(config.ErrUsernameCotainsIlegalChars, config.ErrInvalidUsernameCode)
	}

	if strings.Contains(usr.Password, " ") {
		errs.AddError(config.ErrPwContainsSpace, config.ErrInvalidPasswordCode)
	}

	if len(usr.Password) < u.registerOptions.PasswordOptions.RequiredLength {
		errs.AddError(fmt.Sprintf("El campo de contraseña tiene menos de %d caracteres", u.registerOptions.PasswordOptions.RequiredLength), config.ErrInvalidPasswordCode)
	}

	if u.registerOptions.PasswordOptions.RequireUppercase {
		match, _ := regexp.Match(".*[A-Z].*", []byte(usr.Password))
		if !match {
			errs.AddError(config.ErrPwDoesNotContainsUppercase, config.ErrInvalidPasswordCode)
		}
	}

	if u.registerOptions.PasswordOptions.RequireLowercase {
		match, _ := regexp.Match(".*[a-z].*", []byte(usr.Password))
		if !match {
			errs.AddError(config.ErrPwDoesNotContainsLowercase, config.ErrInvalidPasswordCode)
		}
	}

	// List of avalaible chars: ~!@#$%^&*()-+=?/<>|{}_:;.,
	if u.registerOptions.PasswordOptions.RequireNonAlphanumeric {
		match, _ := regexp.Match(".*[~!@#$%^&*()-+=?/<>|{}_:;.,].*", []byte(usr.Password))
		if !match {
			errs.AddError(config.ErrPwDoesNotContainsNonAlphaChars, config.ErrInvalidPasswordCode)
		}
	}

	if u.registerOptions.PasswordOptions.RequireDigit {
		match, _ := regexp.Match(".*\\d.*", []byte(usr.Password))
		if !match {
			errs.AddError(config.ErrPwDoesNotContainsADigit, config.ErrInvalidPasswordCode)
		}
	}

	if len(errs.ErrError) > 0 {
		return false, errs
	}

	return true, nil
}
