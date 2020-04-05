package login

import (
	"crypto/subtle"
	"fmt"
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/config"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/CienciaArgentina/go-enigma/internal/encryption"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
	"net/http"
	"time"
)

const (
	// Login

	// Internal server error
	ErrFailedTryingToLogin     = "Ocurrió un error al momento de loguear, intentá nuevamente o comunicate con sistemas"

	// Failed
	ErrInvalidLoginCode = "invalid_login"
	ErrInvalidLogin     = "El usuario o la contraseña especificados no existe"

	// Locked account
	ErrLockedAccountCode  = "locked_account"
	ErrLockedManyAttempts = "locked_many_attempts"

	// Failed password decryption
	ErrFailedPasswordDecryptionCode = "failed_decryption"

	// Email not verified
	ErrEmailNotVerified                = "Tu dirección de email no fue verificada aún"
	ErrEmailNotVerifiedCode  = "email_not_verified"

	// Invalid Email
	ErrInvalidEmail = "El mail no se encuentra registrado"
	ErrInvalidEmailCode = "invalid_email"

	// User fetch failed
	ErrUserFetchFailed = "user_fetch_failed"
	// Email fetch failed
	ErrEmailFetchFailed = "email_fetch_failed"
)

type loginService struct {
	cfg          *config.Configuration
	loginOptions *config.LoginOptions
	repository   LoginRepository
}

func NewService(c *config.Configuration, l *config.LoginOptions, r LoginRepository) LoginService {
	if l == nil {
		l = defaultLoginOptions()
	}
	return &loginService{
		cfg:          c,
		loginOptions: l,
		repository: r,
	}
}

func defaultLoginOptions() *config.LoginOptions {
	o := config.LoginOptions{}

	o.LockoutOptions.LockoutTimeDuration = 5 * time.Minute
	o.LockoutOptions.MaxFailedAttempts = 5

	o.SignInOptions.RequireConfirmedEmail = true

	return &o
}

func (l *loginService) LoginUser(u *domain.UserLoginDTO) (string, apierror.ApiError) {
	if err := l.UserCanLogin(u); err != nil {
		return "", err
	}

	user, userEmail, err := l.repository.GetUserByUsername(u.Username)
	if err != nil {
		return "", err
	}

	if user == nil || user == (&domain.User{}) || userEmail == nil || userEmail == (&domain.UserEmail{}) {
		return "", apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
	}

	verifyPassword, err := comparePasswordAndHash(u.Password, user.PasswordHash)
	if err != nil {
		// Return friendly message
		return "", err
	}

	if user.LockoutEnabled {
		// If the register is locked but time is up we should unlock the account
		if user.LockoutDate.Time.Add(l.loginOptions.LockoutOptions.LockoutTimeDuration).Before(time.Now()) {
			user.FailedLoginAttempts = 0
			user.LockoutEnabled = false
			err := l.repository.UnlockAccount(user.UserId)
			if err != nil {
				// TODO: Log this
			}
		} else {
			friendlyMessage := fmt.Sprintf("La cuenta se encuentra bloqueada por %v minutos por intentos fallidos de login",
				l.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
			return "", apierror.New(http.StatusBadRequest, friendlyMessage, apierror.NewErrorCause(friendlyMessage, ErrLockedAccountCode))
		}
	}

	if !verifyPassword {
		if user.FailedLoginAttempts >= l.loginOptions.LockoutOptions.MaxFailedAttempts {
			err := l.repository.LockAccount(user.UserId, l.loginOptions.LockoutOptions.LockoutTimeDuration)
			if err != nil {
				// TODO: Log this
			}
			friendlyMsg := fmt.Sprintf("Debido a repetidos intentos tu cuenta fue bloqueada por %v minutos", l.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
			return "", apierror.New(http.StatusBadRequest, friendlyMsg, apierror.NewErrorCause(friendlyMsg, ErrLockedManyAttempts))
		}
		err := l.repository.IncrementLoginFailAttempt(user.UserId)
		if err != nil {
			// TODO: Log this
		}
		return "", apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
	}

	if l.loginOptions.SignInOptions.RequireConfirmedEmail && !userEmail.VerfiedEmail {
		return "", apierror.New(http.StatusBadRequest, ErrEmailNotVerified, apierror.NewErrorCause(ErrEmailNotVerified, ErrEmailNotVerifiedCode))
	}

	e := l.repository.ResetLoginFails(user.UserId)
	if e != nil {
		// TODO: Log this
	}

	// TODO: Add role

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.UserId,
		"email":  userEmail.Email,
		"timestamp": time.Now().Unix(),
	})

	jwtString, _ := jwt.SignedString([]byte(l.cfg.Keys.PasswordHashingKey))

	return jwtString, nil
}

func (l *loginService) UserCanLogin(u *domain.UserLoginDTO) apierror.ApiError {
	if u.Username == "" {
		return apierror.New(http.StatusBadRequest, config.ErrEmptyUsername, apierror.NewErrorCause(config.ErrEmptyUsername, config.ErrEmptyFieldUserCodeLogin))
	}

	if u.Password == "" {
		return apierror.New(http.StatusBadRequest, config.ErrEmptyPassword, apierror.NewErrorCause(config.ErrEmptyPassword, config.ErrEmptyFieldUserCodeLogin))
	}

	return nil
}

func comparePasswordAndHash(password, encodedHash string) (bool, apierror.ApiError) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := encryption.DecodeHash(encodedHash)
	if err != nil {
		return false, apierror.New(http.StatusInternalServerError, ErrFailedTryingToLogin, apierror.NewErrorCause(err.Error(), ErrFailedPasswordDecryptionCode))
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
