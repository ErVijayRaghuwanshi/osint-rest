package jaco

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"osint-scraper/internal/header"
	"osint-scraper/internal/httpclient"

	"github.com/rs/zerolog"
)

var ErrUserNotFound = errors.New("user not found")
var ErrCommentsNotFound = errors.New("comment not found")

type Service struct {
	*httpclient.BaseService
}

func NewService(hm *header.Manager, log zerolog.Logger) *Service {
	return &Service{
		BaseService: httpclient.NewBaseService("jaco", "https://jaco.live", hm, log),
	}
}

// CheckWebsite checks if Jaco is reachable
func (s *Service) CheckWebsite(ctx context.Context) bool {
	return s.BaseService.CheckWebsite(ctx)
}

// searchUser resolves a Jaco username to a user ID using the search API.
func (s *Service) searchUser(ctx context.Context, session *httpclient.Session, username string) (string, error) {
	searchURL := "https://api.jaco.live/search/c/offical/associate.json"
	q := url.Values{}
	q.Set("keyword", username)
	q.Set("show_agent", "0")

	resp, err := session.Get(ctx, searchURL+"?"+q.Encode(), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("search failed %d: %s", resp.StatusCode, body)
	}

	var searchResp struct {
		Data struct {
			List []struct {
				Content struct {
					User struct {
						ID       string `json:"id"`
						Username string `json:"username"`
					} `json:"user"`
				} `json:"content"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &searchResp); err != nil {
		return "", err
	}

	if len(searchResp.Data.List) == 0 {
		return "", ErrUserNotFound
	}

	user := searchResp.Data.List[0].Content.User
	if user.Username != username {
		return "", ErrUserNotFound
	}

	return user.ID, nil
}

// GetUserInfo fetches user information from Jaco by username
func (s *Service) GetUserInfo(ctx context.Context, username string) ([]byte, error) {
	return s.WithRetry(ctx, 5, func(session *httpclient.Session, headers map[string]string, headerID string) ([]byte, error) {
		uid, err := s.searchUser(ctx, session, username)
		if err != nil {
			return nil, err
		}

		profileURL := "https://api.jaco.live/bff/profile/header/user/web"
		q := url.Values{}
		q.Set("uid", uid)
		q.Set("show_badge", "1")
		q.Set("show_count", "1")
		q.Set("show_tab", "1")
		q.Set("show_business", "1")
		q.Set("show_live", "1")

		profileResp, err := session.Get(ctx, profileURL+"?"+q.Encode(), nil)
		if err != nil {
			return nil, err
		}
		defer profileResp.Body.Close()

		profileBody, _ := io.ReadAll(profileResp.Body)
		if profileResp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("profile fetch failed %d: %s", profileResp.StatusCode, profileBody)
		}

		return profileBody, nil
	})
}

// GetUserReelByUsername fetches a user's reel feed by username
func (s *Service) GetUserReelByUsername(ctx context.Context, username string, limit int, nextSinceID string) ([]byte, error) {
	s.Log.Info().Str("username", username).Int("limit", limit).Str("cursor", nextSinceID).Msg("GetUserReelByUsername")

	session, err := s.EnsureSession()
	if err != nil {
		return nil, err
	}

	uid, err := s.searchUser(ctx, session, username)
	if err != nil {
		return nil, err
	}

	s.Log.Debug().Str("uid", uid).Msg("Resolved user ID")

	reelURL := "https://api.jaco.live/reel/2/reel/c/flow.json"
	cursor := nextSinceID
	var finalRes map[string]any

	for {
		params := url.Values{}
		params.Set("target_uid", uid)
		if cursor != "" {
			params.Set("next_since_id", cursor)
		}

		endpoint := reelURL + "?" + params.Encode()

		resp, err := session.Get(ctx, endpoint, nil)
		if err != nil {
			return nil, err
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var page map[string]any
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, err
		}

		data := page["data"].(map[string]any)
		list := data["list"].([]any)

		if finalRes == nil {
			finalRes = page
		} else {
			finalData := finalRes["data"].(map[string]any)
			finalData["list"] = append(finalData["list"].([]any), list...)
		}

		finalData := finalRes["data"].(map[string]any)
		finalData["next_since_id"] = data["next_since_id"]
		finalData["pre_since_id"] = data["pre_since_id"]

		s.Log.Debug().Int("accumulated", len(finalData["list"].([]any))).Msg("Reel pagination")

		if len(finalData["list"].([]any)) >= limit {
			break
		}

		if fmt.Sprintf("%v", data["next_since_id"]) == "-1" && len(list) == 0 {
			s.Log.Debug().Msg("End of feed reached")
			break
		}

		cursor = fmt.Sprintf("%v", data["next_since_id"])
	}

	if finalRes != nil && limit > 0 {
		data := finalRes["data"].(map[string]any)
		list := data["list"].([]any)
		if len(list) > limit {
			data["list"] = list[:limit]
		}
	}

	return json.Marshal(finalRes)
}

// GetReelComments fetches comments for a reel with header rotation
func (s *Service) GetReelComments(ctx context.Context, rid string, limit int, baseID string) ([]byte, error) {
	return s.WithRetry(ctx, 5, func(session *httpclient.Session, headers map[string]string, headerID string) ([]byte, error) {
		const commentURL = "https://api.jaco.live/reel/1/comment/c/tree_root.json"

		if baseID == "" {
			baseID = "0"
		}

		var finalRes map[string]any
		cursor := baseID

		for {
			params := url.Values{}
			params.Set("rid", rid)
			params.Set("base_id", cursor)
			params.Set("type", "0")
			params.Set("max_id_type", "2")

			resp, err := session.Get(ctx, commentURL+"?"+params.Encode(), nil)
			if err != nil {
				return nil, err
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, ErrCommentsNotFound
			}

			var page map[string]any
			if err := json.Unmarshal(body, &page); err != nil {
				return nil, err
			}

			data, ok := page["data"].(map[string]any)
			if !ok {
				break
			}

			rawList, ok := data["list"]
			if !ok || rawList == nil {
				break
			}

			list, ok := rawList.([]any)
			if !ok {
				break
			}

			if finalRes == nil {
				finalRes = page
			} else {
				finalData := finalRes["data"].(map[string]any)
				finalData["list"] = append(finalData["list"].([]any), list...)
			}

			if len(list) == 0 {
				break
			}

			cursor = fmt.Sprintf("%v", data["since_id"])
			if cursor == "-1" {
				break
			}
		}

		return json.Marshal(finalRes)
	})
}
