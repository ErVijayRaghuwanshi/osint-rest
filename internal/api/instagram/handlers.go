package instagram

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
// @Summary      Check Instagram availability
// @Description  Pings instagram.com to ensure site is reachable
// @Tags         Instagram
// @Produce      json
// @Success      200  {object}  PingResponse
// @Router       /api/instagram/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite(c.Request.Context())

	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "instagram pong"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "instagram unreachable"})
	}
}

// GetUserInfo godoc
// @Summary      Get Instagram user information
// @Description  Fetches detailed user information for a given Instagram username
// @Tags         Instagram
// @Produce      json
// @Param        username  query     string  false "Instagram username (default: sakshi_raghu_1c_)"
// @Success      200       {object}  UserInfoResponse
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/instagram/userinfo [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	username := c.Query("username")
	h.log.Info().Str("username", username).Msg("Requested Instagram user info")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	data, err := h.svc.GetUserInfo(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userInfoResp UserInfoResponse
	if err := json.Unmarshal(data, &userInfoResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, userInfoResp)
}

// GetTimeline godoc
// @Summary      Get Instagram user timeline
// @Description  Fetches timeline/posts for a given Instagram username
// @Tags         Instagram
// @Produce      json
// @Param        username  query     string  true  "Instagram username"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/instagram/timeline [get]
func (h *Handler) GetTimeline(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	_ = username
	c.JSON(http.StatusNotImplemented, gin.H{"error": "GetTimeline not implemented"})
}

// SearchHashtags godoc
// @Summary      Search Instagram hashtags
// @Description  Searches for Instagram hashtags based on a query string
// @Tags         Instagram
// @Produce      json
// @Param        query  query     string  true  "Search query for hashtags"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/instagram/hashtags [get]
func (h *Handler) SearchHashtags(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	_ = query
	c.JSON(http.StatusNotImplemented, gin.H{"error": "SearchHashtags not implemented"})
}
