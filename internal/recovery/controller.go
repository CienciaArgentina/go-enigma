package recovery

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/config"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	ErrMissingUserId = "No puede faltar el campo de ID"
	ErrMissingUserIdCode = "missing_user_id"
)

type recoveryController struct {
	svc RecoveryService
}

func NewController(s RecoveryService) RecoveryController {
	return &recoveryController{svc:s}
}

func (r *recoveryController) SendConfirmationEmail(c *gin.Context) {
	userIdParam := c.Param("id")
	if userIdParam == "" {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, ErrMissingUserId, apierror.NewErrorCause(ErrMissingUserId, ErrMissingUserIdCode)))
		return
	}

	var err error
	parsedUserId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, ErrMissingUserId, apierror.NewErrorCause(err.Error(), ErrMissingUserIdCode)))
		return
	}

	_, e := r.svc.SendConfirmationEmail(parsedUserId)
	if e != nil {
		c.JSON(e.Status(), e)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ConfirmEmail(c *gin.Context) {
	email := c.Query("email")
	token := c.Query("token")
	if email == "" || token == "" {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, config.ErrEmptyField, apierror.NewErrorCause(config.ErrEmptyField, config.ErrEmptyFieldCode)))
		return
	}

	_, err := r.svc.ConfirmEmail(email, token)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ResendEmailConfirmation(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, config.ErrEmptyField, apierror.NewErrorCause(config.ErrEmptyField, config.ErrEmptyFieldCode)))
		return
	}

	_, err := r.svc.ResendEmailConfirmationEmail(email)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ForgotUsername(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, config.ErrEmptyField, apierror.NewErrorCause(config.ErrEmptyField, config.ErrEmptyFieldCode)))
		return
	}

	_, err := r.svc.SendUsername(email)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) SendPasswordReset(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, config.ErrEmptyField, apierror.NewErrorCause(config.ErrEmptyField, config.ErrEmptyFieldCode)))
		return
	}

	_, err := r.svc.SendPasswordReset(email)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ConfirmPasswordReset(c *gin.Context) {
	var dto domain.PasswordResetDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, config.ErrEmptyField, apierror.NewErrorCause(config.ErrEmptyField, config.ErrEmptyFieldCode)))
			return
		}
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, err.Error(), apierror.NewErrorCause(err.Error(), config.ErrEmptyFieldCode)))
		return
	}

	_, err := r.svc.ResetPassword(dto.Email, dto.Password, dto.ConfirmPassword, dto.Token)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) GetUserByUserId(c *gin.Context) {
	id := c.Param("id")

	userid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, config.ErrEmptyField, apierror.NewErrorCause(config.ErrEmptyField, config.ErrEmptyFieldCode)))
	}

	usr, e := r.svc.GetUserByUserId(int64(userid))
	if e != nil {
		 c.JSON(e.Status(), e)
		 return
	}

	c.JSON(http.StatusOK, usr)
}

