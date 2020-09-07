package register

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"
	domain2 "github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type RegisterRepository interface {
	GetUserById(userId int64) (*domain2.User, error)
	AddUser(tx *sqlx.Tx, u *domain2.User) (int64, error)
	AddUserEmail(tx *sqlx.Tx, e *domain2.UserEmail) (int64, error)
	DeleteUser(userId int64) error
	CheckUsernameExists(username string) (bool, error)
	CheckEmailExists(email string) (bool, error)
}

type RegisterService interface {
	UserCanSignUp(u *domain2.UserSignupDTO) (bool, apierror.ApiError)
	CreateUser(u *domain2.UserSignupDTO, ctx *rest.ContextInformation) (int64, apierror.ApiError)
}

type RegisterController interface {
	SignUp(c *gin.Context)
}
