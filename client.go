package dtrack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultTimeout   = 10 * time.Second
	DefaultUserAgent = "github.com/nscuro/dtrack-client"
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
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

func (c Client) newRequest(ctx context.Context, method, path string, params map[string]string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		values := url.Values{}

		for pk, pv := range params {
			values.Add(pk, pv)
		}

		u.RawQuery = values.Encode()
		log.Printf("url: %s", u.String())
	}

	var contentType string
	var bodyBuf io.ReadWriter
	if body != nil {
		switch body := body.(type) {
		case url.Values:
			bodyBuf = bytes.NewBufferString("")
			if _, err = fmt.Fprint(bodyBuf, body.Encode()); err != nil {
				return nil, err
			}
			contentType = "application/x-www-form-urlencoded"
		default:
			bodyBuf = new(bytes.Buffer)
			if err = json.NewEncoder(bodyBuf).Encode(body); err != nil {
				return nil, err
			}
			contentType = "application/json"
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyBuf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", DefaultUserAgent)

	return req, nil
}

type PageOptions struct {
	Offset     int // Offset of the elements to return
	PageNumber int // Page to return
	PageSize   int // Amount of elements to return per page
}

func (c Client) newPagingRequest(ctx context.Context, method, path string, params map[string]string, body interface{}, po PageOptions) (*http.Request, error) {
	paramsCopy := make(map[string]string)
	for k, v := range params {
		paramsCopy[k] = v
	}

	if po.Offset > 0 {
		paramsCopy["offset"] = strconv.Itoa(po.Offset)
	} else if po.PageNumber > 0 {
		paramsCopy["pageNumber"] = strconv.Itoa(po.PageNumber)
	}

	if po.PageSize > 0 {
		paramsCopy["pageSize"] = strconv.Itoa(po.PageSize)
	}

	return c.newRequest(ctx, method, path, paramsCopy, body)
}

func (c Client) doRequest(req *http.Request, v interface{}) (*APIResponse, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = checkResponse(res); err != nil {
		return nil, err
	}

	if err = json.NewDecoder(res.Body).Decode(v); err != nil {
		return nil, err
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
