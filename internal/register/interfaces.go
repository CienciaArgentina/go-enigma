package register

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type RegisterService interface {
	UserCanSignUp(u *domain.UserSignupDTO) (bool, apierror.ApiError)
	CreateUser(u *domain.UserSignupDTO) (int64, apierror.ApiError)
}

type RegisterRepository interface {
	GetUserById(userId int64) (*domain.User, error)
	AddUser(tx *sqlx.Tx, u *domain.User) (int64, error)
	AddUserEmail(tx *sqlx.Tx, e *domain.UserEmail) (int64, error)
	DeleteUser(userId int64) error
	CheckUsernameExists(username string) (bool, error)
	CheckEmailExists(email string) (bool, error)
}

type RegisterController interface {
	SignUp(c *gin.Context)
}
