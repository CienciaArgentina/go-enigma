package register

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/middleware"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type RegisterRepository interface {
	GetUserById(userId int64) (*domain.User, error)
	AddUser(tx *sqlx.Tx, u *domain.User) (int64, error)
	AddUserEmail(tx *sqlx.Tx, e *domain.UserEmail) (int64, error)
	DeleteUser(userId int64) error
	CheckUsernameExists(username string) (bool, error)
	CheckEmailExists(email string) (bool, error)
}

type RegisterService interface {
	UserCanSignUp(u *domain.UserSignupDTO) (bool, apierror.ApiError)
	CreateUser(u *domain.UserSignupDTO, ctx *middleware.ContextInformation) (int64, apierror.ApiError)
}

type RegisterController interface {
	SignUp(c *gin.Context)
}
