package repositories

import (
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/recovery"
	"github.com/jmoiron/sqlx"
)

type recoveryRepository struct {
	db *sqlx.DB
}

func NewRecoveryRepository(db *sqlx.DB) recovery.Repository {
	return &recoveryRepository{db: db}
}

func (r *recoveryRepository) GetEmailByUserId(userId int64) (string, *recovery.UserEmail, error) {
	var user recovery.User

	err := r.db.Get(&user, "SELECT * FROM users where user_id = ?", userId)
	if err != nil {
		return "", nil, err
	}

	if user == (recovery.User{}) {
		return "", nil, nil
	}

	var userEmail recovery.UserEmail

	err = r.db.Get(&userEmail, "SELECT * FROM users_email where user_id = ?", userId)
	if err != nil {
		return "", nil, err
	}

	return user.VerificationToken, &userEmail, nil
}

func (r *recoveryRepository) ConfirmUserEmail(email string, token string) error {
	var userEmail recovery.UserEmail

	err := r.db.Get(&userEmail, "SELECT * FROM users_email where email = ?", email)
	if err != nil {
		return err
	}

	if userEmail == (recovery.UserEmail{}) {
		return config.ErrEmailValidationFailed
	}

	if userEmail.VerfiedEmail {
		return config.ErrEmailAlreadyVerified
	}

	var user recovery.User

	err = r.db.Get(&user, "SELECT * FROM users where user_id = ?", userEmail.UserId)
	if err != nil {
		return err
	}

	if user == (recovery.User{}) {
		return config.ErrEmailValidationFailed
	}

	if token != user.VerificationToken {
		return config.ErrValidationTokenFailed
	}

	result, err := r.db.Exec("UPDATE users_email SET verified_email = 1, verification_date = now() WHERE user_id = ?", user.UserId)
	if err != nil {
		return err
	}

	if num, _ := result.RowsAffected(); num == 0 {
		return config.ErrEmailValidationFailed
	}

	return nil
}