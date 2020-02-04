package entities

import "time"

type UsersEmails struct {
	UserEmailId              int       `db:"user_email_id"`
	UserId                   int       `db:"user_id"`
	Email                    string    `db:"email"`
	NormalizedEmail          string    `db:"normalized_email"`
	VerifiedEmail            bool      `db:"verified_email"`
	VerificationDate         time.Time `db:"verification_date"`
	DateCreated              time.Time `db:"date_created"`
	ModificationDate         time.Time `db:"modification_time"`
	ModificationMadeByUserId int       `db:"modification_made_by_user_id"`
	DateDeleted              time.Time `db:"date_deleted"`
}
