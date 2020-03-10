package rest

import (
	"github.com/CienciaArgentina/go-enigma/internal/listing"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type listingontroller struct {
	svc listing.Service
}

func NewListingController(svc listing.Service) *listingontroller {
	return &listingontroller{svc: svc}
}


func (l *listingontroller) GetUserByUserId(c *gin.Context) {
	userIdParam := c.Param("userId")
	if userIdParam == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, errEmptyBody, false))
		return
	}

	var err error
	parsedId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, false))
		return
	}

	user, err := l.svc.GetUserByUserId(parsedId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewBaseResponse(http.StatusBadRequest, nil, err, false))
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse(http.StatusOK, user, nil, user != nil))
}
