package register

import (
	"database/sql"
	"strings"

	"github.com/CienciaArgentina/go-enigma/internal/domain"

	"github.com/jmoiron/sqlx"
)

type registerRepository struct {
	db *sqlx.DB
}

// NewRepository Returns new productive register repository implementation
func NewRepository(db *sqlx.DB) RegisterRepository {
	return &registerRepository{db: db}
}

// GetUserById Returns user with given ID from repository
func (u *registerRepository) GetUserById(userId int64) (*domain.User, error) {
	var usr domain.User

	err := u.db.Get(&usr, "SELECT username FROM users WHERE user_id = ?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &usr, nil
}

// AddUser Creates a new user in repository
func (u *registerRepository) AddUser(tx *sqlx.Tx, usr *domain.User) (int64, error) {
	res, err := tx.Exec("INSERT INTO users (username, normalized_username, password_hash,  date_created, verification_token, security_token) VALUES (?, ?, ?, now(), ?, ?)",
		usr.Username, usr.NormalizedUsername, usr.PasswordHash, usr.VerificationToken, usr.SecurityToken)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, err
}

// AddUserEmail Updates user's email
func (u *registerRepository) AddUserEmail(tx *sqlx.Tx, e *domain.UserEmail) (int64, error) {
	res, err := tx.Exec("INSERT INTO users_email (user_id, email, normalized_email, date_created) VALUES (?, ?, ?, now())", e.UserId, e.Email, e.NormalizedEmail)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, err
}

// DeleteUser Deletes user's email
func (u *registerRepository) DeleteUser(userId int64) error {
	res, err := u.db.Exec("DELETE FROM users WHERE user_id = ?", userId)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	if rows == 0 {
		return errCannotDelete
	}

	return nil
}

// CheckUsernameExists Check's for username existence
func (u *registerRepository) CheckUsernameExists(username string) (bool, error) {
	var exists int
	err := u.db.Get(&exists, "SELECT count(*) FROM users where username = ?", username)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if exists > 0 {
		return true, nil
	}

	return false, nil
}

// CheckEmailExists Check's user's email existence
func (u *registerRepository) CheckEmailExists(email string) (bool, error) {
	var exists int
	err := u.db.Get(&exists, "SELECT count(*) FROM users_email WHERE normalized_email = ?", strings.ToUpper(email))
	if err != nil && err != sql.ErrNoRows {
		return false, err
	} else if exists > 0 {
		return true, nil
	}

	return false, nil
}
