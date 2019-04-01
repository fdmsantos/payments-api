package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"payments/app/models"
	"payments/infrastructure"
	"payments/utils"
	"testing"
)

func TestCreateAccountWithInvalidBody(t *testing.T) {

	jsonBytes := []byte("{ malformed json }")

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_INVALID_JSON}, response.Errors)
}

func TestCreateAccountWithInvalidEmail(t *testing.T) {

	account := models.Account{
		Email:    "dummyemail",
		Password: "dummy",
	}

	jsonBytes, err := json.Marshal(account)

	if err != nil {
		t.Fatalf("Failed to encode to JSON: %s", err)
	}

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_EMAIL_REQUIRED}, response.Errors)
}

func TestCreateAccountWithInvalidPassword(t *testing.T) {

	account := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummy",
	}

	jsonBytes, err := json.Marshal(account)

	if err != nil {
		t.Fatalf("Failed to encode to JSON: %s", err)
	}

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_PASSWORD_REQUIRED}, response.Errors)
}

func TestCreateAccountWithExistsEmail(t *testing.T) {

	deleteDatabase()

	account := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummypassword",
	}

	if err := account.CreateHashedPassword(); err != nil {
		t.Fatal(err)
	}

	if err := infrastructure.GetDB().Create(&account).Error; err != nil {
		t.Fatal(err)
	}

	jsonBytes, err := json.Marshal(account)

	if err != nil {
		t.Fatalf("Failed to encode to JSON: %s", err)
	}

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_EMAIL_ALREADY_EXISTS}, response.Errors)
}

func TestCreateAccount(t *testing.T) {
	deleteDatabase()

	account := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummypassword",
	}

	jsonBytes, err := json.Marshal(account)

	if err != nil {
		t.Fatalf("Failed to encode to JSON: %s", err)
	}

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user", bytes.NewBuffer(jsonBytes), http.StatusCreated)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	var accountNew models.Account
	if err := json.Unmarshal(response.Data, &accountNew); err != nil {
		t.Fatalf("Failed to decode response to payment: %s", err)
	}

	var existingAccountInDB models.Account
	if err := infrastructure.GetDB().Where("ID = ?", accountNew.ID).First(&existingAccountInDB).Error; err != nil {
		t.Fatal(err)
	}

	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(accountNew.Token, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	assert.EqualValues(t, account.Email, existingAccountInDB.Email)
	assert.EqualValues(t, accountNew.Password, "", "Password returned to client")
	assert.True(t, token.Valid, "Token Invalid")
	assert.NotEqual(t, account.Password, existingAccountInDB.Password, "Password not encrypted")

}

func TestLoginWithInvalidBody(t *testing.T) {

	jsonBytes := []byte("{ malformed json }")

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user/login", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_INVALID_JSON}, response.Errors)
}

func TestLoginWithInvalidPassword(t *testing.T) {

	deleteDatabase()

	account := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummypassword",
	}

	if err := account.CreateHashedPassword(); err != nil {
		t.Fatal(err)
	}

	if err := infrastructure.GetDB().Create(&account).Error; err != nil {
		t.Fatal(err)
	}

	accountToLogin := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummy",
	}

	jsonBytes, err := json.Marshal(accountToLogin)

	if err != nil {
		t.Fatalf("Failed to encode to JSON: %s", err)
	}

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user/login", bytes.NewBuffer(jsonBytes), http.StatusUnauthorized)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_INVALID_LOGIN}, response.Errors)
}

func TestLoginWithNonExistsEmail(t *testing.T) {

	deleteDatabase()

	account := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummypassword",
	}

	jsonBytes, err := json.Marshal(account)

	if err != nil {
		t.Fatalf("Failed to encode to JSON: %s", err)
	}

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user/login", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_EMAIL_NON_EXISTS}, response.Errors)
}

func TestLogin(t *testing.T) {

	deleteDatabase()

	account := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummypassword",
	}

	if err := account.CreateHashedPassword(); err != nil {
		t.Fatal(err)
	}

	if err := infrastructure.GetDB().Create(&account).Error; err != nil {
		t.Fatal(err)
	}

	accountToLogin := models.Account{
		Email:    "dummyemail@dummy.com",
		Password: "dummypassword",
	}

	jsonBytes, err := json.Marshal(accountToLogin)

	if err != nil {
		t.Fatalf("Failed to encode to JSON: %s", err)
	}

	rw := doRequestWithoutLogin(t, http.MethodPost, "/v1/user/login", bytes.NewBuffer(jsonBytes), http.StatusOK)
	validateHeaderContentType(t, rw)
	response := decodeApiResponse(t, rw)

	var accountLogged models.Account
	if err := json.Unmarshal(response.Data, &accountLogged); err != nil {
		t.Fatalf("Failed to decode response to payment: %s", err)
	}

	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(accountLogged.Token, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	assert.EqualValues(t, []string(nil), response.Errors)
	assert.EqualValues(t, account.Email, accountLogged.Email)
	assert.EqualValues(t, accountLogged.Password, "", "Password Returned")
	assert.True(t, token.Valid, "Token Invalid")
}
