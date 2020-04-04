package rest

import (
	"github.com/CienciaArgentina/go-enigma/internal_old/recovery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type recoveryController struct {
	svc recovery.Service
}

func NewRecoveryController(svc recovery.Service) *recoveryController {
	return &recoveryController{svc: svc}
}

func (r *recoveryController) SendConfirmationEmail(c *gin.Context) {
	userIdParam := c.Param("id")
	if userIdParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody, false))
		return
	}

	var err error
	parsedUserId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, false))
		return
	}

	sent, err := r.svc.SendConfirmationEmail(parsedUserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, sent))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, nil, nil, sent))
}

func (r *recoveryController) ConfirmEmail(c *gin.Context) {

	email := c.Query("email")
	token := c.Query("token")
	if email == "" || token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody, false))
		return
	}

	var err error
	confirmed, err := r.svc.ConfirmEmail(email, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, confirmed))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, nil, nil, confirmed))
}

func (r *recoveryController) ResendEmailConfirmation(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody, false))
		return
	}

	var err error
	sent, err := r.svc.ResendEmailConfirmationEmail(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, sent))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, nil, nil, sent))
}

func (r *recoveryController) ForgotUsername(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody, false))
		return
	}

	var err error
	sent, err := r.svc.SendUsername(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, sent))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, nil, nil, sent))
}

func (r *recoveryController) SendPasswordReset(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody, false))
		return
	}

	var err error
	sent, err := r.svc.SendPasswordReset(email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, sent))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, nil, err, sent))
}

func (r *recoveryController) ConfirmPasswordReset(c *gin.Context) {
	var dto recovery.PasswordResetDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			err = errEmptyBody
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, false))
		return
	}

	var err error
	reset, err := r.svc.ResetPassword(dto.Email, dto.Password, dto.ConfirmPassword, dto.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, reset))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, nil, err, reset))
}
