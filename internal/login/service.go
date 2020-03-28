package login

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/CienciaArgentina/go-enigma/config"
	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
	"strings"
	"time"
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
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Iniciando service de Login")
	canLogin, err := s.VerifyCanLogin(u)
	if !canLogin {
		return "", err
	}

	logrus.Info("Iniciando GetUserByUsername")
	start := time.Now()

	user, userEmail, err := s.repo.GetUserByUsername(u.Username)
	if err != nil {
		return "", err
	}

	elapsed := time.Since(start)
	logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Terminó GetUserByUsername")

	if user == nil || user == (&User{}) || userEmail == nil || userEmail == (&UserEmail{}) {
		return "", config.ErrInvalidLogin
	}

	logrus.Info("Iniciando evaluación del hash del password")
	start = time.Now()

	verifyPassword, err := comparePasswordAndHash(u.Password, user.PasswordHash)
	if err != nil {
		// Return friendly message
		return "", config.ErrThroughLogin
	}

	elapsed = time.Since(start)
	logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Terminó evaluación del hash del password")

	if user.LockoutEnabled {
		// If the user is locked but time is up we should unlock the account
		a := user.LockoutDate.Time.Add(s.loginOptions.LockoutOptions.LockoutTimeDuration)
		fmt.Print(a)
		if user.LockoutDate.Time.Add(s.loginOptions.LockoutOptions.LockoutTimeDuration).Before(time.Now()) {
			user.FailedLoginAttempts = 0
			user.LockoutEnabled = false
			err := s.repo.UnlockAccount(user.UserId)
			if err != nil {
				// TODO: Log this
			}
		} else {
			return "", fmt.Errorf("La cuenta se encuentra bloqueada por %v minutos por intentos fallidos de login", s.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
		}
	}

	if !verifyPassword {
		if user.FailedLoginAttempts >= s.loginOptions.LockoutOptions.MaxFailedAttempts {
			err := s.repo.LockAccount(user.UserId, s.loginOptions.LockoutOptions.LockoutTimeDuration)
			if err != nil {
				// TODO: Log this
			}
			return "", fmt.Errorf("Debido a repetidos intentos tu cuenta fue bloqueada por %v minutos", s.loginOptions.LockoutOptions.LockoutTimeDuration.Minutes())
		}

		logrus.Info("Iniciando IncrementLoginFailAttempt")
		start = time.Now()

		err := s.repo.IncrementLoginFailAttempt(user.UserId)
		if err != nil {
			// TODO: Log this
		}

		elapsed = time.Since(start)
		logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Terminó IncrementLoginFailAttempt")
		return "", config.ErrInvalidLogin
	}

	if s.loginOptions.SignInOptions.RequireConfirmedEmail && !userEmail.VerfiedEmail {
		return "", config.ErrEmailNotVerified
	}

	err = s.repo.ResetLoginFails(user.UserId)
	if err != nil {
		// TODO: Log this
	}

	// TODO: add role
	//role, err := s.repo.GetUserRole(user.UserId)
	//if err != nil {
	//	// TODO: log this
	//}

	// TODO: add role
	jwt := jwt2.NewWithClaims(jwt2.SigningMethodHS256, jwt2.MapClaims{
		"userId": user.UserId,
		"email":  userEmail.Email,
		//"role":   role,
		"timestamp": time.Now().Unix(),
	})

	jwtString, _ := jwt.SignedString([]byte(s.config.Keys.PasswordHashingKey))

	return jwtString, nil
}

func (s *loginService) VerifyCanLogin(u *UserLogin) (bool, error) {
	if u.Username == "" {
		return false, config.ErrEmptyUsername
	}

	if u.Password == "" {
		return false, config.ErrEmptyPassword
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
		return nil, nil, nil, config.ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		// TODO: Log this
		return nil, nil, nil, config.ErrIncompatibleVersion
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
