package register

import "time"

type UserEmail struct {
	UserEmailId      string    `json:user_email_id`
	UserId           string    `json:user_id`
	Email            string    `json:email`
	NormalizedEmail  string    `json:normalized_email`
	VerfiedEmail     bool      `json:verified_email`
	VerificationDate time.Time `json:verfication_date`
	DateCreated      time.Time `json:date_created`
	DateDeleted      time.Time `json:date_deleted`
}
