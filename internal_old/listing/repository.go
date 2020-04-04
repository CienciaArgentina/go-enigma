package listing

type Repository interface {
	GetUserByUserId(id int64) (*User, error)
}
