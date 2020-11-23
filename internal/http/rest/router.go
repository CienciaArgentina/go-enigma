package rest

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/CienciaArgentina/go-backend-commons/pkg/rest"

	config2 "github.com/CienciaArgentina/go-backend-commons/config"
	"github.com/CienciaArgentina/go-backend-commons/pkg/clog"
	"github.com/CienciaArgentina/go-backend-commons/pkg/injector"
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/login"
	"github.com/CienciaArgentina/go-enigma/internal/recovery"
	"github.com/CienciaArgentina/go-enigma/internal/register"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(rest.SetContextInformation)
	MapRoutes(router)
	return router
}

func MapRoutes(r *gin.Engine) {
	injector.Initilize()

	dbname := os.Getenv(config2.EnvDBName)
	db := injector.GetDB(dbname).Database

	enigmaConfig, err := config.NewEnigmaConfig()
	if err != nil {
		msg := "error building enigma config"
		clog.Panic(msg, "map-routes", err, nil)
		return
	}

	loginRepo := login.NewRepository(db)
	loginSvc := login.NewService(enigmaConfig, loginRepo)
	loginCtrl := login.NewController(loginSvc)

	recoveryRepo := recovery.NewRepository(db)
	recoverySvc := recovery.NewService(enigmaConfig, recoveryRepo)
	recoveryCtrl := recovery.NewController(recoverySvc)

	registerRepo := register.NewRepository(db)
	registerSvc := register.NewService(enigmaConfig, db, registerRepo, recoverySvc)
	registerCtrl := register.NewController(registerSvc)

	r.GET("/ping", Ping)

	user := r.Group("/users")
	{
		user.POST("/", registerCtrl.SignUp)
		user.POST("/login", loginCtrl.Login)
		user.POST("/confirm_password_reset", recoveryCtrl.ConfirmPasswordReset)
		user.GET("/:id", func(c *gin.Context) {
			GetHandler(c, recoveryCtrl)
		})
	}
}

// I have to do this just because gin router can't handle REST standards.
func GetHandler(c *gin.Context, rc recovery.RecoveryController) {
	id := c.Param("id")

	if strings.Contains(c.Request.RequestURI, "send_confirmation_email") {
		// /users/sendconfirmationemail
		rc.SendConfirmationEmail(c)
	} else if _, err := strconv.Atoi(id); err == nil {
		// /users/1
		rc.GetUserByUserId(c)
	} else if strings.Contains(c.Request.RequestURI, "confirm_email") {
		// /users/confirmemail
		rc.ConfirmEmail(c)
	} else if strings.Contains(c.Request.RequestURI, "resend_confirmation_email") {
		// /users/resendconfirmationemail
		rc.ResendEmailConfirmation(c)
	} else if strings.Contains(c.Request.RequestURI, "forgot_username") {
		// /users/forgotusername
		rc.ForgotUsername(c)
	} else if strings.Contains(c.Request.RequestURI, "send_password_reset") {
		// /users/sendpasswordreset
		rc.SendPasswordReset(c)
	}
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}
