package recovery

type Repository interface {
	GetEmailByUserId(userId int64) (string, *UserEmail, error)
	ConfirmUserEmail(email string, token string) error
}
