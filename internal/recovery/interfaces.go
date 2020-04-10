package recovery

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/gin-gonic/gin"
)

type RecoveryRepository interface {
	GetEmailByUserId(userId int64) (string, *domain.UserEmail, apierror.ApiError)
	ConfirmUserEmail(email string, token string) apierror.ApiError
	GetuserIdByEmail(email string) (int64, apierror.ApiError)
	GetUsernameByEmail(email string) (string, apierror.ApiError)
	GetSecurityToken(email string) (string, apierror.ApiError)
	UpdatePasswordHash(userId int64, passwordHash string) (bool, apierror.ApiError)
	UpdateSecurityToken(userId int64, newSecurityToken string) (bool, apierror.ApiError)
}

type RecoveryService interface {
	SendConfirmationEmail(userId int64) (bool, apierror.ApiError)
	ConfirmEmail(email string, token string) (bool, apierror.ApiError)
	ResendEmailConfirmationEmail(email string) (bool, apierror.ApiError)
	SendUsername(email string) (bool, apierror.ApiError)
	SendPasswordReset(email string) (bool, apierror.ApiError)
	ResetPassword(email, password, confirmPassword, token string) (bool, apierror.ApiError)
}

type RecoveryController interface {
	SendConfirmationEmail(c *gin.Context)
	ConfirmEmail(c *gin.Context)
	ResendEmailConfirmation(c *gin.Context)
	ForgotUsername(c *gin.Context)
	SendPasswordReset(c *gin.Context)
	ConfirmPasswordReset(c *gin.Context)
}
