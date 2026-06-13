package snapchat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"osint-scraper/internal/header"
	"osint-scraper/internal/httpclient"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Service struct {
	*httpclient.BaseService
}

func NewService(hm *header.Manager, log zerolog.Logger) *Service {
	return &Service{
		BaseService: httpclient.NewBaseService("snapchat", "https://www.snapchat.com", hm, log),
	}
}

// Name returns the platform identifier.
func (s *Service) Name() string {
	return "snapchat"
}

// Ping checks if the platform website is reachable.
func (s *Service) Ping(ctx context.Context) bool {
	return s.BaseService.CheckWebsite(ctx)
}

// RegisterRoutes registers the routes for this platform.
func (s *Service) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, s, s.Log)
}

// getNestedMap safely traverses nested maps
func getNestedMap(m map[string]interface{}, keys ...string) (map[string]interface{}, bool) {
	curr := m
	for _, key := range keys {
		next, ok := curr[key].(map[string]interface{})
		if !ok {
			return nil, false
		}
		curr = next
	}
	return curr, true
}

// GetNextData fetches and parses __NEXT_DATA__ from a Snapchat user page.
// This is the single source of truth — all other methods reuse it.
func (s *Service) GetNextData(ctx context.Context, username string) ([]byte, error) {
	session, err := s.EnsureSession()
	if err != nil {
		return nil, err
	}

	url := "https://www.snapchat.com/@" + username

	resp, err := session.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetNextData failed: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	script := doc.Find("script#__NEXT_DATA__").First()
	if script.Length() == 0 {
		return nil, fmt.Errorf("__NEXT_DATA__ script not found")
	}

	jsonText := script.Text()
	if jsonText == "" {
		return nil, fmt.Errorf("__NEXT_DATA__ is empty")
	}

	// Validate JSON and check pageMetadata exists
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonText), &raw); err != nil {
		return nil, err
	}

	if _, ok := getNestedMap(raw, "props", "pageProps", "pageMetadata"); !ok {
		return nil, fmt.Errorf("pageMetadata not found")
	}

	return []byte(jsonText), nil
}

// GetUserInfo fetches user information by reusing GetNextData
func (s *Service) GetUserInfo(ctx context.Context, username string) (*UserInfoResponse, int, error) {
	jsonBytes, err := s.GetNextData(ctx, username)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var nextData struct {
		Props struct {
			PageProps struct {
				UserProfile UserInfoResponse `json:"userProfile"`
			} `json:"pageProps"`
		} `json:"props"`
	}

	if err := json.Unmarshal(jsonBytes, &nextData); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return &nextData.Props.PageProps.UserProfile, http.StatusOK, nil
}

// GetCuratedHighlights fetches curated highlights by reusing GetNextData
func (s *Service) GetCuratedHighlights(ctx context.Context, username string) ([]CuratedHighlight, int, error) {
	jsonBytes, err := s.GetNextData(ctx, username)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var nextData struct {
		Props struct {
			PageProps struct {
				Highlights struct {
					CuratedHighlights []CuratedHighlight `json:"curatedHighlights"`
				}
			} `json:"pageProps"`
		} `json:"props"`
	}

	if err := json.Unmarshal(jsonBytes, &nextData); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nextData.Props.PageProps.Highlights.CuratedHighlights, http.StatusOK, nil
}

// GetSpotlightHighlights fetches spotlight highlights by reusing GetNextData
func (s *Service) GetSpotlightHighlights(ctx context.Context, username string) ([]SpotlightHighlight, int, error) {
	jsonBytes, err := s.GetNextData(ctx, username)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var nextData struct {
		Props struct {
			PageProps struct {
				SpotlightHighlights []SpotlightHighlight `json:"spotlightHighlights"`
			} `json:"pageProps"`
		} `json:"props"`
	}

	if err := json.Unmarshal(jsonBytes, &nextData); err != nil {
		return nil, http.StatusBadRequest, err
	}

	return nextData.Props.PageProps.SpotlightHighlights, http.StatusOK, nil
}
