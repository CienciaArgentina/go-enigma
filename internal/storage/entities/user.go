package entities

import "time"

type User struct {
	UserId                   string `db:"user_id"`
	Username                 string `db:"username"`
	NormalizedUsername       string `db:"normalized_username"`
	PasswordHash             string `db:"password_hash"`
	LockoutEnabled           bool `db:"lockout_enabled"`
	LockoutEndDate           time.Time `db:"lockout_end_date"`
	FailedLoginAttempts      int `db:"failed_login_attempts"`
	DateCreated              time.Time `db:date_created`
	DateModified             time.Time `db:"date_modified"`
	ModificationMadeByUserId int `db:"modification_made_by_user_id"`
	SecurityToken            string `db:"security_token"`
	VerificationToken        string `db:"verification_token"`
	DateDeleted              time.Time `db:"date_deleted"`
}
