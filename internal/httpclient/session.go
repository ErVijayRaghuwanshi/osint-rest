package httpclient

import (
	"context"
	"io"
	"maps"
	"net/http"
	"net/http/cookiejar"
	"time"
)

// Session wraps an http.Client with cookie jar and header management.
// It is the single shared implementation used by all platform services.
type Session struct {
	Client  *http.Client
	Jar     *cookiejar.Jar
	Headers map[string]string
}

// NewSession creates a new HTTP session with cookie jar and sensible defaults.
func NewSession() (*Session, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	return &Session{
		Client:  client,
		Jar:     jar,
		Headers: make(map[string]string),
	}, nil
}

// SetAuthHeaders merges the given headers into the session-level headers.
func (s *Session) SetAuthHeaders(headers map[string]string) {
	if headers == nil {
		return
	}
	maps.Copy(s.Headers, headers)
}

// DoRequest performs an HTTP request with session cookies and headers.
// It accepts a context for cancellation/timeout propagation.
func (s *Session) DoRequest(
	ctx context.Context,
	method, url string,
	headers map[string]string,
	body io.Reader,
) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Default headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("Connection", "keep-alive")

	// Session-level headers (auth, csrf, etc.)
	for k, v := range s.Headers {
		req.Header.Set(k, v)
	}

	// Per-request headers (override session headers if needed)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return s.Client.Do(req)
}

// Get performs a GET request with context.
func (s *Session) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return s.DoRequest(ctx, "GET", url, headers, nil)
}

// Post performs a POST request with context.
func (s *Session) Post(ctx context.Context, url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	return s.DoRequest(ctx, "POST", url, headers, body)
}

// Close closes idle connections on the underlying client.
func (s *Session) Close() error {
	s.Client.CloseIdleConnections()
	return nil
}
