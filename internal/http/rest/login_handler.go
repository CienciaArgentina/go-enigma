package rest

import (
	"fmt"
	"github.com/CienciaArgentina/go-enigma/internal/login"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

type loginController struct {
	svc login.Service
}

func NewLoginController(svc login.Service) *loginController {
	return &loginController{svc: svc}
}

func (l *loginController) Login(c *gin.Context) {
	var dto login.UserLogin

	logrus.SetLevel(logrus.InfoLevel)
	logrus.Info("Iniciando request de Login")
	start := time.Now()

	if err := c.ShouldBindJSON(&dto); err != nil {
		if strings.Contains(err.Error(), "EOF") {
			err = errEmptyBody
		}
		c.JSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, false))
		return
	}

	jwt, err := l.svc.Login(&dto)
	if err != nil {
		elapsed := time.Since(start)
		logrus.WithField("elapsed", fmt.Sprintf("%dms", elapsed.Milliseconds())).Info("Terminó request de login")
		c.JSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, false))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, map[string]string{"jwt": jwt}, nil, jwt != "" && len(jwt) > 0))
	return
}
