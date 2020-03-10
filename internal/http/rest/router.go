package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	Router       *gin.Engine
	errEmptyBody = errors.New("El cuerpo del request no puede estar vacío")
)

func InitRouter(h *healthController, ur *registerController, l *loginController, rc *recoveryController, lc *listingontroller) *gin.Engine {
	r := gin.Default()
	MapRoutes(r, h, ur, l, rc, lc)
	return r
}

func MapRoutes(r *gin.Engine, h *healthController, ur *registerController, l *loginController, rc *recoveryController, lc *listingontroller) {
	// Health
	health := r.Group("/")
	{
		health.GET("/ping", h.Ping)
	}

	user := r.Group("/users")
	{
		user.POST("/", ur.SignUp)
		user.POST("/login", l.Login)
		user.GET("/sendconfirmationemail/:userId", rc.SendConfirmationEmail)
		user.GET("/:id", lc.GetUserByUserId)
		user.GET("/confirmemail", rc.ConfirmEmail)
		user.GET("/resendconfirmationemail", rc.ResendEmailConfirmation)
		user.GET("/forgotusername", rc.ForgotUsername)
		user.GET("/sendpasswordreset", rc.SendPasswordReset)
		user.POST("/confirmpasswordreset", rc.ConfirmPasswordReset)
	}
}
