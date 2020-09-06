package login

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	domain2 "github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/jmoiron/sqlx"
)

type loginRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &loginRepository{db: db}
}

func (l *loginRepository) GetUserByUsername(username string) (*domain2.User, *domain2.UserEmail, apierror.ApiError) {
	var user domain2.User

	err := l.db.Get(&user, "SELECT * FROM users where username = ?", username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
		}
		return nil, nil, apierror.New(http.StatusInternalServerError, ErrFailedTryingToLogin, apierror.NewErrorCause(err.Error(), ErrUserFetchFailed))
	}

	if user == (domain2.User{}) {
		return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
	}

	var userEmail domain2.UserEmail

	err = l.db.Get(&userEmail, "SELECT * FROM users_email WHERE user_id = ?", user.AuthId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidLogin, ErrInvalidLoginCode))
		}
		return nil, nil, apierror.New(http.StatusInternalServerError, ErrFailedTryingToLogin, apierror.NewErrorCause(err.Error(), ErrEmailFetchFailed))
	}

	if userEmail == (domain2.UserEmail{}) {
		return nil, nil, apierror.New(http.StatusBadRequest, ErrInvalidLogin, apierror.NewErrorCause(ErrInvalidEmail, ErrInvalidEmailCode))
	}

	return &user, &userEmail, nil
}

func (l *loginRepository) IncrementLoginFailAttempt(userID int64) error {
	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = failed_login_attempts + 1 where user_id = ?", userID)
	return err
}

func (l *loginRepository) ResetLoginFails(userID int64) error {
	_, err := l.db.Exec("UPDATE users SET failed_login_attempts = 0 where user_id = ?", userID)
	return err
}

func (l *loginRepository) UnlockAccount(userID int64) error {
	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 0, lockout_date = null, failed_login_attempts = 0 where user_id = ?", userID)
	return err
}

func (l *loginRepository) LockAccount(userID int64, duration time.Duration) error {
	_, err := l.db.Exec("UPDATE users SET lockout_enabled = 1, lockout_date = ? where user_id = ?", time.Now().Add(duration), userID)
	return err
}
