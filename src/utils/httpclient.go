package utils

import (
	"bytes"
	"encoding/json"
	"github.com/art-es/ac-export/src/logging"
	"net/http"
	"time"
)

var Client = &http.Client{Timeout: time.Minute * 10}

func SendRequestWithJsonBody(method string, uri string, payload interface{}, headers map[string]string) *http.Response {
	var err error
	jsonPayload := []byte("{}")

	if payload != nil {
		jsonPayload, err = json.Marshal(&payload)
		if err != nil {
			logging.ErrorMarshal(payload, err)
		}
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonPayload))
	if err != nil {
		logging.ErrorCreatingRequest(uri, err)
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return sendRequest(req)
}

func SendRequestWithQueryBody(method string, uri string, payload map[string]string, headers map[string]string) *http.Response {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		logging.ErrorCreatingRequest(uri, err)
	}

	query := req.URL.Query()
	for key, value := range payload {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return sendRequest(req)
}

func sendRequest(req *http.Request) *http.Response {
	resp, err := Client.Do(req)
	if err != nil {
		logging.ErrorSendingRequest(req.URL.String(), err)
	}

	return resp
}
