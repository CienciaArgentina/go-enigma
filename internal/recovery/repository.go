package recovery

type Repository interface {
	GetEmailByUserId(userId int64) (string, *UserEmail, error)
	ConfirmUserEmail(email string, token string) error
	GetuserIdByEmail(email string) (int64, error)
	GetUsernameByEmail(email string) (string, error)
}
