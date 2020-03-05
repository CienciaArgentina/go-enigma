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

func (r *recoveryRepository) GetuserIdByEmail(email string) (int64, error) {
	var userId int64

	err := r.db.Get(&userId, "SELECT user_id FROM users_email where email = ?", email)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *recoveryRepository) GetUsernameByEmail(email string) (string, error) {
	var userId int64

	err := r.db.Get(&userId, "SELECT user_id FROM users_email WHERE email = ?", email)
	if err != nil {
		return "", nil
	}

	var username string
	err = r.db.Get(&username, "SELECT username FROM users where user_id = ?", userId)
	if err != nil {
		return "", nil
	}

	return username, nil
}

func (r *recoveryRepository) GetSecurityToken(email string) (string, error) {
	var userEmail recovery.UserEmail

	err := r.db.Get(&userEmail, "SELECT * FROM users_email where email = ?", email)
	if err != nil {
		return "", err
	}

	if userEmail == (recovery.UserEmail{}) {
		return "", config.ErrEmptySearch
	}

	if !userEmail.VerfiedEmail {
		return "", config.ErrEmailNotVerified
	}

	var securityToken string

	err = r.db.Get(&securityToken, "SELECT security_token FROM users where user_id = ?", userEmail.UserId)
	if err != nil {
		return "", err
	}

	return securityToken, nil
}

func (r *recoveryRepository) UpdatePasswordAndResetSecurityToken(userId int64, passwordHash, newSecurityToken string) (bool, error) {
	if passwordHash == "" || newSecurityToken == "" {
		return false, config.ErrEmptyField
	}

	result := r.db.MustExec("UPDATE users SET password_hash = ?, security_token = ? where user_id = ?", passwordHash, newSecurityToken, userId)

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if updatedRows == 0 {
		return false, config.ErrUnexpectedError
	}

	return true, nil
}
