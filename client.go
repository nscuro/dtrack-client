package dtrack

import (
	"net/url"
)

type Client struct {
	apiBaseURL *url.URL
	apiKey     string
}

func NewClient(apiBaseURL, apiKey string) (*Client, error) {
	return nil, nil
}
