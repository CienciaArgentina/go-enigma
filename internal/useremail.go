package domain

import "time"

type UserEmail struct {
	UserEmailId      int64     `json:user_email_id`
	UserId           int64     `json:user_id`
	Email            string    `json:email`
	NormalizedEmail  string    `json:normalized_email`
	VerfiedEmail     bool      `json:verified_email`
	VerificationDate time.Time `json:verfication_date`
	DateCreated      time.Time `json:date_created`
	DateDeleted      time.Time `json:date_deleted`
}

type UserEmailService interface {
	GetUserEmailByUserId(email string) (*User, error)
	GetUserIdByEmail(email string) (int64, error)
}
