package dtrack

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	restClient *resty.Client
	apiBaseURL *url.URL
	apiKey     string
}

func NewClient(apiBaseURL, apiKey string) (*Client, error) {
	if apiBaseURL == "" {
		return nil, fmt.Errorf("no api base url provided")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("no api key provided")
	}

	u, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, err
	}

	rc := resty.New()
	rc.SetHeader("Accept", "application/json")
	rc.SetHeader("X-Api-Key", apiKey)
	rc.SetHostURL(u.String())

	return &Client{restClient: rc, apiBaseURL: u, apiKey: apiKey}, nil
}
