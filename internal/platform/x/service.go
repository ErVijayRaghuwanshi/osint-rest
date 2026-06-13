package x

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"osint-scraper/internal/header"
	"osint-scraper/internal/httpclient"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Service struct {
	*httpclient.BaseService
}

func NewService(hm *header.Manager, log zerolog.Logger) *Service {
	return &Service{
		BaseService: httpclient.NewBaseService("x", "https://www.x.com", hm, log),
	}
}

// Name returns the platform identifier.
func (s *Service) Name() string {
	return "x"
}

// Ping checks if the platform website is reachable.
func (s *Service) Ping(ctx context.Context) bool {
	return s.BaseService.CheckWebsite(ctx)
}

// RegisterRoutes registers the routes for this platform.
func (s *Service) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, s, s.Log)
}

// GetUserInfo fetches user information from X by screen name
func (s *Service) GetUserInfo(ctx context.Context, username string) ([]byte, error) {
	baseURL := "https://x.com/i/api/graphql/-oaLodhGbbnzJBACb1kk2Q/UserByScreenName"

	q := url.Values{}
	q.Set("variables", mustJSON(map[string]interface{}{
		"screen_name":           username,
		"withGrokTranslatedBio": false,
	}))
	q.Set("features", mustJSON(map[string]bool{
		"hidden_profile_subscriptions_enabled":                              true,
		"profile_label_improvements_pcf_label_in_post_enabled":              true,
		"responsive_web_profile_redirect_enabled":                           false,
		"rweb_tipjar_consumption_enabled":                                   true,
		"verified_phone_label_enabled":                                      true,
		"subscriptions_verification_info_is_identity_verified_enabled":      true,
		"subscriptions_verification_info_verified_since_enabled":            true,
		"highlights_tweets_tab_ui_enabled":                                  true,
		"responsive_web_twitter_article_notes_tab_enabled":                  true,
		"subscriptions_feature_can_gift_premium":                            true,
		"creator_subscriptions_tweet_preview_api_enabled":                   true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled": false,
		"responsive_web_graphql_timeline_navigation_enabled":                true,
	}))
	q.Set("fieldToggles", mustJSON(map[string]bool{
		"withPayments":            false,
		"withAuxiliaryUserLabels": true,
	}))

	endpoint := baseURL + "?" + q.Encode()

	return s.WithRetry(ctx, 3, func(session *httpclient.Session, headers map[string]string, headerID string) ([]byte, error) {
		resp, err := session.Get(ctx, endpoint, headers)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == http.StatusOK {
			return bodyBytes, nil
		}

		s.Log.Warn().Int("status", resp.StatusCode).Str("header_id", headerID).Msg("X non-OK response")

		// Retry only on 403
		if resp.StatusCode != http.StatusForbidden {
			return nil, fmt.Errorf("x api error %d: %s", resp.StatusCode, bodyBytes)
		}

		return nil, fmt.Errorf("403 forbidden")
	})
}

func mustJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
