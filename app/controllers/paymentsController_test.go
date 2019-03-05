package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"payments/app/models"
	"payments/utils"
	"testing"
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

	server = &http.Server{Addr: ":8000", Handler: router}

	code := m.Run()

	os.Exit(code)
}

func deleteDatabase(t *testing.T) {
	utils.GetDB().Delete(&models.Payment{})
	utils.GetDB().Delete(&models.Attributes{})
	utils.GetDB().Delete(&models.BeneficiaryParty{})
	utils.GetDB().Delete(&models.DebtorParty{})
	utils.GetDB().Delete(&models.SponsorParty{})
	utils.GetDB().Delete(&models.ChargesInformation{})
	utils.GetDB().Delete(&models.Charge{})
	utils.GetDB().Delete(&models.FX{})
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
	var payment models.Payment
	if err := json.Unmarshal(paymentExample(paymentId), &payment); err != nil {
		t.Fatal(err)
	}

	if err := utils.GetDB().Create(&payment).Error; err != nil {
		t.Fatal(err)
	}

	return payment
}

func convertToJson(t *testing.T, payment models.Payment) []byte {
	paymentJson, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to encode json")
	}
	return paymentJson
}

func doRequest(t *testing.T, method string, url string, body io.Reader, expectedResultCode int) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	rw := httptest.NewRecorder()
	server.Handler.ServeHTTP(rw, req)

	if rw.Code != expectedResultCode {
		t.Fatalf("Status code was not %d: %d\n", expectedResultCode, rw.Code)
	}

	return rw
}

func validateHeaderContenType(t *testing.T, rw *httptest.ResponseRecorder) {
	if rw.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Content type was not application/json")
	}
}

func decodeApiResponse(t *testing.T, rw *httptest.ResponseRecorder) utils.Response {
	var response utils.Response
	err := json.NewDecoder(rw.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode API response: %s", err)
	}
	return response
}

func convertJsonToPayment(t *testing.T, rw *httptest.ResponseRecorder) (utils.Response, models.Payment) {
	response := decodeApiResponse(t, rw)
	var payment models.Payment
	if err := json.Unmarshal(response.Data, &payment); err != nil {
		t.Fatalf("Failed to decode response to payment: %s", err)
	}
	return response, payment
}

func convertJsonToPayments(t *testing.T, rw *httptest.ResponseRecorder) (utils.Response, []models.Payment) {
	var payments []models.Payment
	response := decodeApiResponse(t, rw)
	if err := json.Unmarshal(response.Data, &payments); err != nil {
		t.Fatalf("Failed to decode response to payments slice: %s", err)
	}

	return response, payments
}

func TestGetPaymentsWithEmptyDatabase(t *testing.T) {

	deleteDatabase(t)

	rw := doRequest(t, http.MethodGet, "/v1/payments", nil, http.StatusOK)
	validateHeaderContenType(t, rw)
	response, payments := convertJsonToPayments(t, rw)

	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: "/v1/payments"}}, response.Links)
	assert.Len(t, payments, 0, "Payments array must be empty when database is empty")
}

func TestGetPaymentsWithOnePayment(t *testing.T) {

	deleteDatabase(t)

	payment := insertPayments(t, uuid.NewV1())

	rw := doRequest(t, http.MethodGet, "/v1/payments", nil, http.StatusOK)
	validateHeaderContenType(t, rw)

	response, payments := convertJsonToPayments(t, rw)

	expectedPaymentJson := convertToJson(t, payment)
	atualPaymentJson := convertToJson(t, payments[0])

	require.Len(t, payments, 1, "Payments array must contain one payment when database has one payment")
	assert.JSONEq(t, string(expectedPaymentJson), string(atualPaymentJson))
	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: "/v1/payments"}}, response.Links)
}

func TestGetPaymentsWithMultiplePayments(t *testing.T) {

	deleteDatabase(t)

	var expectedPayments []models.Payment

	expectedPayments = append(expectedPayments, insertPayments(t, uuid.NewV1()))
	expectedPayments = append(expectedPayments, insertPayments(t, uuid.NewV1()))

	rw := doRequest(t, http.MethodGet, "/v1/payments", nil, http.StatusOK)
	validateHeaderContenType(t, rw)
	response, payments := convertJsonToPayments(t, rw)

	var expectedPaymentsJson []string
	var atualPayments []string

	for _, payment := range expectedPayments {
		expectedPaymentsJson = append(expectedPaymentsJson, string(convertToJson(t, payment)))
	}

	for _, payment := range payments {
		atualPayments = append(atualPayments, string(convertToJson(t, payment)))
	}

	require.Len(t, payments, 2, "Payments array must contain two payments when database has two payments")
	for i := 0; i < len(atualPayments); i++ {
		assert.JSONEq(t, expectedPaymentsJson[i], atualPayments[i])
	}
	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: "/v1/payments"}}, response.Links)
}

func TestGetSinglePaymentWithOnePayment(t *testing.T) {

	deleteDatabase(t)

	testPayment := insertPayments(t, uuid.NewV1())

	rw := doRequest(t, http.MethodGet, fmt.Sprintf("/v1/payments/%s", testPayment.ID.String()), nil, http.StatusOK)
	validateHeaderContenType(t, rw)
	response, payment := convertJsonToPayment(t, rw)

	assert.EqualValues(t, convertToJson(t, testPayment), convertToJson(t, payment))
	assert.EqualValues(t, []utils.Link{{Rel: "self", Href: fmt.Sprintf("/v1/payments/%s", testPayment.ID.String())}}, response.Links)
}

func TestGetSinglePaymentForNonExistingPayment(t *testing.T) {

	deleteDatabase(t)

	rw := doRequest(t, http.MethodGet, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), nil, http.StatusNotFound)
	validateHeaderContenType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_RESOURCE_NOT_FOUND}, response.Errors)
}

func TestGetSinglePaymentForInvalidUUID(t *testing.T) {

	deleteDatabase(t)

	rw := doRequest(t, http.MethodGet, fmt.Sprintf("/v1/payments/%s", "TestUUID"), nil, http.StatusBadRequest)
	validateHeaderContenType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{ERROR_REQUESTED_UUID_INVALID}, response.Errors)
}

func TestGetSinglePaymentForNonExistingPaymentWhenOtherPaymentExists(t *testing.T) {

	deleteDatabase(t)

	_ = insertPayments(t, uuid.NewV1())

	rw := doRequest(t, http.MethodGet, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), nil, http.StatusNotFound)
	validateHeaderContenType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_RESOURCE_NOT_FOUND}, response.Errors)
}

func TestCreateSinglePayment(t *testing.T) {

	deleteDatabase(t)

	testPaymentBytes := paymentExample(uuid.NewV1())

	rw := doRequest(t, http.MethodPost, "/v1/payments", bytes.NewBuffer(testPaymentBytes), http.StatusCreated)
	validateHeaderContenType(t, rw)

	var testPayment models.Payment
	if err := json.Unmarshal(testPaymentBytes, &testPayment); err != nil {
		t.Fatalf("Failed to decode payment: %s", err)
	}
	assert.Equal(t, fmt.Sprintf("/v1/payments/%s", testPayment.ID.String()), rw.Header().Get("Location"))

	actualPayment := models.Payment{
		ID: testPayment.ID,
	}

	utils.GetDB().Set("gorm:auto_preload", true).Find(&actualPayment)

	assert.JSONEq(t, string(convertToJson(t, testPayment)), string(convertToJson(t, actualPayment)))
}

func TestCreateSinglePaymentThatExists(t *testing.T) {

	deleteDatabase(t)

	examplePayment := insertPayments(t, uuid.NewV1())
	jsonBytes, err := json.Marshal(examplePayment)
	require.Nil(t, err)
	_ = doRequest(t, http.MethodPost, "/v1/payments", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
}

func TestCreateSinglePaymentWithInvalidBody(t *testing.T) {

	deleteDatabase(t)

	jsonBytes := []byte("{ malformed json }")

	rw := doRequest(t, http.MethodPost, "/v1/payments", bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	validateHeaderContenType(t, rw)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_INVALID_JSON}, response.Errors)
}

func TestUpdatePayment(t *testing.T) {

	deleteDatabase(t)

	testPayment := insertPayments(t, uuid.NewV1())
	testPayment.Attributes.Currency = "Euro"

	jsonBytes, err := json.Marshal(testPayment)
	require.Nil(t, err)

	rw := doRequest(t, http.MethodPut, fmt.Sprintf("/v1/payments/%s", testPayment.ID), bytes.NewBuffer(jsonBytes), http.StatusNoContent)

	assert.Equal(t, fmt.Sprintf("/v1/payments/%s", testPayment.ID.String()), rw.Header().Get("Location"))

	actualPayment := models.Payment{
		ID: testPayment.ID,
	}

	utils.GetDB().Set("gorm:auto_preload", true).Find(&actualPayment)
	assert.JSONEq(t, string(convertToJson(t, testPayment)), string(convertToJson(t, actualPayment)))
}

func TestUpdateSinglePaymentWithIDThatDoesNotMatchURL(t *testing.T) {

	deleteDatabase(t)

	testPayment := insertPayments(t, uuid.NewV1())

	jsonBytes, err := json.Marshal(testPayment)
	require.Nil(t, err)

	rw := doRequest(t, http.MethodPut, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{ERROR_ID_MISMATCH}, response.Errors)
}

func TestUpdateNonExistentPayment(t *testing.T) {

	deleteDatabase(t)

	id := uuid.NewV1()
	testPayment := paymentExample(id)

	rw := doRequest(t, http.MethodPut, fmt.Sprintf("/v1/payments/%s", id.String()), bytes.NewBuffer(testPayment), http.StatusNotFound)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_RESOURCE_NOT_FOUND}, response.Errors)
}

func TestUpdateSinglePaymentWithInvalidBody(t *testing.T) {

	deleteDatabase(t)

	examplePayment := insertPayments(t, uuid.NewV1())

	jsonBytes := []byte("{ Malformed json }")
	rw := doRequest(t, http.MethodPut, fmt.Sprintf("/v1/payments/%s", examplePayment.ID.String()), bytes.NewBuffer(jsonBytes), http.StatusBadRequest)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_INVALID_JSON}, response.Errors)
}

func TestDeletePayment(t *testing.T) {

	deleteDatabase(t)

	testPayment := insertPayments(t, uuid.NewV1())
	_ = doRequest(t, http.MethodDelete, fmt.Sprintf("/v1/payments/%s", testPayment.ID), nil, http.StatusNoContent)

	err := utils.GetDB().Where("ID = ?", testPayment.ID).First(&models.Payment{}).Error
	assert.True(t, gorm.IsRecordNotFoundError(err))
}

func TestDeleteNonExistingPayment(t *testing.T) {

	deleteDatabase(t)
	rw := doRequest(t, http.MethodDelete, fmt.Sprintf("/v1/payments/%s", uuid.NewV1().String()), nil, http.StatusNotFound)
	response := decodeApiResponse(t, rw)

	assert.EqualValues(t, []string{utils.ERROR_RESOURCE_NOT_FOUND}, response.Errors)

}
