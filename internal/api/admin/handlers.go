package admin

import (
	"encoding/json"
	"net/http"
	"os"

	"osint-scraper/internal/header"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	hm             *header.Manager
	log            zerolog.Logger
	headerFilePath string
}

func NewHandler(hm *header.Manager, log zerolog.Logger, headerFilePath string) *Handler {
	return &Handler{hm: hm, log: log, headerFilePath: headerFilePath}
}

// ListPlatforms godoc
// @Summary      List all platforms and their header sets
// @Description  Returns the full header configuration for all platforms
// @Tags         Admin
// @Security     ApiKeyAuth
// @Produce      json
// @Success      200  {object}  header.HeaderConfig
// @Router       /admin/headers [get]
func (h *Handler) ListPlatforms(c *gin.Context) {
	cfg := h.hm.GetConfig()
	c.JSON(http.StatusOK, cfg)
}

// GetPlatformHeaders godoc
// @Summary      Get headers for a specific platform
// @Description  Returns the header configuration for a single platform
// @Tags         Admin
// @Security     ApiKeyAuth
// @Produce      json
// @Param        platform  path  string  true  "Platform name (e.g. instagram, snapchat, x, jaco)"
// @Success      200  {object}  header.PlatformConfig
// @Failure      404  {object}  map[string]string
// @Router       /admin/headers/{platform} [get]
func (h *Handler) GetPlatformHeaders(c *gin.Context) {
	platform := c.Param("platform")
	cfg := h.hm.GetConfig()

	p, ok := cfg.Platforms[platform]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "platform not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// AddHeaderSet godoc
// @Summary      Add a new header set to a platform
// @Description  Appends a new header entry to the platform's header list
// @Tags         Admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        platform  path  string              true  "Platform name"
// @Param        body      body  header.HeaderEntry   true  "Header entry to add"
// @Success      201  {object}  header.HeaderEntry
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/headers/{platform} [post]
func (h *Handler) AddHeaderSet(c *gin.Context) {
	platform := c.Param("platform")

	var entry header.HeaderEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if entry.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := h.hm.AddHeader(platform, entry); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.persistConfig()
	h.log.Info().Str("platform", platform).Str("id", entry.ID).Msg("Header set added")
	c.JSON(http.StatusCreated, entry)
}

// UpdateHeaderSet godoc
// @Summary      Update an existing header set
// @Description  Replaces the headers for a specific header entry by ID
// @Tags         Admin
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        platform  path  string              true  "Platform name"
// @Param        id        path  string              true  "Header set ID"
// @Param        body      body  header.HeaderEntry   true  "Updated header entry"
// @Success      200  {object}  header.HeaderEntry
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/headers/{platform}/{id} [put]
func (h *Handler) UpdateHeaderSet(c *gin.Context) {
	platform := c.Param("platform")
	id := c.Param("id")

	var entry header.HeaderEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.ID = id

	if err := h.hm.UpdateHeader(platform, entry); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.persistConfig()
	h.log.Info().Str("platform", platform).Str("id", id).Msg("Header set updated")
	c.JSON(http.StatusOK, entry)
}

// DeleteHeaderSet godoc
// @Summary      Delete a header set
// @Description  Removes a header entry by ID from a platform
// @Tags         Admin
// @Security     ApiKeyAuth
// @Produce      json
// @Param        platform  path  string  true  "Platform name"
// @Param        id        path  string  true  "Header set ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /admin/headers/{platform}/{id} [delete]
func (h *Handler) DeleteHeaderSet(c *gin.Context) {
	platform := c.Param("platform")
	id := c.Param("id")

	if err := h.hm.DeleteHeader(platform, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.persistConfig()
	h.log.Info().Str("platform", platform).Str("id", id).Msg("Header set deleted")
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ToggleHeaderSet godoc
// @Summary      Toggle a header set enabled/disabled
// @Description  Enables or disables a specific header entry
// @Tags         Admin
// @Security     ApiKeyAuth
// @Produce      json
// @Param        platform  path  string  true  "Platform name"
// @Param        id        path  string  true  "Header set ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /admin/headers/{platform}/{id}/toggle [post]
func (h *Handler) ToggleHeaderSet(c *gin.Context) {
	platform := c.Param("platform")
	id := c.Param("id")

	enabled, err := h.hm.ToggleHeader(platform, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.persistConfig()
	h.log.Info().Str("platform", platform).Str("id", id).Bool("enabled", enabled).Msg("Header set toggled")
	c.JSON(http.StatusOK, gin.H{"id": id, "enabled": enabled})
}

// persistConfig writes the current config back to the header file.
func (h *Handler) persistConfig() {
	cfg := h.hm.GetConfig()
	data, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to marshal header config for persistence")
		return
	}

	if err := os.WriteFile(h.headerFilePath, data, 0644); err != nil {
		h.log.Error().Err(err).Msg("Failed to persist header config")
	}
}
