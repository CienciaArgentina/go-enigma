package login

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/gin-gonic/gin"
	"time"
)

type LoginService interface {
	LoginUser(user *domain.UserLoginDTO) (string, apierror.ApiError)
	UserCanLogin(user *domain.UserLoginDTO) apierror.ApiError
}

type LoginRepository interface {
	GetUserByUsername(username string) (*domain.User, *domain.UserEmail, apierror.ApiError)
	IncrementLoginFailAttempt(userId int64) error
	ResetLoginFails(userId int64) error
	UnlockAccount(userId int64) error
	LockAccount(userId int64, duration time.Duration) error
}

type LoginController interface {
	Login(c *gin.Context)
}
