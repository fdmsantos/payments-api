package controllers

import (
	"encoding/json"
	"net/http"
	"payments/app/models"
	"payments/utils"
)

// CreateAccount handler to create new user
// Receives email and password and create a new user in accounts table
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := models.Account{}
	// Decode the request body into struct and failed if any error occur
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		utils.CreateErrorResponse(w, utils.ERROR_INVALID_JSON, http.StatusBadRequest)
		return
	}

	// Check if Email is valid
	if err := account.IsEmailValid(); err != nil {
		if err.Error() != utils.ERROR_SERVER {
			utils.CreateErrorResponse(w, err.Error(), http.StatusBadRequest)
		} else {
			utils.CreateErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	// Check if Password has 6 or more characters
	if !account.IsPasswordValid() {
		utils.CreateErrorResponse(w, utils.ERROR_PASSWORD_REQUIRED, http.StatusBadRequest)
		return
	}

	// Create Hashed password
	if err := account.CreateHashedPassword(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Create Account
	if utils.GetDB().Create(&account).Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	account.CreateToken()
	account.Password = "" // Delete password

	data, err := json.Marshal(&account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(utils.Response{
		Data: data,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)

}

// Authenticate handler to login user
// Receives email and password and return token if user was authenticated with successfully
var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	request := models.Account{}
	// Decode the request body into struct and failed if any error occur
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.CreateErrorResponse(w, utils.ERROR_INVALID_JSON, http.StatusBadRequest)
		return
	}

	// Verify if email exists
	account, err := models.GetAccountByEmail(request.Email)
	if err != nil {
		utils.CreateErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	} else if account.Email == "" {
		utils.CreateErrorResponse(w, utils.ERROR_EMAIL_NON_EXISTS, http.StatusBadRequest)
		return
	}

	if err := account.CheckPassword(request.Password); err != nil {
		if err.Error() != utils.ERROR_SERVER {
			utils.CreateErrorResponse(w, err.Error(), http.StatusUnauthorized)
		} else {
			utils.CreateErrorResponse(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	//Worked! Logged In
	account.Password = ""

	// Create JWT token
	account.CreateToken()

	message, err := json.Marshal(&account)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(utils.Response{
		Data: message,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
