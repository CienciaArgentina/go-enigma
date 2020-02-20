package rest

import (
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err.Error()))
	}

	userId, errs := r.svc.SignUp(&dto)
	if errs != nil {
		c.JSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errs))
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, map[string]string{"userId": strconv.FormatInt(userId, 10)}, nil))
}
