package snapchat

import (
	"fmt"
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
// @Summary      Check Snapchat availability
// @Description  Pings snapchat.com to ensure site is reachable
// @Tags         Snapchat
// @Produce      json
// @Success      200  {object}  PingResponse
// @Router       /api/snapchat/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite(c.Request.Context())

	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "snapchat pong"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "snapchat unreachable"})
	}
}

// GetUserInfo godoc
// @Summary      Get Snapchat user information
// @Description  Fetches detailed user information for a given Snapchat username
// @Tags         Snapchat
// @Produce      json
// @Param        username  query     string  false "Snapchat username (default: arora_girl)"
// @Success      200       {object}  UserInfoResponse
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/snapchat/userinfo [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	username := c.Query("username")
	h.log.Info().Str("username", username).Msg("Requested Snapchat user info")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	userInfo, statusCode, err := h.svc.GetUserInfo(c.Request.Context(), username)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": fmt.Sprintf("%v for username: %s", err.Error(), username)})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}

// GetCuratedHighlights godoc
// @Summary      Get Snapchat user curated highlights
// @Description  Fetches curated highlights for a given Snapchat username
// @Tags         Snapchat
// @Produce      json
// @Param        username  query     string  false "Snapchat username (default: arora_girl)"
// @Success      200       {object}  CuratedHighlightsResponse
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/snapchat/curatedhighlights [get]
func (h *Handler) GetCuratedHighlights(c *gin.Context) {
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	highlights, statusCode, err := h.svc.GetCuratedHighlights(c.Request.Context(), username)
	h.log.Info().Int("count", len(highlights)).Str("username", username).Msg("Curated highlights fetched")
	if err != nil {
		c.JSON(statusCode, gin.H{"error": fmt.Sprintf("%v for username: %s", err.Error(), username)})
		return
	}

	c.JSON(http.StatusOK, highlights)
}

// GetSpotlightHighlights godoc
// @Summary      Get Snapchat user spotlight highlights
// @Description  Fetches spotlight highlights for a given Snapchat username
// @Tags         Snapchat
// @Produce      json
// @Param        username  query     string  false "Snapchat username (default : arora_girl)"
// @Success      200       {object}  []SpotlightHighlight
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/snapchat/spotlighthighlights [get]
func (h *Handler) GetSpotlightHighlights(c *gin.Context) {
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	highlights, statusCode, err := h.svc.GetSpotlightHighlights(c.Request.Context(), username)
	h.log.Info().Int("count", len(highlights)).Str("username", username).Msg("Spotlight highlights fetched")
	if err != nil {
		c.JSON(statusCode, gin.H{"error": fmt.Sprintf("%v for username: %s", err.Error(), username)})
		return
	}

	c.JSON(http.StatusOK, highlights)
}
