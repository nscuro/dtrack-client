package dtrack

import (
	"fmt"
	"net/http"
)

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) error {
		if apiKey == "" {
			return fmt.Errorf("no api key provided")
		}

		currentTransport := c.httpClient.Transport
		if currentTransport == nil {
			currentTransport = http.DefaultTransport
		}

		c.httpClient.Transport = &apiKeyTransport{
			apiKey:    apiKey,
			transport: currentTransport,
		}
		return nil
	}
}

type apiKeyTransport struct {
	apiKey    string
	transport http.RoundTripper
}

func (t apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqCopy := *req // Shallow copy of req

	// Deep copy of request headers, because we'll modify them
	reqCopy.Header = make(http.Header, len(req.Header))
	for hn, hv := range req.Header {
		reqCopy.Header[hn] = append([]string(nil), hv...)
	}

	reqCopy.Header.Set("X-Api-Key", t.apiKey)

	return t.transport.RoundTrip(&reqCopy)
}
