package register

import (
	"net/http"
	"time"

	"github.com/CienciaArgentina/go-backend-commons/pkg/performance"
	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"

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
	var errs apierror.ApiError
	ctx := rest.GetContextInformation("SignUp", c)

	if err := c.ShouldBindJSON(&usr); err != nil {
		c.JSON(http.StatusBadRequest, apierror.New(http.StatusBadRequest, domain.ErrInvalidBody, apierror.NewErrorCause(domain.ErrInvalidBody, domain.ErrInvalidBodyCode)))
		return
	}

	var userId int64
	performance.TrackTime(time.Now(), "CreateUser", ctx, func() {
		userId, errs = u.svc.CreateUser(&usr, ctx)
	})

	if errs != nil {
		c.JSON(errs.Status(), errs)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userId})
}
