package register

import (
	"net/http"

	"github.com/CienciaArgentina/go-enigma/internal/domain"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/gin-gonic/gin"
)

type registerController struct {
	svc RegisterService
}

func NewController(s RegisterService) RegisterController {
	return &registerController{svc: s}
}

func (u *registerController) SignUp(c *gin.Context) {
	var usr domain.UserSignupDTO

	if err := c.ShouldBindJSON(&usr); err != nil {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, domain.ErrInvalidBody, apierror.NewErrorCause(domain.ErrInvalidBody, domain.ErrInvalidBodyCode)))
		return
	}

	userId, errs := u.svc.CreateUser(&usr)
	if errs != nil {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userId})
}
