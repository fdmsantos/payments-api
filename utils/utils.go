package utils

import (
	"encoding/json"
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
