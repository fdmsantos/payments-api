package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"payments/app/handlers"
	"payments/app/middleware"
	"payments/models"
	"payments/utils"
)

func main() {

	router := mux.NewRouter()
	handlers.Routes(router)
	router.Use(middleware.JwtAuthentication)
	//router.NotFoundHandler = utils.NotFoundHandler

	utils.GetDB().AutoMigrate(
		&models.Account{},
		&models.Payment{},
		&models.Attributes{},
		&models.BeneficiaryParty{},
		&models.DebtorParty{},
		&models.SponsorParty{},
		&models.ChargesInformation{},
		&models.Charge{},
		&models.FX{},
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
