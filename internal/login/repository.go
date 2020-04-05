package login

import (
	"database/sql"
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

type loginRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) LoginRepository {
	return &loginRepository{db: db}
}

func (l *loginRepository) GetUserByUsername(username string) (*domain.User, *domain.UserEmail, apierror.ApiError) {
	var user domain.User

	err := l.db.Get(&user, "SELECT * FROM users where username = ?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
		}
		return nil, nil, apierror.New(http.StatusInternalServerError, ErrFailedTryingToLogin, apierror.NewErrorCause(err.Error(), ErrUserFetchFailed))
	}

	if user == (domain.User{}) {
		return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
	}

	var userEmail domain.UserEmail

	err = l.db.Get(&userEmail, "SELECT * FROM users_email WHERE user_id = ?", user.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
		}
		return nil, nil, apierror.New(http.StatusInternalServerError, ErrFailedTryingToLogin, apierror.NewErrorCause(err.Error(), ErrEmailFetchFailed))
	}

	if userEmail == (domain.UserEmail{}) {
		return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidEmail, ErrInvalidEmailCode))
	}

	return &user, &userEmail, nil
}

func (l *loginRepository) IncrementLoginFailAttempt(userId int64) error {
	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = failed_login_attempts + 1 where user_id = ?", userId)
	return err
}

func (l *loginRepository) ResetLoginFails(userId int64) error {
	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = 0 where user_id = ?", userId)
	return err
}

func (l *loginRepository) UnlockAccount(userId int64) error {
	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 0, lockout_date = null, failed_login_attempts = 0 where user_id = ?", userId)
	return err
}

func (l *loginRepository) LockAccount(userId int64, duration time.Duration) error {
	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 1, lockout_date = ? where user_id = ?", time.Now().Add(duration), userId)
	return err
}
