package x

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
// @Summary      Check X availability
// @Description  Pings x.com to ensure site is reachable
// @Tags         X
// @Produce      json
// @Success      200  {object}  PingResponse
// @Router       /api/x/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite(c.Request.Context())

	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "x pong"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "x unreachable"})
	}
}

// GetUserInfo godoc
// @Summary      Get X user information
// @Description  Fetches user information from X by screen name
// @Tags         X
// @Produce      json
// @Param        screen_name  query     string  false  "X screen name (default: urstrulymahesh)"
// @Success      200          {object}  UserInfoResponse
// @Failure      400          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /api/x/userinfo [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	username := c.Query("screen_name")
	h.log.Info().Str("screen_name", username).Msg("Requested X user info")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "screen_name parameter is required"})
		return
	}

	data, err := h.svc.GetUserInfo(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userInfoResp UserInfoResponse
	if err := json.Unmarshal(data, &userInfoResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info"})
		return
	}

	c.JSON(http.StatusOK, userInfoResp)
}
