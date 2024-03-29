package camunda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/rs/zerolog/log"
)

const (
	PackageVersion        = "{{version}}"
	DefaultUserAgent      = "CamundaClientGo/" + PackageVersion
	DefaultEndpointUrl    = "http://localhost:8080/engine-rest"
	DefaultTimeoutSec     = 60
	DefaultDateTimeFormat = "2006-01-02T15:04:05.000-0700"
)

// ClientOptions a client options
type ClientOptions struct {
	UserAgent   string
	EndpointUrl string
	Timeout     time.Duration
	ApiUser     string
	ApiPassword string
}

// Client a client for Camunda API
type Client struct {
	httpClient  *http.Client
	endpointURL string
	userAgent   string
	apiUser     string
	apiPassword string

	// TaskManager      *TaskManager
	// Deployment        *Deployment
	// ProcessDefinition *ProcessDefinition
	// UserTask          *userTaskApi
}

var ErrorNotFound = &Error{
	Type:    "NotFound",
	Message: "Not found",
}

// Error a custom error type
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Error error message
func (e *Error) Error() string {
	return e.Message
}

// NewClient creates new instance of Client
func NewClient(options *ClientOptions) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: time.Second * DefaultTimeoutSec,
		},
		endpointURL: DefaultEndpointUrl,
		userAgent:   DefaultUserAgent,
		apiUser:     options.ApiUser,
		apiPassword: options.ApiPassword,
	}

	if options.EndpointUrl != "" {
		client.endpointURL = options.EndpointUrl
	}

	if options.UserAgent != "" {
		client.userAgent = options.UserAgent
	}

	if options.Timeout.Nanoseconds() != 0 {
		client.httpClient.Timeout = options.Timeout
	}

	return client
}

func (c *Client) TaskManager() *TaskManager {
	return &TaskManager{
		client: c,
	}
}

func (c *Client) ProcessManager() *ProcessManager {
	return &ProcessManager{
		client: c,
	}
}

// SetCustomTransport set new custom transport
func (c *Client) SetCustomTransport(customHTTPTransport http.RoundTripper) {
	if c.httpClient != nil {
		c.httpClient.Transport = customHTTPTransport
	}
}

func (c *Client) Post(path string, query interface{}, v interface{}, contentType ...string) (res *http.Response, err error) {
	body := new(bytes.Buffer)

	ct := "application/json"
	if len(contentType) > 0 {
		ct = contentType[0]
	}

	if r, ok := v.(io.Reader); ok {
		return c.do(http.MethodPost, path, query, r, ct)
	} else {
		if err := json.NewEncoder(body).Encode(v); err != nil {
			return nil, err
		}
		res, err = c.do(http.MethodPost, path, query, body, ct)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func (c *Client) doPutJSON(path string, query map[string]string, v interface{}) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(v); err != nil {
		return err
	}

	_, err := c.do(http.MethodPut, path, query, body, "application/json")
	return err
}

func (c *Client) Delete(path string, query interface{}) (res *http.Response, err error) {
	return c.do(http.MethodDelete, path, query, nil, "")
}

func (c *Client) doPost(path string, query interface{}) (res *http.Response, err error) {
	return c.do(http.MethodPost, path, query, nil, "")
}

func (c *Client) doPut(path string, query map[string]string) (res *http.Response, err error) {
	return c.do(http.MethodPut, path, query, nil, "")
}

func (c *Client) do(method, path string, q interface{}, body io.Reader, contentType string) (res *http.Response, err error) {
	u, err := c.buildURL(path, q)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	req.SetBasicAuth(c.apiUser, c.apiPassword)

	res, err = c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err = c.checkResponse(res); err != nil {
		return nil, err
	}

	return
}

func (c *Client) Get(path string, query interface{}) (res *http.Response, err error) {
	return c.do(http.MethodGet, path, query, nil, "")
}

func (c *Client) checkResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}

	defer res.Body.Close()

	if res.Header.Get("Content-Type") == "application/json" {
		if res.StatusCode == 404 {
			return ErrorNotFound
		}

		jsonErr := &Error{}
		err := json.NewDecoder(res.Body).Decode(jsonErr)
		if err != nil {
			return fmt.Errorf("response error with status code %d: failed unmarshal error response: %w", res.StatusCode, err)
		}

		return jsonErr
	}

	errText, err := ioutil.ReadAll(res.Body)
	if err == nil {
		return fmt.Errorf("response error with status code %d: %s", res.StatusCode, string(errText))
	}

	return fmt.Errorf("response error with status code %d", res.StatusCode)
}

func (c *Client) Marshal(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) buildURL(path string, q interface{}) (string, error) {
	// TODO: full refactor to use hard typed interfaces
	if q != nil && reflect.ValueOf(q).Kind() == reflect.Map {
		bb, _ := json.Marshal(q)
		log.Debug().Caller().Str("path", path).
			RawJSON("map", bb).
			Msg("Deprecated map query usage. Use struct with query tags instead!")

		m, ok := q.(map[string]string)
		if !ok {
			return "", errors.New("cannot convert query to map[string]string")
		}

		if len(m) == 0 {
			return c.endpointURL + path, nil
		}

		u, err := url.Parse(c.endpointURL + path)
		if err != nil {
			return "", err
		}
		q := u.Query()

		for k, v := range m {
			q.Set(k, v)
		}

		u.RawQuery = q.Encode()

		return u.String(), nil
	}

	// Mapping the interface with google's query tool into raw query string
	v, err := query.Values(q)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(c.endpointURL + path)
	if err != nil {
		return "", err
	}

	u.RawQuery = v.Encode()

	// log.Debug().Str("url", u.String()).Msg("URL built")

	return u.String(), nil
}
