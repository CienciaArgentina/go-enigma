package domain

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

type User struct {
	AuthId              int64          `json:"user_id" db:"user_id"`
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

type UserSignupDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PasswordResetDto struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Token           string `json:"token"`
}
