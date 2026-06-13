package jaco

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
// @Summary      Check Jaco availability
// @Description  Pings jaco.live to ensure site is reachable
// @Tags         Jaco
// @Produce      json
// @Success      200  {object}  PingResponse
// @Router       /api/jaco/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	ok := h.svc.CheckWebsite(c.Request.Context())
	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "jaco pong"})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "jaco unreachable"})
	}
}

// GetUserInfo godoc
// @Summary      Get Jaco user information
// @Description  Fetches user information from jaco by username
// @Tags         Jaco
// @Produce      json
// @Param        username  query     string  false  "Jaco username (default: Uaegirl)"
// @Success      200 {object} UserInfoResponse
// @Failure 400 {object} BadRequestError "Bad request"
// @Failure 404 {object} NotFoundError "User not found"
// @Failure 500 {object} InternalServerError "Internal server error"
// @Router       /api/jaco/userinfo [get]
func (h *Handler) GetUserInfo(c *gin.Context) {
	username := c.Query("username")
	h.log.Info().Str("username", username).Msg("Requested Jaco user info")

	if username == "" {
		c.JSON(http.StatusBadRequest, BadRequestError{
			Error: "username parameter is required",
		})
		return
	}

	data, err := h.svc.GetUserInfo(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, NotFoundError{
				Error: "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, InternalServerError{
			Error: "internal server error",
		})
		return
	}

	var userInfoResp UserInfoResponse
	if err := json.Unmarshal(data, &userInfoResp); err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError{
			Error: "failed to parse user info",
		})
		return
	}

	c.JSON(http.StatusOK, userInfoResp)
}

// GetUserReel godoc
// @Summary      Get Jaco user reel
// @Description  Fetches user reel (short videos) by username
// @Tags         Jaco
// @Produce      json
// @Param        username        query     string  true   "Jaco username (default: Amal16600)"
// @Param        limit           query     int     false  "Number of reels to fetch (default: 12)"
// @Param        next_since_id   query     string  false  "Pagination cursor"
// @Success      200 {object} UserReelResponse
// @Failure      400 {object} BadRequestError "Bad request"
// @Failure      404 {object} NotFoundError "User not found"
// @Failure      500 {object} InternalServerError "Internal server error"
// @Router       /api/jaco/reel [get]
func (h *Handler) GetUserReel(c *gin.Context) {
	username := c.Query("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, BadRequestError{
			Error: "username parameter is required",
		})
		return
	}

	limit := 12
	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	nextSinceID := c.Query("next_since_id")

	data, err := h.svc.GetUserReelByUsername(c.Request.Context(), username, limit, nextSinceID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			c.JSON(http.StatusNotFound, NotFoundError{
				Error: "user not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, InternalServerError{
			Error: "internal server error",
		})
		return
	}

	var resp UserReelResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError{
			Error: "failed to parse reel response",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetReelComments godoc
// @Summary      Get reel comments
// @Description  Fetch top-level comments for a Jaco reel
// @Tags         Jaco
// @Produce      json
// @Param        rid       query string true  "Reel ID (default 23194424305516544)"
// @Param        limit     query int    false "Max comments (default 20)"
// @Param        base_id   query string false "Pagination cursor"
// @Success      200 {object} CommentTreeResponse
// @Failure      400 {object} BadRequestError
// @Failure      500 {object} InternalServerError
// @Router       /api/jaco/comments [get]
func (h *Handler) GetReelComments(c *gin.Context) {
	rid := c.Query("rid")
	if rid == "" {
		c.JSON(http.StatusBadRequest, BadRequestError{
			Error: "rid parameter is required",
		})
		return
	}

	limit := 20
	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	baseID := c.Query("base_id")

	data, err := h.svc.GetReelComments(c.Request.Context(), rid, limit, baseID)
	if err != nil {
		if errors.Is(err, ErrCommentsNotFound) {
			c.JSON(http.StatusNotFound, NotFoundError{
				Error: "comments not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, InternalServerError{
			Error: "failed to fetch comments",
		})
		return
	}

	var resp CommentTreeResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError{
			Error: "invalid comment response",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
