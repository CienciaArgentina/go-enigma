package domain

import (
	"database/sql"
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

type User struct {
	UserId              int64          `json:"user_id" db:"user_id"`
	Username            string         `json:"username" db:"username"`
	NormalizedUsername  string         `json:"normalized_username" db:"normalized_username"`
	PasswordHash        string         `json:"password_hash" db:"password_hash"`
	LockoutEnabled      bool           `json:"lockout_enabled" db:"lockout_enabled"`
	LockoutDate         mysql.NullTime `json:"lockout_date" db:"lockout_date"`
	FailedLoginAttempts int            `json:"failed_login_attempts" db:"failed_login_attempts"`
	DateCreated         string         `json:"date_created" db:"date_created"`
	SecurityToken       sql.NullString `json:"security_token" db:"security_token"`
	VerificationToken   string         `json:"verification_token" db:"verification_token"`
	DateDeleted         *time.Time     `json:"date_deleted" db:"date_deleted"`
}

type UserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserService interface {
	UserCanSignUp(u *UserDTO) (bool, apierror.ApiError)
	CreateUser(u *UserDTO) (int64, apierror.ApiError)
}

type UserRepository interface {
	GetUserById(userId int64) (*User, error)
	AddUser(tx *sqlx.Tx, u *User) (int64, error)
	AddUserEmail(tx *sqlx.Tx, e *UserEmail) (int64, error)
	DeleteUser(userId int64) error
	CheckUsernameExists(username string) (bool, error)
	CheckEmailExists(email string) (bool, error)
}

type UserController interface {
	SignUp(c *gin.Context)
}
