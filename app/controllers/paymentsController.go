package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"payments/app/models"
	"payments/infrastructure"
	"payments/utils"
)

// CreatePayment handler to create a single payment
// Receives the payment and inserts in database
var CreatePayment = func(w http.ResponseWriter, r *http.Request) {

	infrastructure.LogApiRequest(r)

	//user := r.Context().Value("user") . (uint) //Grab the id of the user that send the request

	// Decode the request body into payment struct and failed if any error occur
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		utils.CreateApiErrorResponse(w, utils.ERROR_INVALID_JSON, http.StatusBadRequest)
		return
	}

	// Verify if the requested payment already exists in DB
	if _, err := models.GetPaymentByID(payment.ID); err == nil || (err != nil && err.Error() != utils.ERROR_RESOURCE_NOT_FOUND) {
		utils.CreateApiErrorResponse(w, utils.ERROR_PAYMENT_ALREADY_EXISTS, http.StatusBadRequest)
		return
	}

	// Creates the payment in DB
	if infrastructure.GetDB().Create(&payment).Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create Api Response
	links := []utils.Link{{
		Rel:  "self",
		Href: fmt.Sprintf("/v1/payments/%s", payment.ID.String()),
	}}
	utils.CreateApiResponse(w, nil, http.StatusCreated, links)
}

// GetPayments handler to get all payments
// Returns all payments from database
var GetPayments = func(w http.ResponseWriter, r *http.Request) {

	infrastructure.LogApiRequest(r)

	var payments []models.Payment

	// Fetch all payments from DB
	if err := infrastructure.GetDB().Set("gorm:auto_preload", true).Find(&payments).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create Links
	links := []utils.Link{{
		Rel:  "self",
		Href: "/v1/payments",
	}}
	for _, payment := range payments {
		links = append(links, utils.Link{
			Rel:  payment.ID.String(),
			Href: fmt.Sprintf("/v1/payments/%s", payment.ID.String()),
		})
	}

	// Create Api Response
	utils.CreateApiResponse(w, payments, http.StatusOK, links)
}

// GetPayment handler to get a single payment
// Receives the payment id and returns the payment
var GetPayment = func(w http.ResponseWriter, r *http.Request) {

	infrastructure.LogApiRequest(r)

	// Read the ID from the mux vars
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok { // this should not be possible as muxer will only route requests with an ID
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse the UUID
	uuid, err := utils.ConvertStringToUUID(id)
	if err != nil {
		utils.CreateApiErrorResponse(w, utils.ERROR_REQUESTED_UUID_INVALID, http.StatusBadRequest)
		return
	}

	// Fetch the requested payment from the db
	payment, err := models.GetPaymentByID(uuid)
	if err != nil {
		if err.Error() != utils.ERROR_SERVER {
			utils.CreateApiErrorResponse(w, err.Error(), http.StatusNotFound)
		} else {
			utils.CreateApiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create Api Response
	links := []utils.Link{{
		Rel:  "self",
		Href: fmt.Sprintf("/v1/payments/%s", payment.ID.String()),
	}}

	utils.CreateApiResponse(w, payment, http.StatusOK, links)
}

// UpdatePayment handler update a single payment
// Receives payment id and updates the payment in database
var UpdatePayment = func(w http.ResponseWriter, r *http.Request) {

	infrastructure.LogApiRequest(r)

	// Read the ID from the mux vars
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok { // the muxer should not assign this handler if the id is missing, so internal error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse the UUID
	uuid, err := utils.ConvertStringToUUID(id)
	if err != nil {
		utils.CreateApiErrorResponse(w, utils.ERROR_REQUESTED_UUID_INVALID, http.StatusBadRequest)
		return
	}

	// Decode the request body into payment struct and failed if any error occur
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		utils.CreateApiErrorResponse(w, utils.ERROR_INVALID_JSON, http.StatusBadRequest)
		return
	}

	// Ensure the payment being updated matches the one specified in the URL
	if payment.ID.String() != uuid.String() {
		utils.CreateApiErrorResponse(w, utils.ERROR_ID_MISMATCH, http.StatusBadRequest)
		return
	}

	// Verify if the payment exists before editing/replacing it
	oldPayment, err := models.GetPaymentByID(uuid)
	if err != nil {
		if err.Error() != utils.ERROR_SERVER {
			utils.CreateApiErrorResponse(w, err.Error(), http.StatusNotFound)
		} else {
			utils.CreateApiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	oldPayment = payment
	// Update the payment in DB
	if err := infrastructure.GetDB().Save(&oldPayment).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create Api Response
	links := []utils.Link{{
		Rel:  "self",
		Href: fmt.Sprintf("/v1/payments/%s", payment.ID.String()),
	}}
	utils.CreateApiResponse(w, nil, http.StatusOK, links)
}

// DeletePayment handler to delete a single payment
// Receives the payment id and deletes the payment in database
var DeletePayment = func(w http.ResponseWriter, r *http.Request) {

	infrastructure.LogApiRequest(r)

	// Read the ID from the mux vars
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok { // the muxer should not assign this handler if the id is missing, so internal error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Parse the UUID
	uuid, err := utils.ConvertStringToUUID(id)
	if err != nil {
		utils.CreateApiErrorResponse(w, utils.ERROR_REQUESTED_UUID_INVALID, http.StatusBadRequest)
		return
	}

	// Verify if the payment exists before attempting to delete it

	payment, err := models.GetPaymentByID(uuid)
	if err != nil {
		if err.Error() != utils.ERROR_SERVER {
			utils.CreateApiErrorResponse(w, err.Error(), http.StatusNotFound)
		} else {
			utils.CreateApiErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Delete the payment
	if err := infrastructure.GetDB().Delete(&payment).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create Api Response
	utils.CreateApiResponse(w, nil, http.StatusNoContent, nil)
}
