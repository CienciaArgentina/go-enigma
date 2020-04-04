package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/CienciaArgentina/go-enigma/internal_old/login"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	errInvalidLogin                = errors.New("El usuario o la contraseña especificados no existe")
	errInvalidEmail                = errors.New("Por alguna razón tu usuario no tiene email asociado")
	errUserIdMustBeGreaterThanZero = errors.New("El user id tiene que ser mayor que 0")
)

type loginRepository struct {
	db *sqlx.DB
}

func NewLoginRepository(db *sqlx.DB) login.Repository {
	return &loginRepository{db: db}
}

func (l *loginRepository) GetUserByUsername(username string) (*login.User, *login.UserEmail, error) {
	var user login.User
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Consultando si el usuario existe (SELECT * FROM users where username)")
	start := time.Now()

	err := l.db.Get(&user, "SELECT * FROM users where username = ?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	elapsed := time.Since(start)
	logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Terminó consulta sobre usuario")

	if user == (login.User{}) {
		return nil, nil, errInvalidLogin
	}

	var userEmail login.UserEmail

	logrus.Info("Consultando si el email existe (SELECT * FROM users_email WHERE user_id)")
	start = time.Now()

	err = l.db.Get(&userEmail, "SELECT * FROM users_email WHERE user_id = ?", user.UserId)
	if err != nil {
		return nil, nil, err
	}

	elapsed = time.Since(start)
	logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Terminó consulta sobre email")

	if userEmail == (login.UserEmail{}) {
		return nil, nil, errInvalidEmail
	}

	return &user, &userEmail, nil
}

func (l *loginRepository) IncrementLoginFailAttempt(userId int64) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Actualizando failed login attempts (UPDATE users SET failed_login_attempts = failed_login_attempts + 1 where user_id)")
	start := time.Now()

	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = failed_login_attempts + 1 where user_id = ?", userId)

	elapsed := time.Since(start)
	logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Terminó failed login attempts")
	return err
}

func (l *loginRepository) ResetLoginFails(userId int64) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = 0 where user_id = ?", userId)
	return err
}

func (l *loginRepository) UnlockAccount(userId int64) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 0, lockout_date = null, failed_login_attempts = 0 where user_id = ?", userId)
	return err
}

func (l *loginRepository) LockAccount(userId int64, duration time.Duration) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 1, lockout_date = ? where user_id = ?", time.Now().Add(duration), userId)
	return err
}

func (l *loginRepository) GetUserRole(userId int64) (string, error) {
	if userId == 0 {
		return "", errUserIdMustBeGreaterThanZero
	}

	var roleId int
	err := l.db.Get(&roleId, "SELECT role_id FROM user_roles where user_id = ?", userId)
	if err != nil {
		return "", err
	}

	var role string
	err = l.db.Get(&role, "SELECT name FROM roles where role_id = ?", roleId)

	return role, err
}
