package recovery

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"
	domain2 "github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/gin-gonic/gin"
)

type RecoveryRepository interface {
	GetEmailByUserId(userId int64) (string, *domain2.UserEmail, apierror.ApiError)
	ConfirmUserEmail(email string, token string) apierror.ApiError
	GetuserIdByEmail(email string) (int64, apierror.ApiError)
	GetUsernameByEmail(email string) (string, apierror.ApiError)
	GetSecurityToken(email string) (string, apierror.ApiError)
	UpdatePasswordHash(userId int64, passwordHash string) (bool, apierror.ApiError)
	UpdateSecurityToken(userId int64, newSecurityToken string) (bool, apierror.ApiError)
	GetUserByUserId(userId int64) (*domain2.User, apierror.ApiError)
}

type RecoveryService interface {
	SendConfirmationEmail(userId int64, ctx *rest.ContextInformation) (bool, apierror.ApiError)
	ConfirmEmail(email string, token string, ctx *rest.ContextInformation) (bool, apierror.ApiError)
	ResendEmailConfirmationEmail(email string, ctx *rest.ContextInformation) (bool, apierror.ApiError)
	SendUsername(email string, ctx *rest.ContextInformation) (bool, apierror.ApiError)
	SendPasswordReset(email string, ctx *rest.ContextInformation) (bool, apierror.ApiError)
	ResetPassword(email, password, confirmPassword, token string, ctx *rest.ContextInformation) (bool, apierror.ApiError)
	GetUserByUserId(userId int64) (*domain2.User, apierror.ApiError)
}

type RecoveryController interface {
	SendConfirmationEmail(c *gin.Context)
	ConfirmEmail(c *gin.Context)
	ResendEmailConfirmation(c *gin.Context)
	ForgotUsername(c *gin.Context)
	SendPasswordReset(c *gin.Context)
	ConfirmPasswordReset(c *gin.Context)
	GetUserByUserId(c *gin.Context)
}
