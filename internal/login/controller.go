package login

import (
	"net/http"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"

	"github.com/CienciaArgentina/go-backend-commons/pkg/apierror"
	"github.com/CienciaArgentina/go-backend-commons/pkg/performance"
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
	ctx := rest.GetContextInformation("login", c)

	if err := c.ShouldBindJSON(&usr); err != nil {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, domain.ErrInvalidBody, apierror.NewErrorCause(domain.ErrInvalidBody, domain.ErrInvalidBodyCode)))
		return
	}

	var jwt string
	var apierr apierror.ApiError
	performance.TrackTime(time.Now(), "CompleteLogin", ctx, func() {
		jwt, apierr = l.svc.LoginUser(&usr, ctx)
	})
	if apierr != nil {
		c.JSON(apierr.Status(), apierr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"jwt": jwt})
}
