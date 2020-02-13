package repositories

import (
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/jmoiron/sqlx"
)

type registerRepository struct {
	db *sqlx.DB
}

func NewRegisterRepository(db *sqlx.DB) register.Repository {
	return &registerRepository{db:db}
}

func (r registerRepository) SignUp(u *register.User, e *register.UserEmail) (int64, error) {

}