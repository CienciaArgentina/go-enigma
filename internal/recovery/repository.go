package recovery

type Repository interface {
	GetEmailByUserId(userId int64) (string, *UserEmail, error)
	ConfirmUserEmail(email string, token string) error
	GetuserIdByEmail(email string) (int64, error)
	GetUsernameByEmail(email string) (string, error)
	GetSecurityToken(email string) (string, error)
	UpdatePasswordAndResetSecurityToken(userId int64, passwordHash, newSecurityToken string) (bool, error)
}
