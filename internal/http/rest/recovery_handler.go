package rest

import (
	"github.com/CienciaArgentina/go-enigma/internal/recovery"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type recoveryController struct {
	svc recovery.Service
}

func NewRecoveryController(svc recovery.Service) *recoveryController {
	return &recoveryController{svc: svc}
}

func (r *recoveryController) SendConfirmationEmail(c *gin.Context) {
	userIdParam := c.Param("userId")
	if userIdParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody))
		return
	}

	var err error
	parsedUserId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err))
		return
	}

	sent, err := r.svc.SendConfirmationEmail(parsedUserId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, map[string]string{"sentEmail": strconv.FormatBool(sent)}, nil))
}

func (r *recoveryController) ConfirmEmail(c *gin.Context) {

	email := c.Query("email")
	token := c.Query("token")
	if email == "" || token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody))
		return
	}

	var err error
	confirmed, err := r.svc.ConfirmEmail(email, token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, map[string]string{"confirmed": strconv.FormatBool(confirmed)}, err))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, map[string]string{"confirmed": strconv.FormatBool(confirmed)}, nil))
}

