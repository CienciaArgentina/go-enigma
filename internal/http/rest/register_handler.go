package rest

import (
	"errors"
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

var (
	errEmptyBody = errors.New("El cuerpo del mensaje no puede estar vacío")
)

type registerController struct {
	svc register.Service
}

func NewRegisterController(svc register.Service) *registerController {
	return &registerController{svc: svc}
}

func (r *registerController) SignUp(c *gin.Context) {
	var dto register.UserSignUp

	if err := c.ShouldBindJSON(&dto); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			err = errEmptyBody
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err.Error()))
		return
	}

	userId, errs := r.svc.SignUp(&dto)
	if errs != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errs))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, map[string]string{"userId": strconv.FormatInt(userId, 10)}, nil))
}
