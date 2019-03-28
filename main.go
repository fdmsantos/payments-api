package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"payments/app/handlers"
	"payments/app/middleware"
	"payments/app/models"
	"payments/utils"
)

func main() {
	router := mux.NewRouter()
	handlers.Routes(router)
	router.Use(middleware.JwtAuthentication)
	//router.NotFoundHandler = utils.NotFoundHandler

	createDB()

	err := http.ListenAndServe(":8000", router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}

// createDB Creates the tables on database
func createDB() {
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
}
