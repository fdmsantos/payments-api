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

// CreateApiErrorResponse to create an error response
func CreateApiErrorResponse(w http.ResponseWriter, error string, httpStatusCode int) {
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

func CreateApiResponse(w http.ResponseWriter, response interface{}, httpStatusCode int, links []Link) {

	// Encode the response to JSON
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	apiResponse, err := json.Marshal(Response{
		Data:  data,
		Links: links,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, err = w.Write(apiResponse)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
