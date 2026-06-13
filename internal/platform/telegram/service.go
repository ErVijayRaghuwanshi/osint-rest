package telegram

import (
	"context"
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
		BaseService: httpclient.NewBaseService("telegram", "https://www.telegram.com", hm, log),
	}
}

// Name returns the platform identifier.
func (s *Service) Name() string {
	return "telegram"
}

// Ping checks if the platform website is reachable.
func (s *Service) Ping(ctx context.Context) bool {
	return s.BaseService.CheckWebsite(ctx)
}

// RegisterRoutes registers the routes for this platform.
func (s *Service) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, s, s.Log)
}

// GetUserInfo fetches user information from Telegram
func (s *Service) GetUserInfo(ctx context.Context, username string) ([]byte, error) {
	endpoint := "https://www.telegram.com/api/userinfo?username=" +
		url.QueryEscape(username)

	return s.WithRetry(ctx, 3, func(session *httpclient.Session, headers map[string]string, headerID string) ([]byte, error) {
		resp, err := session.Get(ctx, endpoint, headers)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("telegram read body: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return body, nil
		}

		s.Log.Warn().Int("status", resp.StatusCode).Str("header_id", headerID).Msg("Telegram non-OK response")

		if resp.StatusCode < 400 || resp.StatusCode >= 500 {
			return nil, fmt.Errorf("telegram error: %d", resp.StatusCode)
		}

		return nil, fmt.Errorf("telegram %d forbidden", resp.StatusCode)
	})
}
