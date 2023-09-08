package cmc

import "net/http"

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Result struct {
	Status struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message,omitempty"`
	}
	Data []struct {
		Quote map[string]struct {
			Price float64 `json:"price"`
		} `json:"quote"`
	} `json:"data"`
}
