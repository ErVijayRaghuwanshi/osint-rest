package instagram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"osint-scraper/internal/header"
	"osint-scraper/internal/httpclient"

	"github.com/rs/zerolog"
)

type Service struct {
	*httpclient.BaseService
}

func NewService(hm *header.Manager, log zerolog.Logger) *Service {
	return &Service{
		BaseService: httpclient.NewBaseService("instagram", "https://www.instagram.com", hm, log),
	}
}

// CheckWebsite checks if Instagram is reachable
func (s *Service) CheckWebsite(ctx context.Context) bool {
	return s.BaseService.CheckWebsite(ctx)
}

// GetUserInfo fetches user information using header rotation with retries
func (s *Service) GetUserInfo(ctx context.Context, username string) ([]byte, error) {
	endpoint := "https://www.instagram.com/api/v1/users/web_profile_info/?username=" +
		url.QueryEscape(username)

	return s.WithRetry(ctx, 3, func(session *httpclient.Session, headers map[string]string, headerID string) ([]byte, error) {
		resp, err := session.Get(ctx, endpoint, headers)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return io.ReadAll(resp.Body)
		}

		s.Log.Warn().Int("status", resp.StatusCode).Str("header_id", headerID).Msg("Instagram non-OK response")

		// retry on 4xx (rate-limit / forbidden), fail on 5xx
		if resp.StatusCode < 400 || resp.StatusCode >= 500 {
			return nil, fmt.Errorf("instagram error: %d", resp.StatusCode)
		}

		return nil, fmt.Errorf("instagram %d forbidden", resp.StatusCode)
	})
}
