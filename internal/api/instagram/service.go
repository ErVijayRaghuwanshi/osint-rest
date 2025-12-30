package instagram

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	// "net/url"
	"net/http/cookiejar"
	"strings"
	"time"
)

type Session struct {
	client  *http.Client
	jar     *cookiejar.Jar
	headers map[string]string
}


type Service struct {
	session *Session
}

// NewSession creates a new authenticated Instagram session
// It maintains cookies and headers across multiple requests
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
		client:  client,
		jar:     jar,
		headers: make(map[string]string),
	}, nil
}


// SetAuthHeaders sets the authentication headers for the session
func (s *Session) SetAuthHeaders(headers map[string]string) {
	if headers == nil {
		return
	}
	for k, v := range headers {
		s.headers[k] = v
	}
}


// DoRequest performs an HTTP request with the session's client and maintains cookies
func (s *Session) DoRequest(
	method, url string,
	headers map[string]string,
	body io.Reader,
) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Default headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
	req.Header.Set("x-csrftoken", "PUesp9RgFVHZJ9Ejx-j6wT")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("x-ig-app-id", "936619743392459")



	// Session-level headers (auth, csrf, etc.)
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	// Per-request headers (override session headers if needed)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return s.client.Do(req)
}


// Get performs a GET request
func (s *Session) Get(url string, headers map[string]string) (*http.Response, error) {
	return s.DoRequest("GET", url, headers, nil)
}

// Post performs a POST request
func (s *Session) Post(url string, headers map[string]string, body io.Reader) (*http.Response, error) {
	return s.DoRequest("POST", url, headers, body)
}

// Close closes the session
func (s *Session) Close() error {
	s.client.CloseIdleConnections()
	return nil
}

func NewService() *Service {
	return &Service{}
}

// InitializeSession creates and initializes an authenticated session
func (s *Service) InitializeSession() (*Session, error) {
	session, err := NewSession()
	if err != nil {
		return nil, err
	}
	s.session = session
	return session, nil
}

// GetSession returns the current session
func (s *Service) GetSession() *Session {
	return s.session
}

func (s *Service) CheckWebsite() bool {
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get("https://www.instagram.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200
}

// GetUserInfo fetches user information using the authenticated session
func (s *Service) GetUserInfo(username string) ([]byte, error) {
	if s.session == nil {
		_, err := s.InitializeSession()
		if err != nil {
			return nil, err
		}
	}

	url := "https://www.instagram.com/api/v1/users/web_profile_info/?username=" + username

	// Only endpoint-specific headers (optional)
	headers := map[string]string{
		"Accept": "application/json",
	}

	resp, err := s.session.Get(url, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetUserInfo failed: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}


// GetTimeline fetches the user's timeline using the authenticated session
func (s *Service) GetTimeline(username string) ([]byte, error) {
	if s.session == nil {
		session, err := s.InitializeSession()
		if err != nil {
			return nil, err
		}
		s.session = session
	}

	url := "https://www.instagram.com/api/v1/feed/user/" + username + "/username/?count=12"

	headers := map[string]string{
		"Accept":           "application/json",
		"X-Requested-With": "XMLHttpRequest",
	}

	resp, err := s.session.Get(url, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// SearchHashtags searches for hashtags using the authenticated session
func (s *Service) SearchHashtags(query string) ([]byte, error) {
	if s.session == nil {
		_, err := s.InitializeSession()
		if err != nil {
			return nil, err
		}
	}

	endpoint := "https://www.instagram.com/graphql/query"

	payload := fmt.Sprintf(
		`av=17841409322006652&__d=www&__user=0&__a=1&__req=1d&__hs=20451.HCSV2%%3Ainstagram_web_pkg.2.1...0&dpr=2&__ccg=MODERATE&__rev=1031507717&__s=6wva3i%%3Afmlmkn%%3An0n4yp&__hsi=7589342574035025574&__dyn=7xeUjG1mxu1syaxG4Vp41twWwIxu13wvoKewSAwHwNw9G2Saxa0DU6u3y4o0B-q1ew6ywMwto2awgo9oO0n24oaEnxO1ywOwv89k2C1Fwc60D85m263ifK0zEkxe2GewGw9a361qwuEjUlwhEe88o5i7U1oEbUGdG1QwTU9UaQ0z8c86-bwHwKG1pg2fwxyo6O1FwlA3a3zhAq4rwIDyXxui2qiUqwm8jxK2K2G0EoKmUhw4rxOi6oGq2K18whE984O&__comet_req=7&server_timestamps=true&variables=%%7B%%22data%%22%%3A%%7B%%22context%%22%%3A%%22blended%%22%%2C%%22include_reel%%22%%3A%%22true%%22%%2C%%22query%%22%%3A%%22%s%%22%%2C%%22rank_token%%22%%3A%%22%%22%%2C%%22search_session_id%%22%%3A%%22%%22%%2C%%22search_surface%%22%%3A%%22web_top_search%%22%%7D%%2C%%22hasQuery%%22%%3Atrue%%7D&doc_id=24146980661639222`,
		url.QueryEscape(query),
	)

	resp, err := s.session.Post(endpoint, nil, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SearchHashtags failed: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
