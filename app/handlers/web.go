package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"payments/app/controllers"
)

// Handlers
var Routes = func(router *mux.Router) {
	router.HandleFunc("/v1/user", controllers.CreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/v1/user/login", controllers.Authenticate).Methods(http.MethodPost)
	router.HandleFunc("/v1/payments", controllers.CreatePayment).Methods(http.MethodPost)
	router.HandleFunc("/v1/payments", controllers.GetPayments).Methods(http.MethodGet)
	router.HandleFunc("/v1/payments/{id}", controllers.GetPayment).Methods(http.MethodGet)
	router.HandleFunc("/v1/payments/{id}", controllers.UpdatePayment).Methods(http.MethodPut)
	router.HandleFunc("/v1/payments/{id}", controllers.DeletePayment).Methods(http.MethodDelete)
	router.HandleFunc("/v1/health", controllers.HealthCheck).Methods(http.MethodGet)
}
