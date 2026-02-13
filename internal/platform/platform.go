package platform

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Platform is the interface that all social media platform modules must implement.
// This enables expandable platform integration — new platforms just implement
// this interface and register themselves.
type Platform interface {
	// Name returns the platform identifier (e.g. "instagram", "telegram").
	Name() string

	// RegisterRoutes mounts the platform's API endpoints onto the given route group.
	RegisterRoutes(rg *gin.RouterGroup)

	// Ping checks if the platform's target website is reachable.
	Ping(ctx context.Context) bool
}
