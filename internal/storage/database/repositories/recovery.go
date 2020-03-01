package repositories

import (
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
