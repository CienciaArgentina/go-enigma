package register

type Repository interface {
	AddUser(u *User) (int64, error)
	AddEmail(e *UserEmail) (int64, error)
	VerifyIfEmailExists(email string) (bool, error)
}
