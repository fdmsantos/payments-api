package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"payments/app"
	"payments/controllers"
	_ "payments/models"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/v1/user/new", controllers.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/v1/user/login", controllers.Authenticate).Methods(http.MethodPost)
	router.HandleFunc("/v1/payments", controllers.CreatePayment).Methods(http.MethodPost)
	router.HandleFunc("/v1/payments", controllers.GetPayments).Methods(http.MethodGet)
	router.HandleFunc("/v1/payments/{id}", controllers.GetPayment).Methods(http.MethodGet)
	router.HandleFunc("/v1/payments/{id}", controllers.UpdatePayment).Methods(http.MethodPut)
	router.HandleFunc("/v1/payments/{id}", controllers.DeletePayment).Methods(http.MethodDelete)

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
