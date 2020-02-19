package login

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/CienciaArgentina/go-enigma/config"
	jwt2 "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/argon2"
	"strings"
	"time"
)

var (
	errEmptyUsername       = errors.New("El nombre de usuario es necesario para el login")
	errEmptyPassword       = errors.New("La contraseña es necesaria para el login")
	errInvalidLogin        = errors.New("El usuario o la contraseña especificados no existe")
	errInvalidHash         = errors.New("El hash no usa el encoding correcto")
	errIncompatibleVersion = errors.New("Versión de argon2 incompatible")
	errThroughLogin        = errors.New("Ocurrió un error al momento de loguear")
	errInvalidEmail        = errors.New("Por alguna razón tu usuario no tiene email asociado")
	errEmailNotVerified = errors.New("Tu dirección de email no fue verificada aún")
)

type LoginOptions struct {
	LockoutOptions struct {
		LockoutTimeDuration time.Duration
		MaxFailedAttempts   int
	}
	SignInOptions struct {
		RequireConfirmedEmail bool
	}
}

type Service interface {
	Login(u *UserLogin) (string, error)
	VerifyCanLogin(u *UserLogin) (bool, error)
}

type loginService struct {
	config       *config.Configuration
	repo         Repository
	loginOptions *LoginOptions
}

func NewService(r Repository, l *LoginOptions, cfg *config.Configuration) Service {
	if l == (&LoginOptions{}) || l == nil {
		l = defaultLoginOptions()
	}
	return &loginService{
		repo:         r,
		config:       cfg,
		loginOptions: l,
	}
}

func defaultLoginOptions() *LoginOptions {
	o := LoginOptions{}

	o.LockoutOptions.LockoutTimeDuration = 5 * time.Minute
	o.LockoutOptions.MaxFailedAttempts = 5

	o.SignInOptions.RequireConfirmedEmail = true

	return &o
}

func (s *loginService) Login(u *UserLogin) (string, error) {
	canLogin, err := s.VerifyCanLogin(u)
	if !canLogin {
		return "", err
	}

	user, userEmail := s.repo.GetUserByUsername(u.Username)
	if user == nil || user == (&User{}) {
		return "", errInvalidLogin
	}

	if userEmail == nil || userEmail == (&UserEmail{}) {
		return "", errInvalidEmail
	}

	verifyPassword, err := comparePasswordAndHash(u.Password, user.PasswordHash)
	if err != nil {
		// Return friendly message
		return "", errThroughLogin
	}

	if !verifyPassword {
		if user.FailedLoginAttempts >= s.loginOptions.LockoutOptions.MaxFailedAttempts {
			err := s.repo.LockAccount(user.UserId)
			if err != nil {
				// TODO: Log this
			}
			return "", fmt.Errorf("Debido a repetidos intentos tu cuenta fue bloqueada por %v minutos", s.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
		}
		err := s.repo.IncrementLoginFailAttempt(user.UserId)
		if err != nil {
			// TODO: Log this
		}
		return "", errInvalidLogin
	}

	if user.LockoutEnabled {
		// If the user is locked but time is up we should unlock the account
		if user.LockoutDate.Add(s.loginOptions.LockoutOptions.LockoutTimeDuration).After(time.Now()) {
			err := s.repo.UnlockAccount(user.UserId)
			if err != nil {
				// TODO: Log this
			}
		} else {
			return "", fmt.Errorf("La cuenta se encuentra bloqueada por %v minutos por intentos fallidos de login", s.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
		}
	}

	if s.loginOptions.SignInOptions.RequireConfirmedEmail && !userEmail.VerfiedEmail {
		return "", errEmailNotVerified
	}

	err = s.repo.ResetLoginFails(user.UserId)
	if err != nil {
		// TODO: Log this
	}

	role, err := s.repo.GetUserRole(user.UserId)
	if err != nil {
		// TODO: log this
	}

	jwt := jwt2.NewWithClaims(jwt2.SigningMethodHS256, jwt2.MapClaims{
		"userId": user.UserId,
		"email": userEmail.Email,
		"role": role,
	})

	jwtString, _ := jwt.SignedString(s.config.Keys.PasswordHashingKey)

	return jwtString, nil
}

func (s *loginService) VerifyCanLogin(u *UserLogin) (bool, error) {
	if u.Username == "" {
		return false, errEmptyUsername
	}

	if u.Password == "" {
		return false, errEmptyPassword
	}

	return true, nil
}

func comparePasswordAndHash(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
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

func decodeHash(encodedHash string) (p *config.ArgonParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		// TODO: Log this
		return nil, nil, nil, errInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		// TODO: Log this
		return nil, nil, nil, errIncompatibleVersion
	}

	p = &config.ArgonParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}
