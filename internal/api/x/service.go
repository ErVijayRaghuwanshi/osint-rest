package x

import (
	"maps"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Session struct {
	client *http.Client
	jar    *cookiejar.Jar
	headers map[string]string // persistent headers
}

type Service struct {
	session *Session
}

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
		headers: make(map[string]string), // IMPORTANT: init map
	}, nil
}

// ------------------------------------------------------------
// STORE HEADERS INTO SESSION
// ------------------------------------------------------------
func (s *Session) SetAuthHeaders(headers map[string]string) {
	maps.Copy(s.headers, headers)
}

// ------------------------------------------------------------
// APPLY SESSION HEADERS + PER REQUEST OVERRIDES
// ------------------------------------------------------------
func (s *Session) DoRequest(method, url string, override map[string]string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// 1️⃣ Apply persistent session headers
	for k, v := range s.headers {
		req.Header.Set(k, v)
	}

	// 2️⃣ Apply temporary override headers (if provided)
	if override != nil {
		for k, v := range override {
			req.Header.Set(k, v)
		}
	}

	return s.client.Do(req)
}

func (s *Session) Get(url string, h map[string]string) (*http.Response, error) {
	return s.DoRequest("GET", url, h, nil)
}

func (s *Session) Post(url string, h map[string]string, body io.Reader) (*http.Response, error) {
	return s.DoRequest("POST", url, h, body)
}

// ------------------------------------------------------------
// SERVICE
// ------------------------------------------------------------

func NewService() *Service {
	return &Service{}
}

func (s *Service) InitializeSession() (*Session, error) {
	session, err := NewSession()
	if err != nil {
		return nil, err
	}

	// Set persistent global headers ONCE
	session.SetAuthHeaders(map[string]string{
		"Accept":     "*/*",
		"Authorization": "Bearer AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA",
		"x-csrf-token": "07a7bd632f372f3665a80d538397d12d200c97894749e8bcb81dbd8e2b7c8caef838f786bc12bcf894b6a0192fef52099396f7a70dc3a65a13a4b49c0d255cdb7ca4bf1768333f725098fc508cbd92e0",
		"Cookie": "night_mode=2; personalization_id=\"v1_ghCHIePjBETcby5cGISyPg==\"; kdt=pLnJLTcX5PJKegezUlAMKUMSAGfR1fAf1hKsF8rT; auth_token=2c7689c7b2bb354954e519d60dc6dfdd49267e6c; ct0=07a7bd632f372f3665a80d538397d12d200c97894749e8bcb81dbd8e2b7c8caef838f786bc12bcf894b6a0192fef52099396f7a70dc3a65a13a4b49c0d255cdb7ca4bf1768333f725098fc508cbd92e0; twid=u%3D1885901238673477632; lang=en;",
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36",
		"Accept-Language": "en-US,en;q=0.9",
	})

	s.session = session
	return session, nil
}

func (s *Service) GetUserInfo(screenName string) ([]byte, error) {
	if s.session == nil {
		_, err := s.InitializeSession()
		if err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf(
		"https://x.com/i/api/graphql/-oaLodhGbbnzJBACb1kk2Q/UserByScreenName?variables={\"screen_name\":\"%s\"}",
		screenName,
	)

	resp, err := s.session.Get(url, nil) // NO HEADERS NEEDED
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("failed to fetch user info, status code: %d", resp.StatusCode)
    }
    // DEBUG: status code
	fmt.Println("Status code:", resp.StatusCode)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (s *Service) CheckWebsite() bool {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://www.x.com")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
