package register

import (
	"time"
)

type User struct {
	Username           string    `json:"username"`
	NormalizedUsername string    `json:"normalized_username"`
	PasswordHash       string    `json:"password_hash"`
	DateCreated        time.Time `json:"date_created"`
	VerificationToken  string    `json:"verification_token"`
	SecurityToken      string    `json:"security_token"`
}
