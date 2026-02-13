package header

import (
	"errors"
	"fmt"
	"sync"
)

type Manager struct {
	cfg      *HeaderConfig
	counters map[string]*counter
	mu       sync.RWMutex
}

func NewManager(cfg *HeaderConfig) *Manager {
	counters := make(map[string]*counter)
	for platform := range cfg.Platforms {
		counters[platform] = &counter{}
	}

	return &Manager{
		cfg:      cfg,
		counters: counters,
	}
}

func (m *Manager) GetHeaders(platform string) (map[string]string, string, error) {
	entry, err := m.Next(platform)
	if err != nil {
		return nil, "", err
	}

	// Copy to prevent mutation
	headers := make(map[string]string, len(entry.Headers))
	for k, v := range entry.Headers {
		headers[k] = v
	}

	return headers, entry.ID, nil
}

// Reload atomically swaps the header configuration and resets counters
// for any new platforms. Existing counters are preserved.
func (m *Manager) Reload(cfg *HeaderConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cfg = cfg

	// Add counters for new platforms, keep existing ones
	for platform := range cfg.Platforms {
		if _, ok := m.counters[platform]; !ok {
			m.counters[platform] = &counter{}
		}
	}
}

// GetCacheTTL returns the configured cache TTL in seconds for a platform.
// Returns 0 if the platform is not found or has no TTL configured.
func (m *Manager) GetCacheTTL(platform string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, ok := m.cfg.Platforms[platform]
	if !ok {
		return 0
	}
	return p.CacheTTLSeconds
}

// GetConfig returns a copy of the current header configuration.
func (m *Manager) GetConfig() *HeaderConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cfg
}

// AddHeader appends a new header entry to a platform.
func (m *Manager) AddHeader(platform string, entry HeaderEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.cfg.Platforms[platform]
	if !ok {
		return fmt.Errorf("platform %q not found", platform)
	}

	// Default to enabled if not specified
	if entry.Enabled == nil {
		t := true
		entry.Enabled = &t
	}

	p.Headers = append(p.Headers, entry)
	m.cfg.Platforms[platform] = p
	return nil
}

// UpdateHeader replaces the headers for an existing entry by ID.
func (m *Manager) UpdateHeader(platform string, entry HeaderEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.cfg.Platforms[platform]
	if !ok {
		return fmt.Errorf("platform %q not found", platform)
	}

	for i, h := range p.Headers {
		if h.ID == entry.ID {
			p.Headers[i].Headers = entry.Headers
			if entry.Enabled != nil {
				p.Headers[i].Enabled = entry.Enabled
			}
			m.cfg.Platforms[platform] = p
			return nil
		}
	}

	return fmt.Errorf("header %q not found in platform %q", entry.ID, platform)
}

// DeleteHeader removes a header entry by ID from a platform.
func (m *Manager) DeleteHeader(platform, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.cfg.Platforms[platform]
	if !ok {
		return fmt.Errorf("platform %q not found", platform)
	}

	for i, h := range p.Headers {
		if h.ID == id {
			p.Headers = append(p.Headers[:i], p.Headers[i+1:]...)
			m.cfg.Platforms[platform] = p
			return nil
		}
	}

	return fmt.Errorf("header %q not found in platform %q", id, platform)
}

// ToggleHeader flips the enabled state of a header entry and returns the new state.
func (m *Manager) ToggleHeader(platform, id string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.cfg.Platforms[platform]
	if !ok {
		return false, fmt.Errorf("platform %q not found", platform)
	}

	for i, h := range p.Headers {
		if h.ID == id {
			current := h.Enabled == nil || *h.Enabled
			newVal := !current
			p.Headers[i].Enabled = &newVal
			m.cfg.Platforms[platform] = p
			return newVal, nil
		}
	}

	return false, fmt.Errorf("header %q not found in platform %q", id, platform)
}

// isEntryEnabled returns true if the entry is enabled (nil defaults to true).
func isEntryEnabled(e *HeaderEntry) bool {
	return e.Enabled == nil || *e.Enabled
}

func (m *Manager) Next(platform string) (*HeaderEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	p, ok := m.cfg.Platforms[platform]
	if !ok || !p.Enabled {
		return nil, errors.New("platform not enabled or not found")
	}

	// Filter to only enabled headers
	var enabled []HeaderEntry
	for _, h := range p.Headers {
		if h.Enabled == nil || *h.Enabled {
			enabled = append(enabled, h)
		}
	}

	if len(enabled) == 0 {
		return nil, errors.New("no enabled headers configured")
	}

	c, ok := m.counters[platform]
	if !ok {
		c = &counter{}
		m.counters[platform] = c
	}

	switch p.Rotation {
	case "round_robin", "":
		idx := c.Next(len(enabled))
		entry := enabled[idx]
		return &entry, nil
	default:
		return nil, errors.New("unsupported rotation strategy")
	}
}
