package login

import (
	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/config"
	domain "github.com/CienciaArgentina/go-enigma/internal"
	"github.com/gin-gonic/gin"
	"net/http"
)

type loginController struct {
	svc LoginService
}

func NewController(s LoginService) LoginController {
	return &loginController{svc: s}
}

func (l *loginController) Login(c *gin.Context) {
	var usr *domain.UserLoginDTO

	if err := c.ShouldBindJSON(usr); err != nil {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, config.ErrInvalidBody, apierror.NewErrorCause(config.ErrInvalidBody, config.ErrInvalidBodyCode)))
		return
	}

	jwt, errs := l.svc.LoginUser(usr)
	if errs != nil {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
