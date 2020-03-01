package login

import "time"

type Repository interface {
	GetUserByUsername(username string) (*User, *UserEmail, error)
	IncrementLoginFailAttempt(userId int) error
	ResetLoginFails(userId int) error
	UnlockAccount(userId int) error
	LockAccount(userId int, duration time.Duration) error
	GetUserRole(userId int) (string, error)
}
