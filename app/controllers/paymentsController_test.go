package controllers

import (
	"bytes"
	"github.com/jinzhu/gorm"

	//"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	//"github.com/jinzhu/gorm"
	//"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"os"
	//"payments/app"
	"payments/app/models"
	"payments/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
	//"github.com/go-pg/pg"
)

var server *http.Server

func TestMain(m *testing.M) {
	router := mux.NewRouter()

	router.HandleFunc("/v1/user/new", CreateAccount).Methods(http.MethodPost)
	router.HandleFunc("/v1/user/login", Authenticate).Methods(http.MethodPost)
	router.HandleFunc("/v1/payments", CreatePayment).Methods(http.MethodPost)
	router.HandleFunc("/v1/payments", GetPayments).Methods(http.MethodGet)
	router.HandleFunc("/v1/payments/{id}", GetPayment).Methods(http.MethodGet)
	router.HandleFunc("/v1/payments/{id}", UpdatePayment).Methods(http.MethodPut)
	router.HandleFunc("/v1/payments/{id}", DeletePayment).Methods(http.MethodDelete)

	//router.Use(app.JwtAuthentication) //attach JWT auth middleware

	//router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	//err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000/api
	//if err != nil {
	//	fmt.Print(err)
	//}

	server = &http.Server{Addr: ":8080", Handler: router}

	code := m.Run()

	os.Exit(code)
}

func emptyDatabase(t *testing.T) {

	// remove all rows from all of the tables
	models.GetDB().Delete(&models.Payment{})
	models.GetDB().Delete(&models.Attributes{})
	models.GetDB().Delete(&models.BeneficiaryParty{})
	models.GetDB().Delete(&models.DebtorParty{})
	models.GetDB().Delete(&models.SponsorParty{})
	models.GetDB().Delete(&models.ChargesInformation{})
	models.GetDB().Delete(&models.Charge{})
	models.GetDB().Delete(&models.FX{})
}

func paymentExample(paymentId uuid.UUID) []byte {
	return []byte(`
	{
		"type": "Payment",
		"id": "` + paymentId.String() + `",
		"version": 0,
		"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
		"attributes": {
			"amount": "100.21",
			"beneficiary_party": {
				"account_name": "W Owens",
				"account_number": "31926819",
				"account_number_code": "BBAN",
				"account_type": 0,
				"address": "1 The Beneficiary Localtown SE2",
				"bank_id": "403000",
				"bank_id_code": "GBDSC",
				"name": "Wilfred Jeremiah Owens"
			},
			"charges_information": {
				"bearer_code": "SHAR",
				"sender_charges": [{
					"amount": "5.00",
					"currency": "GBP"
				},
				{
					"amount": "10.00",
					"currency": "USD"
				}],
				"receiver_charges_amount": "1.00",
				"receiver_charges_currency": "USD"
			},
			"currency": "GBP",
			"debtor_party": {
				"account_name": "EJ Brown Black",
				"account_number": "GB29XABC10161234567801",
				"account_number_code": "IBAN",
				"address": "10 Debtor Crescent Sourcetown NE1",
				"bank_id": "203301",
				"bank_id_code": "GBDSC",
				"name": "Emelia Jane Brown"
			},
			"end_to_end_reference": "Wil piano Jan",
			"fx": {
				"contract_reference": "FX123",
				"exchange_rate": "2.00000",
				"original_amount": "200.42",
				"original_currency": "USD"
			},
			"numeric_reference": "1002001",
			"payment_id": "123456789012345678",
			"payment_purpose": "Paying for goods/services",
			"payment_scheme": "FPS",
			"payment_type": "Credit",
			"processing_date": "2017-01-18",
			"reference": "Payment for Em\u0027s piano lessons",
			"scheme_payment_sub_type": "InternetBanking",
			"scheme_payment_type": "ImmediatePayment",
			"sponsor_party": {
				"account_number": "56781234",
				"bank_id": "123123",
				"bank_id_code": "GBDSC"
			}
		}
	}
	`)
}

func insertPayments(t *testing.T, paymentId uuid.UUID) models.Payment {

	// populate table with example payment
	var payment models.Payment
	if err := json.Unmarshal(paymentExample(paymentId), &payment); err != nil {
		t.Fatal(err)
	}

	if err := models.GetDB().Create(&payment).Error; err != nil {
		t.Fatal(err)
	}

	return payment
}

func TestGetPaymentsWithEmptyTable(t *testing.T) {

	emptyDatabase(t)

	req := httptest.NewRequest(http.MethodGet, "/v1/payments", nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 200 {
		t.Fatalf("Status code was not 200: %d\n", rw.Code)
	}

	if rw.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content type was not application/json")
	}

	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}

	var payments []models.Payment
	if err := json.Unmarshal(response.Data, &payments); err != nil {
		t.Fatalf("Failed to decode response to payments slice: %s", err)
	}

	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: "/v1/payments"}}, response.Links)

	assert.Len(t, payments, 0, "Payments array must be empty when database is empty")
}

func TestGetPaymentsWithOneExistingPayment(t *testing.T) {

	emptyDatabase(t)

	payment := insertPayments(t, uuid.NewV1())

	req := httptest.NewRequest(http.MethodGet, "/v1/payments", nil)
	//req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjJ9.DzOJ7GHkPwiDE3T78dFMriY96VwzytQSBV7-c64dxx8")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)

	if rw.Code != 200 {
		t.Fatalf("Status code was not 200: %d\n", rw.Code)
	}

	if rw.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content type was not application/json")
	}

	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}

	var payments []models.Payment
	if err := json.Unmarshal(response.Data, &payments); err != nil {
		t.Fatalf("Failed to decode response to payments slice: %s", err)
	}

	expectedPaymentJson, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	atualPaymentJson, err := json.Marshal(payments[0])
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	require.Len(t, payments, 1, "Payments array must contain one payment when database has one payment")
	assert.JSONEq(t, string(expectedPaymentJson), string(atualPaymentJson))
	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: "/v1/payments"}}, response.Links)
}

func TestGetPaymentsWithMultipleExistingPayments(t *testing.T) {

	emptyDatabase(t)

	var expectedPayments []models.Payment

	expectedPayments = append(expectedPayments, insertPayments(t, uuid.NewV1()))
	expectedPayments = append(expectedPayments, insertPayments(t, uuid.NewV1()))

	req := httptest.NewRequest(http.MethodGet, "/v1/payments", nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 200 {
		t.Fatalf("Status code was not 200: %d\n", rw.Code)
	}

	if rw.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content type was not application/json")
	}

	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}

	var payments []models.Payment
	if err := json.Unmarshal(response.Data, &payments); err != nil {
		t.Fatalf("Failed to decode response to payments slice: %s", err)
	}

	var expectedPaymentsJson []string
	var atualPayments []string

	for _, payment := range expectedPayments {
		paymentJson, err := json.Marshal(payment)
		if err != nil {
			t.Fatalf("Failed to encode json")
		}
		expectedPaymentsJson = append(expectedPaymentsJson, string(paymentJson))
	}

	for _, payment := range payments {
		paymentJson, err := json.Marshal(payment)
		if err != nil {
			t.Fatalf("Failed to encode json")
		}
		atualPayments = append(atualPayments, string(paymentJson))
	}

	require.Len(t, payments, 2, "Payments array must contain two payments when database has two payments")

	for i := 0; i < len(atualPayments); i++ {
		assert.JSONEq(t, expectedPaymentsJson[i], atualPayments[i])
	}

	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: "/v1/payments"}}, response.Links)
}

func TestGetSinglePaymentWithOneExistingPayment(t *testing.T) {

	emptyDatabase(t)

	testPayment := insertPayments(t, uuid.NewV1())

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/payments/%s", testPayment.ID.String()), nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 200 {
		t.Fatalf("Status code was not 200: %d\n", rw.Code)
	}

	if rw.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content type was not application/json")
	}

	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}

	var payment models.Payment
	if err := json.Unmarshal(response.Data, &payment); err != nil {
		t.Fatalf("Failed to decode response to payment: %s", err)
	}

	expectedPaymentJson, err := json.Marshal(testPayment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	atualPaymentJson, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	assert.EqualValues(t, expectedPaymentJson, atualPaymentJson)

	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: fmt.Sprintf("/v1/payments/%s", testPayment.ID.String())}}, response.Links)
}

//
func TestGetSinglePaymentForNonExistingPayment(t *testing.T) {

	emptyDatabase(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 404 {
		t.Fatalf("Status code was not 404: %d\n", rw.Code)
	}
	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Payment not found"}, response.Errors)
}

func TestGetSinglePaymentForInvalidUUID(t *testing.T) {

	emptyDatabase(t)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/payments/%s", "TestUUID"), nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 400 {
		t.Fatalf("Status code was not 400: %d\n", rw.Code)
	}
	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Requested UUID is Invalid"}, response.Errors)
}

func TestGetSinglePaymentForNonExistingPaymentWhenOtherPaymentExists(t *testing.T) {

	emptyDatabase(t)

	_ = insertPayments(t, uuid.NewV1())

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 404 {
		t.Fatalf("Status code was not 404: %d\n", rw.Code)
	}
	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Payment not found"}, response.Errors)
}

func TestCreateSinglePayment(t *testing.T) {

	emptyDatabase(t)

	testPaymentBytes := paymentExample(uuid.NewV1())

	//jsonBytes, err := json.Marshal(testPayment)
	//require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(testPaymentBytes))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 201 {
		t.Fatalf("Status code was not 201: %d\n", rw.Code)
	}

	var testPayment models.Payment
	if err := json.Unmarshal(testPaymentBytes, &testPayment); err != nil {
		t.Fatalf("Failed to decode payment: %s", err)
	}

	assert.Equal(t, fmt.Sprintf("/v1/payments/%s", testPayment.ID.String()), rw.Header().Get("Location"))

	actualPayment := models.Payment{
		ID: testPayment.ID,
	}

	models.GetDB().Set("gorm:auto_preload", true).Find(&actualPayment)

	expectedPaymentJson, err := json.Marshal(testPayment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	atualPaymentJson, err := json.Marshal(actualPayment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	require.Nil(t, err)

	assert.JSONEq(t, string(expectedPaymentJson), string(atualPaymentJson))
}

func TestCreateSinglePaymentThatAlreadyExists(t *testing.T) {

	emptyDatabase(t)

	// populate table with example payment
	examplePayment := insertPayments(t, uuid.NewV1())

	jsonBytes, err := json.Marshal(examplePayment)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)

	retry := httptest.NewRecorder()
	server.Handler.ServeHTTP(retry, req)
	if retry.Code != 400 {
		t.Fatalf("Status code was not 400: %d\n", retry.Code)
	}
}

func TestCreateSinglePaymentWithInvalidJSON(t *testing.T) {

	emptyDatabase(t)

	jsonBytes := []byte("{ malformed json }")
	req := httptest.NewRequest(http.MethodPost, "/v1/payments", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 400 {
		t.Fatalf("Status code was not 400: %d\n", rw.Code)
	}
	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Invalid JSON"}, response.Errors)
}

func TestUpdatePayment(t *testing.T) {

	emptyDatabase(t)

	// populate table with example payment
	testPayment := insertPayments(t, uuid.NewV1())

	testPayment.Attributes.Currency = "Euro"

	jsonBytes, err := json.Marshal(testPayment)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/payments/%s", testPayment.ID), bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 204 {
		t.Fatalf("Status code was not 204: %d\n", rw.Code)
	}
	assert.Equal(t, fmt.Sprintf("/v1/payments/%s", testPayment.ID.String()), rw.Header().Get("Location"))

	actualPayment := models.Payment{
		ID: testPayment.ID,
	}

	models.GetDB().Set("gorm:auto_preload", true).Find(&actualPayment)

	expectedPaymentJson, err := json.Marshal(testPayment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	atualPaymentJson, err := json.Marshal(actualPayment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}

	assert.JSONEq(t, string(expectedPaymentJson), string(atualPaymentJson))
}

func TestUpdateSinglePaymentWithIDThatDoesNotMatchURL(t *testing.T) {

	emptyDatabase(t)

	// populate table with example payment
	testPayment := insertPayments(t, uuid.NewV1())

	jsonBytes, err := json.Marshal(testPayment)
	require.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 400 {
		t.Fatalf("Status code was not 400: %d\n", rw.Code)
	}
	var response utils.Response
	if err := json.NewDecoder(rw.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Mismatching IDs"}, response.Errors)
}

func TestUpdateNonExistentPayment(t *testing.T) {

	emptyDatabase(t)

	// populate table with example payment
	id := uuid.NewV1()
	testPayment := paymentExample(id)

	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/payments/%s", id.String()), bytes.NewBuffer(testPayment))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 404 {
		t.Fatalf("Status code was not 404: %d\n", rw.Code)
	}
	var response utils.Response
	if err := json.NewDecoder(rw.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Payment not found"}, response.Errors)
}

func TestUpdateSinglePaymentWithInvalidJSON(t *testing.T) {

	emptyDatabase(t)

	examplePayment := insertPayments(t, uuid.NewV1())

	jsonBytes := []byte("{ Malformed json }")
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/v1/payments/%s", examplePayment.ID.String()), bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 400 {
		t.Fatalf("Status code was not 400: %d\n", rw.Code)
	}
	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Invalid JSON"}, response.Errors)
}

func TestDeletePayment(t *testing.T) {

	emptyDatabase(t)

	testPayment := insertPayments(t, uuid.NewV1())

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/payments/%s", testPayment.ID), nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 204 {
		t.Fatalf("Status code was not 204: %d\n", rw.Code)
	}

	err := models.GetDB().Where("ID = ?", testPayment.ID).First(&models.Payment{}).Error
	assert.True(t, gorm.IsRecordNotFoundError(err))
}

func TestDeleteNonExistingPayment(t *testing.T) {

	emptyDatabase(t)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), nil)
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)
	if rw.Code != 404 {
		t.Fatalf("Status code was not 404: %d\n", rw.Code)
	}
	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	assert.EqualValues(t, []string{"Payment not found"}, response.Errors)

}
