package login

import (
	"github.com/go-sql-driver/mysql"
)

type UserEmail struct {
	UserEmailId      int            `json:"user_email_id" db:"user_email_id"`
	UserId           int64          `json:"user_id" db:"user_id"`
	Email            string         `json:"email" db:"email"`
	NormalizedEmail  string         `json:"normalized_email" db:"normalized_email"`
	VerfiedEmail     bool           `json:"verified_email" db:"verified_email"`
	VerificationDate mysql.NullTime `json:"verfication_date" db:"verification_date"`
	DateCreated      string         `json:"date_created" db:"date_created"`
	DateDeleted      mysql.NullTime `json:"date_deleted" db:"date_deleted"`
}
