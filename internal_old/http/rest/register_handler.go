package rest

import (
	"encoding/json"
	"github.com/CienciaArgentina/go-enigma/internal_old/register"
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

	//if err := c.ShouldBindJSON(&dto); err != nil {
	//	if strings.Contains(err.Error(), "EOF") {
	//		err = errEmptyBody
	//	}
	//	c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, false))
	//	return
	//}

	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &dto)

	userId, errs := r.svc.SignUp(&dto)
	if errs != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errs, false))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, map[string]string{"userId": strconv.FormatInt(userId, 10)}, nil, userId > 0))
}
