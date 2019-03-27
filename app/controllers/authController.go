package controllers

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"payments/app/models"
	u "payments/utils"
	"strings"
	"time"
)

const ERROR_EMAIL_REQUIRED = "Email address is required"
const ERROR_PASSWORD_REQUIRED = "Password is required"
const ERROR_EMAIL_EXISTS = "Email address already in use by another user"
const ERROR_EMAIL_NON_EXISTS = "Email address not found"
const ERROR_INVALID_LOGIN = "Invalid login credentials. Please try again"

// CreateAccount handler to create new user
// Receives email and password and create a new user in accounts table
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	// Decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		if response, err := json.Marshal(u.Response{Errors: []string{u.ERROR_INVALID_JSON}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
		}
		return
	}

	// Check if Email contains @ character
	if !strings.Contains(account.Email, "@") {
		if response, err := json.Marshal(u.Response{Errors: []string{ERROR_EMAIL_REQUIRED}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
		}
		return
	}

	// Check if Password has 6 or more characters
	if len(account.Password) < 6 {
		if response, err := json.Marshal(u.Response{Errors: []string{ERROR_PASSWORD_REQUIRED}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
		}
		return
	}

	// Email must be unique
	temp := &models.Account{}

	// Check for errors and duplicate emails
	err = u.GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if temp.Email != "" {
		if response, err := json.Marshal(u.Response{Errors: []string{ERROR_EMAIL_EXISTS}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
		}
		return
	}

	// Create Hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	// Create Account
	if u.GetDB().Create(account).Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create new JWT token for the newly registered account
	tk := &models.Token{
		UserId: account.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		}}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" // Delete password

	data, err := json.Marshal(account)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(u.Response{
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

	request := &models.Account{}
	// Decode the request body into struct and failed if any error occur
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		if response, err := json.Marshal(u.Response{Errors: []string{u.ERROR_INVALID_JSON}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
		}
		return
	}

	account := &models.Account{}

	// Verify if email exists
	if err := u.GetDB().Table("accounts").Where("email = ?", request.Email).First(account).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if response, err := json.Marshal(u.Response{Errors: []string{ERROR_EMAIL_NON_EXISTS}}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(response)
			}
		}
		return
	}

	// Check if the password is valid
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		if response, err := json.Marshal(u.Response{Errors: []string{ERROR_INVALID_LOGIN}}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(response)
		}
		return
	}

	//Worked! Logged In
	account.Password = ""

	// Create JWT token
	tk := &models.Token{
		UserId: account.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		}}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response

	message, err := json.Marshal(account)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(u.Response{
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
