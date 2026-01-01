package snapchat

import (
	// "encoding/json"
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
// @Summary      Check Snapchat availability
// @Description  Pings snapchat.com to ensure site is reachable
// @Tags         Snapchat
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      503  {object}  map[string]string
// @Router       /api/snapchat/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite()

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
	fmt.Println("Requested Snapchat username:", username)

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	userInfo, err := h.svc.GetUserInfo(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Gin automatically serializes structs to JSON
	c.JSON(http.StatusOK, userInfo)
}
