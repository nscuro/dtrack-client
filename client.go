package dtrack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultTimeout   = 10 * time.Second
	DefaultUserAgent = "github.com/nscuro/dtrack-client"
)

type contextKey string

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	debug      bool
}

func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("no api base url provided")
	}

	u, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, err
	}

	client := Client{
		baseURL: u,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	for _, option := range options {
		if err := option(&client); err != nil {
			return nil, err
		}
	}

	return &client, nil
}

func (c Client) newRequest(ctx context.Context, method, path string, options ...requestOption) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", DefaultUserAgent)

	for _, option := range options {
		if err = option(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

type requestOption func(*http.Request) error

func withParams(params map[string]string) requestOption {
	return func(req *http.Request) error {
		if len(params) == 0 {
			return nil
		}

		query := req.URL.Query()

		for pk, pv := range params {
			query.Add(pk, pv)
		}

		req.URL.RawQuery = query.Encode()

		return nil
	}
}

func withBody(body interface{}) requestOption {
	return func(req *http.Request) error {
		if body == nil {
			return nil
		}

		var (
			contentType string
			bodyBuf     io.ReadWriter
		)

		switch body := body.(type) {
		case url.Values:
			bodyBuf = bytes.NewBufferString("")
			if _, err := fmt.Fprint(bodyBuf, body.Encode()); err != nil {
				return err
			}
			contentType = "application/x-www-form-urlencoded"
		default:
			bodyBuf = new(bytes.Buffer)
			if err := json.NewEncoder(bodyBuf).Encode(body); err != nil {
				return err
			}
			contentType = "application/json"
		}

		req.Body = io.NopCloser(bodyBuf)
		req.Header.Set("Content-Type", contentType)

		return nil
	}
}

type PageOptions struct {
	Offset     int // Offset of the elements to return
	PageNumber int // Page to return
	PageSize   int // Amount of elements to return per page
}

func withPageOptions(po PageOptions) requestOption {
	return func(req *http.Request) error {
		query := req.URL.Query()

		if po.Offset > 0 {
			query.Set("offset", strconv.Itoa(po.Offset))
		} else if po.PageNumber > 0 {
			query.Set("pageNumber", strconv.Itoa(po.PageNumber))
		}

		if po.PageSize > 0 {
			query.Set("pageSize", strconv.Itoa(po.PageSize))
		}

		req.URL.RawQuery = query.Encode()

		return nil
	}
}

func (c Client) doRequest(req *http.Request, v interface{}) (*APIResponse, error) {
	if c.debug {
		reqDump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("sending request:\n%s", string(reqDump))
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if c.debug {
		resDump, _ := httputil.DumpResponse(res, true)
		log.Printf("received response:\n%s", string(resDump))
	}

	if err = checkResponse(res); err != nil {
		return nil, err
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(v); err != nil {
			return nil, err
		}
	}

	apiResponse, err := c.newAPIResponse(res)
	if err != nil {
		return nil, err
	}

	return apiResponse, nil
}

type APIResponse struct {
	*http.Response

	TotalCount int
}

func (c Client) newAPIResponse(res *http.Response) (*APIResponse, error) {
	response := APIResponse{Response: res}

	totalCount, ok := response.Header["X-Total-Count"]
	if ok && len(totalCount) > 0 {
		totalCountVal, err := strconv.Atoi(totalCount[0])
		if err != nil {
			return nil, err
		}
		response.TotalCount = totalCountVal
	}

	return &response, nil
}

type ClientOption func(*Client) error

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}
