package dtrack

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func checkResponse(res *http.Response) error {
	if res.StatusCode < 300 {
		return nil
	}

	// TODO: Read body and set .Message

	return &APIError{
		StatusCode: res.StatusCode,
	}
}
