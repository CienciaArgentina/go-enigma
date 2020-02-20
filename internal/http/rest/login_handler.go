package rest

import (
	"github.com/CienciaArgentina/go-enigma/internal/login"
	"github.com/gin-gonic/gin"
	"net/http"
)

type loginController struct {
	svc login.Service
}

func NewLoginController(svc login.Service) *loginController {
	return &loginController{svc: svc}
}

func (l *loginController) Login (c *gin.Context) {
	var dto login.UserLogin

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err.Error()))
	}

	jwt, err := l.svc.Login(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err.Error()))
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, map[string]string{"jwt": jwt}, nil))
}
