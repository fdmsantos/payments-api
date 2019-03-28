package utils

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"net/http"
)

type Response struct {
	Data   json.RawMessage `json:"data,omitempty"`
	Links  []Link          `json:"links,omitempty"`
	Errors []string        `json:"errors,omitempty"`
}

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// ConvertStringToUUID convert a string Id to UUID
func ConvertStringToUUID(id string) (uuid.UUID, error) {
	return uuid.FromString(id)
}

// CreateErrorResponse to create an error response
func CreateErrorResponse(w http.ResponseWriter, error string, httpStatusCode int) {
	// write an error response
	if response, err := json.Marshal(Response{Errors: []string{error}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(httpStatusCode)
		_, err = w.Write(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
