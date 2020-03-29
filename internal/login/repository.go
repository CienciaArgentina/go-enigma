package login

import "time"

type Repository interface {
	GetUserByUsername(username string) (*User, *UserEmail, error)
	IncrementLoginFailAttempt(userId int64) error
	ResetLoginFails(userId int64) error
	UnlockAccount(userId int64) error
	LockAccount(userId int64, duration time.Duration) error
	GetUserRole(userId int64) (string, error)
}
