package snapchat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Session struct {
    client *http.Client
    jar   *cookiejar.Jar
    headers map[string]string
}


type Service struct {
    session *Session
}

// NewSession creates a new authenticated Snapchat session
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
	req.Header.Set("Connection", "keep-alive")



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

// CheckWebsite checks if Snapchat website is reachable
func (s *Service) CheckWebsite() bool {
    client := http.Client{
        Timeout: 9 * time.Second,
    }

    resp, err := client.Get("https://www.snapchat.com")
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    return resp.StatusCode == 200
}


// GetUserInfo fetches user information using the authenticated session
func (s *Service) GetUserInfo(username string) (*UserInfoResponse, error) {
	if s.session == nil {
		_, err := s.InitializeSession()
		if err != nil {
			return nil, err
		}
	}

	url := "https://www.snapchat.com/@" + username

	// Perform request
	resp, err := s.session.Get(url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetUserInfo failed: %d", resp.StatusCode)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Find __NEXT_DATA__
	script := doc.Find("script#__NEXT_DATA__").First()
	if script.Length() == 0 {
		return nil, fmt.Errorf("__NEXT_DATA__ script not found")
	}

	jsonText := script.Text()
	if jsonText == "" {
		return nil, fmt.Errorf("__NEXT_DATA__ is empty")
	}

	// Decode only what we need
	var nextData struct {
		Props struct {
			PageProps struct {
				UserProfile UserInfoResponse `json:"userProfile"`
			} `json:"pageProps"`
		} `json:"props"`
	}

	if err := json.Unmarshal([]byte(jsonText), &nextData); err != nil {
		return nil, err
	}

	return &nextData.Props.PageProps.UserProfile, nil
}

