package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"payments/app/handlers"
	"payments/app/middleware"
	"payments/app/models"
	"payments/infrastructure"
)

func main() {
	router := mux.NewRouter()
	handlers.Routes(router)
	router.Use(middleware.JwtAuthentication)

	provisionDatabase()

	err := http.ListenAndServe(":8000", router) //Launch the app
	if err != nil {
		fmt.Print(err)
	}
}

// provisionDatabase Create tables on Database
func provisionDatabase() {
	infrastructure.GetDB().AutoMigrate(
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
}
