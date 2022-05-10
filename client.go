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
	"strings"
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
	userAgent  string
	debug      bool

	About           AboutService
	Analysis        AnalysisService
	BOM             BOMService
	Component       ComponentService
	Finding         FindingService
	License         LicenseService
	Metrics         MetricsService
	Project         ProjectService
	ProjectProperty ProjectPropertyService
	Repository      RepositoryService
	User            UserService
	Vulnerability   VulnerabilityService
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
			Timeout: DefaultTimeout,
		},
		userAgent: DefaultUserAgent,
		debug:     false,
	}

	for _, option := range options {
		if err := option(&client); err != nil {
			return nil, err
		}
	}

	client.About = AboutService{client: &client}
	client.Analysis = AnalysisService{client: &client}
	client.BOM = BOMService{client: &client}
	client.Component = ComponentService{client: &client}
	client.Finding = FindingService{client: &client}
	client.License = LicenseService{client: &client}
	client.Metrics = MetricsService{client: &client}
	client.Project = ProjectService{client: &client}
	client.ProjectProperty = ProjectPropertyService{client: &client}
	client.Repository = RepositoryService{client: &client}
	client.User = UserService{client: &client}
	client.Vulnerability = VulnerabilityService{client: &client}

	return &client, nil
}

// BaseURL provides a copy of the Dependency-Track base URL.
func (c Client) BaseURL() *url.URL {
	u := *c.baseURL
	return &u
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
	req.Header.Set("User-Agent", c.userAgent)

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

type Page[T any] struct {
	Items      []T
	TotalCount int
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

// FetchAll is a convenience function to retrieve all items of a paginated API resource.
func FetchAll[T any](f func(po PageOptions) (Page[T], error)) (items []T, err error) {
	const pageSize = 50
	pageNumber := 1

	for {
		page, fErr := f(PageOptions{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if fErr != nil {
			err = fErr
			break
		}

		items = append(items, page.Items...)
		if len(items) > page.TotalCount {
			break
		}

		pageNumber++
	}

	return
}

func (c Client) doRequest(req *http.Request, v interface{}) (a apiResponse, err error) {
	if c.debug {
		reqDump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("sending request:\n>>>>>>\n%s\n>>>>>>\n", string(reqDump))
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if c.debug {
		resDump, _ := httputil.DumpResponse(res, true)
		log.Printf("received response:\n<<<<<<\n%s\n<<<<<<\n", string(resDump))
	}

	err = checkResponse(res)
	if err != nil {
		return
	}

	if v != nil {
		switch vt := v.(type) {
		case *string:
			if content, vErr := io.ReadAll(res.Body); vErr != nil {
				err = vErr
				return
			} else {
				*vt = strings.TrimSpace(string(content))
			}
		default:
			if err = json.NewDecoder(res.Body).Decode(v); err != nil {
				return
			}
		}
	}

	a, err = c.newAPIResponse(res)
	return
}

type apiResponse struct {
	*http.Response

	TotalCount int
}

func (c Client) newAPIResponse(res *http.Response) (a apiResponse, err error) {
	a = apiResponse{Response: res}

	totalCount, ok := a.Header["X-Total-Count"]
	if ok && len(totalCount) > 0 {
		totalCountVal, vErr := strconv.Atoi(totalCount[0])
		if vErr != nil {
			err = vErr
			return
		}
		a.TotalCount = totalCountVal
	}

	return
}

type ClientOption func(*Client) error

// WithDebug toggles the debug mode.
// When enabled, HTTP requests and responses will be logged to stderr.
func WithDebug(debug bool) ClientOption {
	return func(c *Client) error {
		c.debug = debug
		return nil
	}
}

// WithUserAgent overrides the default user agent.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) error {
		c.userAgent = userAgent
		return nil
	}
}

// WithTimeout overrides the default timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) error {
		c.httpClient.Timeout = timeout
		return nil
	}
}
