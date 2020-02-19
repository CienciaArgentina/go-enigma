package login

type Repository interface {
	GetUserByUsername(username string) (*User, *UserEmail)
	IncrementLoginFailAttempt(userId int) error
	ResetLoginFails(userId int) error
	UnlockAccount(userId int) error
	LockAccount(userId int) error
	GetUserRole(userId int) (string, error)
}
