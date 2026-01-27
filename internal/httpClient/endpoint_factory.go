package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	dialTimeout    = 1 * time.Second
	requestTimeout = 30 * time.Second
	poolName       = "fintech_sre_client"
)

// EndpointFactory builds HTTP endpoints for a base URL.
type EndpointFactory interface {
	Build(pattern string) Endpoint
}

// Endpoint performs HTTP requests. Path params like {id} are replaced via WithParam.
type Endpoint interface {
	Get(ctx context.Context, opts ...RequestOption) (*http.Response, error)
	Post(ctx context.Context, opts ...RequestOption) (*http.Response, error)
	Patch(ctx context.Context, opts ...RequestOption) (*http.Response, error)
}

// RequestOption configures a request (path params, query, body).
type RequestOption interface {
	apply(*requestConfig)
}

type requestConfig struct {
	pathParams map[string]string
	query      url.Values
	body       interface{}
}

type paramOpt struct{ k, v string }
func (o paramOpt) apply(c *requestConfig) { c.pathParams[o.k] = o.v }

type queryOpt struct{ v url.Values }
func (o queryOpt) apply(c *requestConfig) { c.query = o.v }

type bodyOpt struct{ v interface{} }
func (o bodyOpt) apply(c *requestConfig) { c.body = o.v }

// WithParam sets a path parameter (e.g. "id" -> "123" for pattern "/v1/accounts/{id}").
func WithParam(key, value string) RequestOption { return paramOpt{k: key, v: value} }

// WithQuery sets URL query values.
func WithQuery(v url.Values) RequestOption { return queryOpt{v: v} }

// WithBody sets the JSON body for POST/PATCH.
func WithBody(v interface{}) RequestOption { return bodyOpt{v: v} }

var _ EndpointFactory = (*DefaultEndpointFactory)(nil)

// NewEndpointFactory creates a factory for the given base URL.
func NewEndpointFactory(baseURL string) *DefaultEndpointFactory {
	return &DefaultEndpointFactory{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		client: &http.Client{
			Timeout: requestTimeout,
			Transport: &http.Transport{
				IdleConnTimeout: dialTimeout,
			},
		},
	}
}

// DefaultEndpointFactory implements EndpointFactory using net/http.
type DefaultEndpointFactory struct {
	baseURL string
	client  *http.Client
}

// Build returns an Endpoint for baseURL + pattern. Pattern may contain placeholders like {id}.
func (f *DefaultEndpointFactory) Build(pattern string) Endpoint {
	return &defaultEndpoint{
		baseURL: f.baseURL,
		pattern: strings.TrimPrefix(pattern, "/"),
		client:  f.client,
	}
}

type defaultEndpoint struct {
	baseURL string
	pattern string
	client  *http.Client
}

func (e *defaultEndpoint) urlAndConfig(opts []RequestOption) (string, *requestConfig, error) {
	cfg := &requestConfig{pathParams: make(map[string]string)}
	for _, o := range opts {
		o.apply(cfg)
	}
	path := e.pattern
	for k, v := range cfg.pathParams {
		path = strings.ReplaceAll(path, "{"+k+"}", url.PathEscape(v))
	}
	u := e.baseURL + "/" + path
	if cfg.query != nil && len(cfg.query) > 0 {
		u += "?" + cfg.query.Encode()
	}
	return u, cfg, nil
}

func (e *defaultEndpoint) Get(ctx context.Context, opts ...RequestOption) (*http.Response, error) {
	u, _, err := e.urlAndConfig(opts)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return e.client.Do(req)
}

func (e *defaultEndpoint) Post(ctx context.Context, opts ...RequestOption) (*http.Response, error) {
	u, cfg, err := e.urlAndConfig(opts)
	if err != nil {
		return nil, err
	}
	var body []byte
	if cfg.body != nil {
		body, err = json.Marshal(cfg.body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return e.client.Do(req)
}

func (e *defaultEndpoint) Patch(ctx context.Context, opts ...RequestOption) (*http.Response, error) {
	u, cfg, err := e.urlAndConfig(opts)
	if err != nil {
		return nil, err
	}
	var body []byte
	if cfg.body != nil {
		body, err = json.Marshal(cfg.body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return e.client.Do(req)
}

// HTTPError represents a non-2xx response.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d: %s", e.StatusCode, e.Body)
}
