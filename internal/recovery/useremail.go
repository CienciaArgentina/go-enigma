package recovery

import (
	"database/sql"
	"time"
)

type UserEmail struct {
	UserEmailId      int          `json:"user_email_id" db:"user_email_id"`
	UserId           int64        `json:"user_id" db:"user_id"`
	Email            string       `json:"email" db:"email"`
	NormalizedEmail  string       `json:"normalized_email" db:"normalized_email"`
	VerfiedEmail     bool         `json:"verified_email" db:"verified_email"`
	VerificationDate *time.Time   `json:"verfication_date" db:"verification_date"`
	DateCreated      string       `json:"date_created" db:"date_created"`
	DateDeleted      sql.NullTime `json:"date_deleted" db:"date_deleted"`
}
