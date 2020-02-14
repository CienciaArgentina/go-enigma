package register

import (
	"time"
)

type User struct {
	UserId              int64     `json:user_id`
	Username            string    `json:username`
	NormalizedUsername  string    `json:normalized_username`
	PasswordHash        string    `json:password_hash`
	LockoutEnabled      bool      `json:lockout_enabled`
	LockoutEnd          time.Time `json:lockout_end`
	FailedLoginAttempts int       `json:failed_login_attempts`
	DateCreated         time.Time `json:date_created`
	SecurityToken       string    `json:security_token`
	VerificationToken   string    `json:verification_token`
	DateDeleted         time.Time `json:date_deleted`
}
