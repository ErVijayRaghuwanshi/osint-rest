package instagram

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
// @Summary      Check Instagram availability
// @Description  Pings instagram.com to ensure site is reachable
// @Tags         Instagram
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      503  {object}  map[string]string
// @Router       /api/instagram/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite()

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
	fmt.Println("Requested Instagram username:", username)
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username parameter is required"})
		return
	}

	data, err := h.svc.GetUserInfo(username)
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

	data, err := h.svc.GetTimeline(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, result)
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
	fmt.Println("Requested Instagram hashtag search query:", query)
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	data, err := h.svc.SearchHashtags(query)
	fmt.Println("Hashtag search raw response data:", string(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// var hashtagResp HashTagSearchResponse
	// if err := json.Unmarshal(data, &hashtagResp); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
	// 	return
	// }

	// c.JSON(http.StatusOK, hashtagResp)

	// just return raw data for now
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("Error unmarshaling hashtag search response:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse response"})
		return
	}

	c.JSON(http.StatusOK, result)

}