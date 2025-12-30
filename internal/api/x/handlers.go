package x

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Ping godoc
// @Summary      Check X availability
// @Description  Pings x.com to ensure site is reachable
// @Tags         X
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      503  {object}  map[string]string
// @Router       /api/x/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite()

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
// @Param        screen_name  query     string  false  "X screen name"
// @Success      200          {object}  UserInfo
// @Failure      400          {object}  map[string]string
// @Failure      404          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /api/x/userinfo [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	screenName := c.Query("screen_name")
	if screenName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "screen_name parameter is required"})
		return
	}

	data, err := h.svc.GetUserInfo(screenName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Raw user info data:", string(data))
	var userInfoResp UserInfo
	if err := json.Unmarshal(data, &userInfoResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info"})
		return
	}
	
	c.JSON(http.StatusOK, userInfoResp)
}

