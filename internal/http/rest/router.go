package rest

import (
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/login"
	"github.com/CienciaArgentina/go-enigma/internal/recovery"
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/CienciaArgentina/go-enigma/internal_old/storage/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func InitRouter(cfg *config.Configuration) *gin.Engine {
	router := gin.Default()
	MapRoutes(router, cfg)
	return router
}

func MapRoutes(r *gin.Engine, cfg *config.Configuration) {
	db := database.New(cfg)

	registerRepo := register.NewRepository(db)
	registerSvc := register.NewService(cfg, db, nil, registerRepo)
	registerCtrl := register.NewController(registerSvc)

	loginRepo := login.NewRepository(db)
	loginSvc := login.NewService(cfg, nil, loginRepo)
	loginCtrl := login.NewController(loginSvc)

	recoveryRepo := recovery.NewRepository(db)
	recoverySvc := recovery.NewService(cfg, recoveryRepo)
	recoveryCtrl := recovery.NewController(recoverySvc)

	user := r.Group("/users")
	{
		user.POST("/", registerCtrl.SignUp)
		user.POST("/login", loginCtrl.Login)
		user.POST("/confirmpasswordreset", recoveryCtrl.ConfirmPasswordReset)
		user.GET("/:id", func(c *gin.Context) {
			GetHandler(c, recoveryCtrl)
		})
	}
}

// I have to do this just because gin works like shit
func GetHandler(c *gin.Context, rc recovery.RecoveryController) {
	id := c.Param("id")

	if strings.HasPrefix(c.Request.RequestURI, "/sendconfirmationemail") {
		// /users/sendconfirmationemail
		rc.SendConfirmationEmail(c)
	} else if _, err := strconv.Atoi(id); err == nil {
		// /users/1
		rc.GetUserByUserId(c)
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
		Ping(c)
	}
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}