#!/usr/bin/env bash
set -euo pipefail

# -------------------------------------------------------
# scaffold-platform.sh — Generate boilerplate for a new
# OSINT Scraper platform package.
#
# Usage:  ./scripts/scaffold-platform.sh <platform_name>
# Example: ./scripts/scaffold-platform.sh telegram
# -------------------------------------------------------

if [ -z "${1:-}" ]; then
  echo "Usage: $0 <platform_name>"
  echo "  platform_name: lowercase name (e.g. telegram, tiktok)"
  exit 1
fi

NAME="$1"
# Title-cased name for comments/tags (e.g. telegram -> Telegram)
TITLE="$(echo "${NAME:0:1}" | tr '[:lower:]' '[:upper:]')${NAME:1}"

DIR="internal/api/${NAME}"

if [ -d "$DIR" ]; then
  echo "Error: $DIR already exists. Aborting."
  exit 1
fi

echo "Scaffolding platform: ${NAME} (${TITLE})"
mkdir -p "$DIR"

# -------------------------------------------------------
# models.go
# -------------------------------------------------------
cat > "${DIR}/models.go" <<EOF
package ${NAME}

type PingResponse struct {
	Message string \`json:"message" example:"Pong"\`
}

// UserInfoResponse represents the parsed user info from ${TITLE}
type UserInfoResponse struct {
	// TODO: Define response fields for ${TITLE}
}
EOF

# -------------------------------------------------------
# service.go
# -------------------------------------------------------
cat > "${DIR}/service.go" <<EOF
package ${NAME}

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
		BaseService: httpclient.NewBaseService("${NAME}", "https://www.${NAME}.com", hm, log),
	}
}

// CheckWebsite checks if ${TITLE} is reachable
func (s *Service) CheckWebsite(ctx context.Context) bool {
	return s.BaseService.CheckWebsite(ctx)
}

// GetUserInfo fetches user information from ${TITLE}
func (s *Service) GetUserInfo(ctx context.Context, username string) ([]byte, error) {
	endpoint := "https://www.${NAME}.com/api/userinfo?username=" +
		url.QueryEscape(username)

	return s.WithRetry(ctx, 3, func(session *httpclient.Session, headers map[string]string, headerID string) ([]byte, error) {
		resp, err := session.Get(ctx, endpoint, headers)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("${NAME} read body: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return body, nil
		}

		s.Log.Warn().Int("status", resp.StatusCode).Str("header_id", headerID).Msg("${TITLE} non-OK response")

		if resp.StatusCode < 400 || resp.StatusCode >= 500 {
			return nil, fmt.Errorf("${NAME} error: %d", resp.StatusCode)
		}

		return nil, fmt.Errorf("${NAME} %d forbidden", resp.StatusCode)
	})
}
EOF

# -------------------------------------------------------
# handlers.go
# -------------------------------------------------------
cat > "${DIR}/handlers.go" <<EOF
package ${NAME}

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	svc *Service
	log zerolog.Logger
}

func NewHandler(svc *Service, log zerolog.Logger) *Handler {
	return &Handler{svc: svc, log: log}
}

// Ping godoc
// @Summary      Check ${TITLE} availability
// @Description  Pings ${NAME} to ensure site is reachable
// @Tags         ${TITLE}
// @Produce      json
// @Success      200  {object}  PingResponse
// @Router       /api/${NAME}/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite(c.Request.Context())

	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "${NAME} pong"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "${NAME} unreachable"})
	}
}

// GetUserInfo godoc
// @Summary      Get ${TITLE} user information
// @Description  Fetches detailed user information for a given ${TITLE} username
// @Tags         ${TITLE}
// @Produce      json
// @Param        username  query     string  true  "${TITLE} username"
// @Success      200       {object}  UserInfoResponse
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/${NAME}/userinfo [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	username := c.Query("username")
	h.log.Info().Str("username", username).Msg("Requested ${TITLE} user info")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	data, err := h.svc.GetUserInfo(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp UserInfoResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
EOF

# -------------------------------------------------------
# routes.go
# -------------------------------------------------------
cat > "${DIR}/routes.go" <<EOF
package ${NAME}

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RegisterRoutes(rg *gin.RouterGroup, svc *Service, log zerolog.Logger) {
	h := NewHandler(svc, log)

	rg.GET("/ping", h.Ping)
	rg.GET("/userinfo", h.GetUserInfo)
}
EOF

echo ""
echo "✅ Scaffolded ${TITLE} platform at ${DIR}/"
echo "   - handlers.go  (Ping + GetUserInfo)"
echo "   - service.go   (CheckWebsite + GetUserInfo with retry)"
echo "   - models.go    (PingResponse + UserInfoResponse stub)"
echo "   - routes.go    (RegisterRoutes)"
echo ""
echo "Next steps:"
echo "  1. Update models.go with actual response structs"
echo "  2. Update service.go with the correct base URL and API endpoint"
echo "  3. Wire the platform in internal/api/router.go:"
echo ""
echo "     import \"osint-scraper/internal/api/${NAME}\""
echo ""
echo "     ${NAME}Service := ${NAME}.NewService(hm, log)"
echo "     ${NAME}Service.SetCache(appCache)"
echo "     ${NAME}Service.SetHealthTracker(ht)"
echo "     ${NAME}Group := r.Group(\"/api/${NAME}\")"
echo "     ${NAME}.RegisterRoutes(${NAME}Group, ${NAME}Service, log)"
echo ""
echo "  4. Add headers to config/headers.json under \"${NAME}\" platform"
echo "  5. Run 'make run' to verify"
