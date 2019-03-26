package controllers

import (
	"net/http"
)

var HealthCheck = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
