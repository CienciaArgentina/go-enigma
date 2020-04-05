package main

import (
	"github.com/CienciaArgentina/go-enigma/config"
	"github.com/CienciaArgentina/go-enigma/internal/http/rest"
	"github.com/sirupsen/logrus"
)

func main() {

	logrus.SetLevel(logrus.InfoLevel)

	cfg := config.New()
	//h := rest.NewHealthController()

	//regRepo := repositories.NewRegisterRepository(db)
	//regSvc := register.NewService(regRepo, nil, cfg)
	//regCtrl := rest.NewRegisterController(regSvc)
	//
	//logRepo := repositories.NewLoginRepository(db)
	//logsvc := login.NewService(logRepo, nil, cfg)
	//logCtrl := rest.NewLoginController(logsvc)
	//
	//recRepo := repositories.NewRecoveryRepository(db)
	//recSvc := recovery.NewService(recRepo, cfg)
	//recCtlr := rest.NewRecoveryController(recSvc)
	//
	//lisRepo := repositories.NewListingRepository(db)
	//lisSvc := listing.NewService(lisRepo)
	//lisCtlr := rest.NewListingController(lisSvc)

	//if err := rest.InitRouter(h, regCtrl, logCtrl, recCtlr, lisCtlr).Run(cfg.Server.Port); err != nil {
	//	panic(err)
	//}

	if err := rest.InitRouter(cfg).Run(cfg.Server.Port); err != nil {
		panic(err)
	}
}
