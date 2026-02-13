package platform

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Registry holds all registered platform implementations and provides
// auto-mounting of routes.
type Registry struct {
	mu        sync.RWMutex
	platforms []Platform
	log       zerolog.Logger
}

// NewRegistry creates a new empty platform registry.
func NewRegistry(log zerolog.Logger) *Registry {
	return &Registry{
		platforms: make([]Platform, 0),
		log:       log,
	}
}

// Register adds a platform to the registry.
func (r *Registry) Register(p Platform) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.platforms = append(r.platforms, p)
	r.log.Info().Str("platform", p.Name()).Msg("Platform registered")
}

// All returns a snapshot of all registered platforms.
func (r *Registry) All() []Platform {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Platform, len(r.platforms))
	copy(out, r.platforms)
	return out
}

// MountAll registers routes for every platform under /api/<name>.
func (r *Registry) MountAll(engine *gin.Engine) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.platforms {
		group := engine.Group("/api/" + p.Name())
		p.RegisterRoutes(group)
		r.log.Info().Str("platform", p.Name()).Msg("Routes mounted")
	}
}
