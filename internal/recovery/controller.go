package recovery

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/middleware"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/performance"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/gin-gonic/gin"
)

const (
	ErrMissingUserId     = "No puede faltar el campo de ID"
)

type recoveryController struct {
	svc RecoveryService
}

func NewController(s RecoveryService) RecoveryController {
	return &recoveryController{svc: s}
}

func (r *recoveryController) SendConfirmationEmail(c *gin.Context) {
	ctx := middleware.GetContextInformation("SendConfirmationEmail", c)
	userIdParam := c.Param("id")
	if userIdParam == "" {
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(ErrMissingUserId))
		return
	}

	var err error
	parsedUserId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(ErrMissingUserId))
		return
	}

	_, e := r.svc.SendConfirmationEmail(parsedUserId, ctx)
	if e != nil {
		c.JSON(e.Status(), e)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ConfirmEmail(c *gin.Context) {
	ctx := middleware.GetContextInformation("ConfirmEmail", c)
	email := c.Query("email")
	token := c.Query("token")
	if email == "" || token == "" {
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(domain.ErrEmptyField))
		return
	}

	var err apierror.ApiError
	performance.TrackTime(time.Now(), "ConfirmEmail", ctx, func() {
		_, err = r.svc.ConfirmEmail(email, token, ctx)
	})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ResendEmailConfirmation(c *gin.Context) {
	ctx := middleware.GetContextInformation("ResendEmailConfirmation", c)
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(domain.ErrEmptyField))
		return
	}

	var err apierror.ApiError
	performance.TrackTime(time.Now(), "ResendEmailConfirmationEmail", ctx, func() {
		_, err = r.svc.ResendEmailConfirmationEmail(email, ctx)
	})

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ForgotUsername(c *gin.Context) {
	ctx := middleware.GetContextInformation("ForgotUsername", c)
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(domain.ErrEmptyField))
		return
	}

	var err apierror.ApiError
	performance.TrackTime(time.Now(), "SendUsername", ctx, func() {
		_, err = r.svc.SendUsername(email, ctx)
	})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) SendPasswordReset(c *gin.Context) {
	ctx := middleware.GetContextInformation("ForgotUsername", c)
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(domain.ErrEmptyField))
		return
	}

	var err apierror.ApiError
	performance.TrackTime(time.Now(), "SendPasswordReset", ctx, func() {
		_, err = r.svc.SendPasswordReset(email, ctx)
	})

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) ConfirmPasswordReset(c *gin.Context) {
	ctx := middleware.GetContextInformation("ConfirmPasswordReset", c)
	var dto domain.PasswordResetDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(domain.ErrEmptyField))
			return
		}
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(err.Error()))
		return
	}

	var err apierror.ApiError
	performance.TrackTime(time.Now(), "ResetPassword", ctx, func() {
		_, err = r.svc.ResetPassword(dto.Email, dto.Password, dto.ConfirmPassword, dto.Token, ctx)
	})
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.Status(http.StatusOK)
}

func (r *recoveryController) GetUserByUserId(c *gin.Context) {
	id := c.Param("id")

	userid, err := strconv.Atoi(id)
	if err != nil || id == "" {
		c.JSON(http.StatusBadRequest, apierror.NewBadRequestApiError(domain.ErrEmptyField))
		return
	}

	usr, e := r.svc.GetUserByUserId(int64(userid))
	if e != nil {
		c.JSON(e.Status(), e)
		return
	}

	c.JSON(http.StatusOK, usr)
}
