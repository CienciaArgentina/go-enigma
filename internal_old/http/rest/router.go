package rest

import (
	"errors"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
	"strings"
)

var (
	Router       *gin.Engine
	errEmptyBody = errors.New("El cuerpo del request no puede estar vacío")
)

func InitRouter(h *healthController, ur *registerController, l *loginController, rc *recoveryController, lc *listingontroller) *gin.Engine {
	r := gin.Default()

	if os.Getenv(config.Scope) == config.Development {
		c := cors.Config{
			AllowAllOrigins: true,
			AllowMethods:    []string{"POST", "GET", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:    []string{"*"},
		}
		r.Use(cors.New(c))
	}
	MapRoutes(r, h, ur, l, rc, lc)
	return r
}

func MapRoutes(r *gin.Engine, h *healthController, ur *registerController, l *loginController, rc *recoveryController, lc *listingontroller) {
	user := r.Group("/users")
	{
		user.POST("/", ur.SignUp)
		user.POST("/login", l.Login)
		user.POST("/confirmpasswordreset", rc.ConfirmPasswordReset)
		user.GET("/:id", func(c *gin.Context) {
			GetHandler(c, h, rc, lc)
		})
	}
}

// I have to do this just because gin works like shit
func GetHandler(c *gin.Context, h *healthController, rc *recoveryController, lc *listingontroller) {
	id := c.Param("id")

	if strings.HasPrefix(c.Request.RequestURI, "/sendconfirmationemail") {
		// /users/sendconfirmationemail
		rc.SendConfirmationEmail(c)
	} else if _, err := strconv.Atoi(id); err == nil {
		// /users/1
		lc.GetUserByUserId(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/confirmemail") {
		// /users/confirmemail
		rc.ConfirmEmail(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/resendconfirmationemail") {
		// /users/resendconfirmationemail
		rc.ResendEmailConfirmation(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/forgotusername") {
		// /users/forgotusername
		rc.ForgotUsername(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/sendpasswordreset") {
		// /users/sendpasswordreset
		rc.SendPasswordReset(c)
	} else if strings.HasPrefix(c.Request.RequestURI, "/ping") {
		// /users/ping
		h.Ping(c)
	}
}
