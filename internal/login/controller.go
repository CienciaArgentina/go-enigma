package login

import (
	"net/http"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-enigma/internal/domain"
	"github.com/gin-gonic/gin"
)

type loginController struct {
	svc Service
}

func NewController(s Service) Controller {
	return &loginController{svc: s}
}

func (l *loginController) Login(c *gin.Context) {
	var usr domain.UserLoginDTO

	if err := c.ShouldBindJSON(&usr); err != nil {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, domain.ErrInvalidBody, apierror.NewErrorCause(domain.ErrInvalidBody, domain.ErrInvalidBodyCode)))
		return
	}

	jwt, errs := l.svc.LoginUser(&usr)
	if errs != nil {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
