package recovery

import (
	"database/sql"
	"net/http"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/jmoiron/sqlx"
)

const (
	ErrNoUserId     = "AuthId inexistente"
	ErrNoUserIdCode = "invalid_user_id"

	ErrNoUserEmail     = "Email inexistente"
	ErrNoUserEmailCode = "invalid_email"

	ErrFetchingUserCode      = "error_fetching_user"
	ErrFetchingUserEmailCode = "error_fetching_email"

	ErrEmailAlreadyverified     = "Este email ya se encuentra verificado"
	ErrEmailAlreadyverifiedCode = "email_already_verified"

	ErrEmailNotVerified     = "Este email no se encuentra verificado"
	ErrEmailNotVerifiedCode = "email_not_verified"

	ErrValidationTokenFailed     = "La validación del token falló"
	ErrValidationTokenFailedCode = "token_validation_failed"

	ErrUpdatingUserEmail     = "Ocurrió un error al intentar actualizar el email"
	ErrUpdatingUserEmailCode = "email_update_failed"

	ErrUpdatingUser     = "Ocurrió un error al intentar actualizar el usuario"
	ErrUpdatingUserCode = "user_update_failed"
)

type recoveryRepository struct {
	db *sqlx.DB
}

// NewRepository Creates new repository with injected DB
func NewRepository(db *sqlx.DB) RecoveryRepository {
	return &recoveryRepository{db: db}
}

// GetEmailByUserId Retrieves user's email based on it's ID
func (r *recoveryRepository) GetEmailByUserId(userId int64) (string, *domain.UserEmail, apierror.ApiError) {
	var user domain.User

	err := r.db.Get(&user, "SELECT * FROM users where user_id = ?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, apierror.New(http.StatusBadRequest, ErrNoUserId, apierror.NewErrorCause(ErrNoUserId, ErrNoUserIdCode))
		}

		return "", nil, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserCode))
	}

	var userEmail domain.UserEmail

	err = r.db.Get(&userEmail, "SELECT * FROM users_email where user_id = ?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, apierror.New(http.StatusBadRequest, ErrNoUserEmail, apierror.NewErrorCause(ErrNoUserEmail, ErrNoUserEmailCode))
		}

		return "", nil, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserEmailCode))
	}

	return user.VerificationToken, &userEmail, nil
}

func (r *recoveryRepository) ConfirmUserEmail(email, token string) apierror.ApiError {
	var userEmail domain.UserEmail

	err := r.db.Get(&userEmail, "SELECT * FROM users_email where email = ?", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(http.StatusBadRequest, ErrNoUserEmail, apierror.NewErrorCause(ErrNoUserEmail, ErrNoUserEmailCode))
		}
		return apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserEmailCode))
	}

	if userEmail.VerfiedEmail {
		return apierror.New(http.StatusBadRequest, ErrEmailAlreadyverified, apierror.NewErrorCause(ErrEmailAlreadyVerified, ErrEmailAlreadyverifiedCode))
	}

	var user domain.User

	err = r.db.Get(&user, "SELECT * FROM users where user_id = ?", userEmail.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return apierror.New(http.StatusBadRequest, ErrNoUserId, apierror.NewErrorCause(ErrNoUserId, ErrNoUserIdCode))
		}
		return apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserCode))
	}

	if token != user.VerificationToken {
		return apierror.New(http.StatusBadRequest, ErrValidationTokenFailed, apierror.NewErrorCause(ErrValidationTokenFailed, ErrValidationTokenFailedCode))
	}

	result, err := r.db.Exec("UPDATE users_email SET verified_email = 1, verification_date = now() WHERE user_id = ?", user.AuthId)
	if err != nil {
		return apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrUpdatingUserEmailCode))
	}

	if num, _ := result.RowsAffected(); num == 0 {
		return apierror.New(http.StatusBadRequest, ErrUpdatingUserEmail, apierror.NewErrorCause(ErrUpdatingUserEmail, ErrUpdatingUserEmailCode))
	}

	return nil
}

// GetuserIdByEmail Returns user's ID from given Email
func (r *recoveryRepository) GetuserIdByEmail(email string) (int64, apierror.ApiError) {
	var userId int64

	err := r.db.Get(&userId, "SELECT user_id FROM users_email where email = ?", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, apierror.New(http.StatusBadRequest, ErrNoUserEmail, apierror.NewErrorCause(ErrNoUserEmail, ErrNoUserEmailCode))
		}
		return 0, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserEmailCode))
	}

	return userId, nil
}

func (r *recoveryRepository) GetUsernameByEmail(email string) (string, apierror.ApiError) {
	var userId int64

	err := r.db.Get(&userId, "SELECT user_id FROM users_email WHERE email = ?", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apierror.New(http.StatusBadRequest, ErrNoUserEmail, apierror.NewErrorCause(ErrNoUserEmail, ErrNoUserEmailCode))
		}
		return "", apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserEmailCode))
	}

	var username string
	err = r.db.Get(&username, "SELECT username FROM users where user_id = ?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apierror.New(http.StatusBadRequest, ErrNoUserId, apierror.NewErrorCause(ErrNoUserId, ErrNoUserIdCode))
		}
		return "", apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserCode))
	}

	return username, nil
}

// GetSecurityToken Returns security token for given email
func (r *recoveryRepository) GetSecurityToken(email string) (string, apierror.ApiError) {
	var userEmail domain.UserEmail

	err := r.db.Get(&userEmail, "SELECT * FROM users_email where email = ?", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apierror.New(http.StatusBadRequest, ErrNoUserEmail, apierror.NewErrorCause(ErrNoUserEmail, ErrNoUserEmailCode))
		}
		return "", apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserEmailCode))
	}

	if !userEmail.VerfiedEmail {
		return "", apierror.New(http.StatusBadRequest, ErrEmailNotVerified, apierror.NewErrorCause(ErrEmailNotVerified, ErrEmailNotVerifiedCode))
	}

	var securityToken string

	err = r.db.Get(&securityToken, "SELECT security_token FROM users where user_id = ?", userEmail.UserId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", apierror.New(http.StatusBadRequest, ErrNoUserId, apierror.NewErrorCause(ErrNoUserId, ErrNoUserIdCode))
		}
		return "", apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserCode))
	}

	return securityToken, nil
}

// UpdatePasswordHash Updates users's password
func (r *recoveryRepository) UpdatePasswordHash(userId int64, passwordHash string) (bool, apierror.ApiError) {
	if passwordHash == "" {
		return false, apierror.New(http.StatusBadRequest, domain.ErrEmptyField, apierror.NewErrorCause(domain.ErrEmptyField, domain.ErrEmptyFieldCode))
	}

	result := r.db.MustExec("UPDATE users SET password_hash = ?  where user_id = ?", passwordHash, userId)

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return false, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrUpdatingUserCode))
	}

	if updatedRows == 0 {
		return false, apierror.New(http.StatusInternalServerError, ErrUpdatingUser, apierror.NewErrorCause("Error updating password", ErrUpdatingUserCode))
	}

	return true, nil
}

func (r *recoveryRepository) UpdateSecurityToken(userId int64, newSecurityToken string) (bool, apierror.ApiError) {
	if newSecurityToken == "" {
		return false, apierror.New(http.StatusBadRequest, domain.ErrEmptyField, apierror.NewErrorCause(domain.ErrEmptyField, domain.ErrEmptyFieldCode))
	}

	result := r.db.MustExec("UPDATE users SET security_token = ? where user_id = ?", newSecurityToken, userId)

	updatedRows, err := result.RowsAffected()
	if err != nil {
		return false, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrUpdatingUserCode))
	}

	if updatedRows == 0 {
		return false, apierror.New(http.StatusInternalServerError, ErrUpdatingUser, apierror.NewErrorCause("No rows affected", ErrUpdatingUserCode))
	}

	return true, nil
}

// GetUserByUserId Retrieves user with given user ID
func (r *recoveryRepository) GetUserByUserId(userId int64) (*domain.User, apierror.ApiError) {
	if userId == 0 {
		return nil, apierror.New(http.StatusBadRequest, domain.ErrEmptyField, apierror.NewErrorCause(domain.ErrEmptyField, domain.ErrEmptyFieldCode))
	}

	var usr domain.User
	err := r.db.Get(&usr, "SELECT * FROM users where user_id = ?", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apierror.New(http.StatusBadRequest, ErrNoUserId, apierror.NewErrorCause(ErrNoUserId, ErrNoUserIdCode))
		}
		return nil, apierror.New(http.StatusInternalServerError, domain.ErrUnexpectedError, apierror.NewErrorCause(err.Error(), ErrFetchingUserCode))
	}

	return &usr, nil
}
