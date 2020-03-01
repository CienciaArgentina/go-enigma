package repositories

import (
	"errors"
	"github.com/CienciaArgentina/go-enigma/internal/login"
	"github.com/jmoiron/sqlx"
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

	err := l.db.Get(&user, "SELECT * FROM users where username = ?", username)
	if err != nil {
		return nil, nil, err
	}

	if user == (login.User{}) {
		return nil, nil, errInvalidLogin
	}

	var userEmail login.UserEmail

	err = l.db.Get(&userEmail, "SELECT * FROM users_email WHERE user_id = ?", user.UserId)
	if err != nil {
		return nil, nil, err
	}

	if userEmail == (login.UserEmail{}) {
		return nil, nil, errInvalidEmail
	}

	return &user, &userEmail, nil
}

func (l *loginRepository) IncrementLoginFailAttempt(userId int) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = failed_login_attempts + 1 where user_id = ?", userId)
	return err
}

func (l *loginRepository) ResetLoginFails(userId int) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = 0 where user_id = ?", userId)
	return err
}

func (l *loginRepository) UnlockAccount(userId int) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 0, lockout_date = null, failed_login_attempts = 0 where user_id = ?", userId)
	return err
}

func (l *loginRepository) LockAccount(userId int, duration time.Duration) error {
	if userId == 0 {
		return errUserIdMustBeGreaterThanZero
	}

	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 1, lockout_date = ? where user_id = ?", time.Now().Add(duration), userId)
	return err
}

func (l *loginRepository) GetUserRole(userId int) (string, error) {
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