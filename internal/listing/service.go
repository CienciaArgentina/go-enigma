package listing

type Service interface {
	GetUserByUserId(id int64) (*User, error)
}

type listingService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &listingService{repo:r}
}

func (l *listingService) GetUserByUserId(id int64) (*User, error) {
	user, err := l.repo.GetUserByUserId(id)
	if err != nil {
		return nil, err
	}

	return user, err
}