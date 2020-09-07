package login

import (
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	domain2 "github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetUserByUsername(username string) (*domain2.User, *domain2.UserEmail, apierror.ApiError)
	IncrementLoginFailAttempt(userID int64) error
	ResetLoginFails(userID int64) error
	UnlockAccount(userID int64) error
	LockAccount(userID int64, duration time.Duration) error
}

type Service interface {
	LoginUser(user *domain2.UserLoginDTO, ctx *rest.ContextInformation) (string, apierror.ApiError)
	UserCanLogin(user *domain2.UserLoginDTO) apierror.ApiError
}

type Controller interface {
	Login(c *gin.Context)
}
