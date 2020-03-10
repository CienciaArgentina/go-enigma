package repositories

import (
	"github.com/CienciaArgentina/go-enigma/internal/listing"
	"github.com/jmoiron/sqlx"
)

type listingRepository struct {
	db *sqlx.DB
}

func NewListingRepository(db *sqlx.DB) listing.Repository {
	return &listingRepository{db: db}
}

func (l *listingRepository) GetUserByUserId(id int64) (*listing.User, error) {
	var u *listing.User

	var email string
	err := l.db.Get(&email, "SELECT email FROM users_email WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}

	var username string
	err = l.db.Get(&username, "SELECT username FROM users WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}

	u.Username = username
	u.Email = email

	return u, nil
}

