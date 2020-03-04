package form3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"time"
)

//Client struct
type Client struct {
	HTTPClient *http.Client
	BaseURL    *url.URL
	Debug      bool
	AuthKey    string
}

// restClient struct
type restClient struct {
	URI         *url.URL
	Method      string
	ErrorRef    interface{}
	Body        io.Reader
	Debug       bool
	Headers     map[string]string
	HTTPClient  *http.Client
	ResponseRef interface{}
}

const baseURL = "https://api.staging-form3.tech/v1"

// NewHTTPClient creates a new Client
// if httpClient is nil then a DefaultClient is used
func NewHTTPClient(httpClient *http.Client, baseURL *url.URL) *Client {
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 5 * time.Minute,
		}
	}
	client := &Client{
		HTTPClient: httpClient,
		BaseURL:    baseURL,
	}

	return client
}

// Start initializes the start anonymous function and sets the authorozation token
func (c *Client) Start(responseRef interface{}, errorRef interface{}) *restClient {
	return c.StartAnonymous(responseRef, errorRef).SetAuthorization(c.AuthKey)
}

// StartAnonymous creates rest client setting the content type
func (c *Client) StartAnonymous(responseRef interface{}, errorRef interface{}) *restClient {
	rc := &restClient{
		Debug:       c.Debug,
		ErrorRef:    errorRef,
		Headers:     make(map[string]string),
		HTTPClient:  c.HTTPClient,
		ResponseRef: responseRef,
	}
	rc.URI, _ = url.Parse(c.BaseURL.String())
	rc.WithHeader("Accept", "application/vnd.api+json")
	return rc
}

func (rc *restClient) Call() error {
	req, err := http.NewRequest(rc.Method, rc.URI.String(), rc.Body)
	if err != nil {
		return err
	}
	for key, val := range rc.Headers {
		req.Header.Set(key, val)
	}
	resp, err := rc.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if rc.Debug {
		responseDump, _ := httputil.DumpResponse(resp, true)
		fmt.Println(string(responseDump))
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if rc.ErrorRef != nil {
			err = json.NewDecoder(resp.Body).Decode(rc.ErrorRef)
		}
	} else {
		rc.ErrorRef = nil
		if _, ok := rc.ResponseRef.(*BaseHTTPResponse); !ok {
			err = json.NewDecoder(resp.Body).Decode(rc.ResponseRef)
		}
	}
	rc.ResponseRef.(StatusAble).SetStatus(resp.StatusCode)
	return err
}

func (rc *restClient) SetHeader(key string, value string) *restClient {
	rc.Headers[key] = value
	return rc
}

func (rc *restClient) SetAuthorization(key string) *restClient {
	if key != "" {
		rc.WithHeader("Authorization", key)
	}
	return rc
}

func (rc *restClient) SetJSONBody(body interface{}) *restClient {
	rc.SetHeader("Accept", "application/vnd.api+json")
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(body)
	rc.Body = buffer
	return rc
}

func (rc *restClient) SetMethod(method string) *restClient {
	rc.Method = method
	return rc
}

func (rc *restClient) SetParameter(key string, value interface{}) *restClient {
	q := rc.URI.Query()
	if x, ok := value.([]string); ok {
		for _, i := range x {
			q.Add(key, i)
		}
	} else {
		q.Add(key, fmt.Sprintf("%v", value))
	}
	rc.URI.RawQuery = q.Encode()
	return rc
}

func (rc *restClient) WithURI(uri string) *restClient {
	rc.URI.Path = path.Join(rc.URI.Path, uri)
	return rc
}

func (rc *restClient) WithURISegment(segment string) *restClient {
	if segment != "" {
		rc.URI.Path = path.Join(rc.URI.Path, "/"+segment)
	}
	return rc
}

func (rc *restClient) WithHeader(key string, value string) *restClient {
	rc.Headers[key] = value
	return rc
}

func (rc *restClient) WithMethod(method string) *restClient {
	rc.Method = method
	return rc
}

// ListAccounts retrieves the all accounts of an organization
func (c *Client) ListAccounts(pageNUmber int, pageSize int) (*ListResponse, error) {
	var resp ListResponse

	err := c.Start(&resp, nil).
		WithUri("/organisation/accounts").
		SetParameter("page[number]", pageNUmber).
		SetParameter("page[size]", pageSize).
		WithMethod(http.MethodGet).
		Call()
	return &resp, err
}

// CreateAccount creates a new account
func (c *Client) CreateAccount(request CreateAccountRequest) (*CreateAccountResponse, *Errors, error) {
	var resp CreateAccountResponse
	var errors Errors

	restClient := c.Start(&resp, &errors)
	err := restClient.WithURI("/organisation/accounts").
		WithJSONBody(request).
		WithMethod(http.MethodPost).
		Call()
	if restClient.ErrorRef == nil {
		return &resp, nil, err
	}
	return &resp, &errors, err
}

// FetchAccount action
func (c *Client) FetchAccount(accountID string) (*FetchAccountResponse, *Errors, error) {
	var resp FetchAccountResponse
	var errors Errors

	restClient := c.Start(&resp, &errors)
	err := restClient.WithURI("/organisation/accounts").
		WithURISegment(accountID).
		WithMethod(http.MethodGet).
		Call()
	if restClient.ErrorRef == nil {
		return &resp, nil, err
	}
	return &resp, &errors, err
}

// DeleteAccount action
func (c *Client) DeleteAccount(accountID string, version int) (*BaseHTTPResponse, *Errors, error) {
	var resp BaseHTTPResponse
	var errors Errors

	restClient := c.Start(&resp, &errors)
	err := restClient.WithURI("/organisation/accounts").
		WithURISegment(accountID).
		SetParameter("version", version).
		WithMethod(http.MethodDelete).
		Call()
	if restClient.ErrorRef == nil {
		return &resp, nil, err
	}
	return &resp, &errors, err
}
