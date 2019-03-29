package utils

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
)

var log *logrus.Logger
var once sync.Once

func GetLog() *logrus.Logger {
	logrus.New()
	once.Do(func() {
		log = logrus.New()
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.DebugLevel)
	})
	return log
}

func LogApiRequest(r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // To Read body again in handler
	bodyString := string(bodyBytes)
	body := map[string]interface{}{}
	if err := json.Unmarshal([]byte(bodyString), &body); err == nil {
		delete(body, "password") // Delete password from log
	}

	header := r.Header
	header.Del("Authorization") // Delete Authorization from log

	GetLog().WithFields(logrus.Fields{
		"endpoint": r.RequestURI,
		"method":   r.Method,
		"header":   header,
		"body":     body,
		"host":     r.Host,
	}).Info("Request")

}

func LogApiBadRequestResponse(httpStatusCode int, response []byte) {
	GetLog().WithFields(logrus.Fields{
		"httpStatusCode": httpStatusCode,
		"response":       string(response),
	}).Warn("Request Response")
}

func LogError(httpStatusCode int, error string) {
	GetLog().WithFields(logrus.Fields{
		"httpStatusCode": httpStatusCode,
		"error":          error,
	}).Error("Request Response")
}

func LogApiResponse(httpStatusCode int, response []byte) {
	body := map[string]interface{}{}
	if err := json.Unmarshal(response, &body); err == nil {
		switch body["data"].(type) {
		case map[string]interface{}:
			delete(body["data"].(map[string]interface{}), "token") // Delete token from log
		default:

		}

		GetLog().WithFields(logrus.Fields{
			"httpStatusCode": httpStatusCode,
			"response":       body,
		}).Info("Request Response")
	}
}
