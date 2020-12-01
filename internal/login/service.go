package login

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"

	"github.com/CienciaArgentina/go-backend-commons/pkg/performance"

	"github.com/CienciaArgentina/go-backend-commons/pkg/clog"

	"github.com/go-resty/resty/v2"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/CienciaArgentina/go-enigma/internal/encryption"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
)

const (
	// Login.

	// Internal server error.
	ErrFailedTryingToLogin = "Ocurrió un error al momento de loguear, intentá nuevamente o comunicate con sistemas"

	// Failed.
	ErrInvalidLoginCode = "invalid_login"
	ErrInvalidLogin     = "El usuario o la contraseña especificados no existe"

	// Locked account.
	ErrLockedAccountCode  = "locked_account"
	ErrLockedManyAttempts = "locked_many_attempts"

	// Failed password decryption.
	errDecrypt = "failed_decryption" // nolint

	// Email not verified.
	ErrEmailNotVerified     = "Tu dirección de email no fue confirmada aún"
	ErrEmailNotVerifiedCode = "email_not_verified"

	// Invalid Email.
	ErrInvalidEmail     = "El mail no se encuentra registrado"
	ErrInvalidEmailCode = "invalid_email"

	// User fetch failed.
	ErrUserFetchFailed = "user_fetch_failed"
	// Email fetch failed.
	ErrEmailFetchFailed = "email_fetch_failed"
)

type loginService struct {
	cfg          *config.EnigmaConfig
	loginOptions *config.LoginOptions
	repository   Repository
}

func NewService(cfg *config.EnigmaConfig, r Repository) Service {
	return &loginService{
		cfg:          cfg,
		loginOptions: setLoginOptions(),
		repository:   r,
	}
}

func setLoginOptions() *config.LoginOptions {
	o := config.LoginOptions{}

	o.LockoutOptions.LockoutTimeDuration = 5 * time.Minute
	o.LockoutOptions.MaxFailedAttempts = 5

	o.SignInOptions.RequireConfirmedEmail = true

	return &o
}

func (l *loginService) LoginUser(u *domain.UserLoginDTO, ctx *rest.ContextInformation) (string, apierror.ApiError) { // nolint
	var err error
	var apierr apierror.ApiError
	var user *domain.User
	var userEmail *domain.UserEmail

	performance.TrackTime(time.Now(), "UserCanLogin", ctx, func() {
		apierr = l.UserCanLogin(u)
	})

	if apierr != nil {
		return "", apierr
	}

	performance.TrackTime(time.Now(), "GetUserByUsername", ctx, func() {
		user, userEmail, apierr = l.repository.GetUserByUsername(u.Username)
	})
	if apierr != nil {
		clog.Error("Error al obtener el username", "login-user", err, nil)
		return "", apierr
	}

	if user == nil || userEmail == nil {
		return "", apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
	}

	var verifyPassword bool
	performance.TrackTime(time.Now(), "comparePasswordAndHash", ctx, func() {
		verifyPassword, err = comparePasswordAndHash(u.Password, user.PasswordHash)
	})
	if err != nil {
		// Return friendly message
		clog.Error("Error comparing password", "login-user", err, nil)
		return "", apierror.NewInternalServerApiError(domain.ErrUnexpectedError, err, domain.ErrInternalCode)
	}

	if user.LockoutEnabled {
		// If the register is locked but time is up we should unlock the account
		if user.LockoutDate.Time.Add(l.loginOptions.LockoutOptions.LockoutTimeDuration).Before(time.Now()) {
			user.FailedLoginAttempts = 0
			user.LockoutEnabled = false
			err := l.repository.UnlockAccount(user.AuthId)
			if err != nil {
				clog.Error("Can't unlock account", "login-user", err, map[string]string{"auth_id": fmt.Sprintf("%d", user.AuthId)})
			}
		} else {
			friendlyMessage := fmt.Sprintf("La cuenta se encuentra bloqueada por %v minutos por intentos fallidos de login",
				l.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
			return "", apierror.NewBadRequestApiError(friendlyMessage)
		}
	}

	if !verifyPassword {
		if user.FailedLoginAttempts >= l.loginOptions.LockoutOptions.MaxFailedAttempts {
			err := l.repository.LockAccount(user.AuthId, l.loginOptions.LockoutOptions.LockoutTimeDuration)
			if err != nil {
				clog.Error("Can't lock account", "login-user", err, map[string]string{"auth_id": fmt.Sprintf("%d", user.AuthId)})
			}
			friendlyMsg := fmt.Sprintf("Debido a repetidos intentos tu cuenta fue bloqueada por %v minutos", l.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
			return "", apierror.NewBadRequestApiError(friendlyMsg)
		}
		err := l.repository.IncrementLoginFailAttempt(user.AuthId)
		if err != nil {
			clog.Error("Can't increment login fail attemp", "login-user", err, map[string]string{"auth_id": fmt.Sprintf("%d", user.AuthId)})
		}
		return "", apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
	}

	if l.loginOptions.SignInOptions.RequireConfirmedEmail && !userEmail.VerfiedEmail {
		return "", apierror.New(http.StatusBadRequest, ErrEmailNotVerified, apierror.NewErrorCause(userEmail.Email, ErrEmailNotVerifiedCode))
	}

	performance.TrackTime(time.Now(), "ResetLoginFails", ctx, func() {
		err = l.repository.ResetLoginFails(user.AuthId)
	})
	if err != nil {
		clog.Error("can't reset login fails", "login-user", err, map[string]string{"auth_id": fmt.Sprintf("%d", user.AuthId)})
	}

	var role *domain.AssignedRole
	var gErr error
	performance.TrackTime(time.Now(), "getRole", ctx, func() {
		role, gErr = getRole(user.AuthId, ctx)
	})
	if gErr != nil {
		return "", apierror.NewInternalServerApiError("Cannot get role", gErr, "get_role")
	}

	roleb, mErr := json.Marshal(role.Roles)
	if mErr != nil {
		return "", apierror.NewInternalServerApiError("Cannot get role", mErr, "marshal_role")
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"auth_id":   user.AuthId,
		"email":     userEmail.Email,
		"timestamp": time.Now().Unix(),
		"role":      string(roleb),
	})

	jwtString, _ := jwt.SignedString([]byte(l.cfg.JwtSign))

	return jwtString, nil

}

func (l *loginService) UserCanLogin(u *domain.UserLoginDTO) apierror.ApiError {
	if u.Username == "" {
		return apierror.NewBadRequestApiError(domain.ErrEmptyUsername)
	}

	if u.Password == "" {
		return apierror.NewBadRequestApiError(domain.ErrEmptyPassword)
	}

	return nil
}

func comparePasswordAndHash(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := encryption.DecodeHash(encodedHash)
	if err != nil {
		return false, err
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

func getRole(authid int64, ctx *rest.ContextInformation) (*domain.AssignedRole, error) {
	var role *domain.AssignedRole
	var res *resty.Response
	var err error
	baseURL := domain.GetRolesBaseURL()
	authstr := strconv.FormatInt(authid, 10)

	// TODO: Move this to a client
	performance.TrackTime(time.Now(), "GetRoleAPICall", ctx, func() {
		res, err = resty.New().SetHostURL(baseURL).R().SetPathParams(map[string]string{"auth_id": authstr}).Get("/assign/{auth_id}")
	})

	if err != nil {
		clog.Error("Rest client error", "get-role", err, nil)
	}
	if res.IsError() {
		clog.Error("Status error - GetRole", "get-role", errors.New("status error - GetRole"), map[string]string{"status": res.Status()})
		return nil, errors.New(res.String())
	}
	err = json.Unmarshal(res.Body(), &role)
	if err != nil {
		clog.Error("Unmarshal error - GetRole", "get-role", err, nil)
		return nil, err
	}

	return role, nil
}
