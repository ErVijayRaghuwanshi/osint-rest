package telegram

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
// @Summary      Check Telegram availability
// @Description  Pings telegram to ensure site is reachable
// @Tags         Telegram
// @Produce      json
// @Success      200  {object}  PingResponse
// @Router       /api/telegram/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite(c.Request.Context())

	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "telegram pong"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "telegram unreachable"})
	}
}

// GetUserInfo godoc
// @Summary      Get Telegram user information
// @Description  Fetches detailed user information for a given Telegram username
// @Tags         Telegram
// @Produce      json
// @Param        username  query     string  true  "Telegram username"
// @Success      200       {object}  UserInfoResponse
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /api/telegram/userinfo [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	username := c.Query("username")
	h.log.Info().Str("username", username).Msg("Requested Telegram user info")

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
