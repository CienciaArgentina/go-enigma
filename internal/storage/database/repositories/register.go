package repositories

import (
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/jmoiron/sqlx"
	"strings"
)

type registerRepository struct {
	db *sqlx.DB
}

func NewRegisterRepository(db *sqlx.DB) register.Repository {
	return &registerRepository{db: db}
}

func (r *registerRepository) AddUser(u *register.User) (int64, error) {
	res, err := r.db.Exec("INSERT INTO users (`username`, `normalized_username`, `password_hash`, `lockout_enabled`,  "+
		"`date_created`, `verification_token`) VALUES ($1, $2, $3, $4, $5, $6)", u.Username, u.NormalizedUsername, u.PasswordHash, u.DateCreated, u.VerificationToken)

	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, err
}

func (r *registerRepository) AddEmail(e *register.UserEmail) (int64, error) {
	res, err := r.db.Exec("INSERT INTO users_email (`user_id`, `email`, `normalized_email`, `date_created`) VALUES ($1, $2, $3, $4)", e.UserId, e.Email, e.NormalizedEmail,
		e.DateCreated)
	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return lastId, err
}

func (r *registerRepository) VerifyIfEmailExists(email string) (bool, error) {
	normalizedemail := strings.ToUpper(email)
	var ue register.UserEmail
	err := r.db.Select(&ue, "SELECT * FROM users_email WHERE normalized_email = ?", normalizedemail)
	if err != nil {
		return true, err
	}

	return &ue == nil, err
}
