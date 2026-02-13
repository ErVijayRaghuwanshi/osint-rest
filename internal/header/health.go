package header

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// HealthTracker tracks success/failure counts per header entry and
// auto-disables entries that exceed the failure threshold.
type HealthTracker struct {
	mu        sync.Mutex
	stats     map[string]*EntryStats // key: "platform:headerID"
	threshold int
	manager   *Manager
	log       zerolog.Logger
}

// EntryStats holds runtime health data for a single header entry.
type EntryStats struct {
	Platform          string    `json:"platform"`
	HeaderID          string    `json:"header_id"`
	Successes         int       `json:"successes"`
	Failures          int       `json:"failures"`
	ConsecutiveFails  int       `json:"consecutive_fails"`
	LastUsed          time.Time `json:"last_used"`
	AutoDisabled      bool      `json:"auto_disabled"`
}

// NewHealthTracker creates a tracker that auto-disables headers after
// consecutiveFailThreshold consecutive failures.
func NewHealthTracker(m *Manager, threshold int, log zerolog.Logger) *HealthTracker {
	return &HealthTracker{
		stats:     make(map[string]*EntryStats),
		threshold: threshold,
		manager:   m,
		log:       log,
	}
}

func statsKey(platform, headerID string) string {
	return platform + ":" + headerID
}

// ReportSuccess records a successful request for the given header.
func (ht *HealthTracker) ReportSuccess(platform, headerID string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	key := statsKey(platform, headerID)
	s := ht.getOrCreate(key, platform, headerID)
	s.Successes++
	s.ConsecutiveFails = 0
	s.LastUsed = time.Now()
}

// ReportFailure records a failed request. If consecutive failures exceed
// the threshold, the header is auto-disabled via the Manager.
func (ht *HealthTracker) ReportFailure(platform, headerID string) {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	key := statsKey(platform, headerID)
	s := ht.getOrCreate(key, platform, headerID)
	s.Failures++
	s.ConsecutiveFails++
	s.LastUsed = time.Now()

	if s.ConsecutiveFails >= ht.threshold && !s.AutoDisabled {
		s.AutoDisabled = true
		ht.log.Warn().
			Str("platform", platform).
			Str("header_id", headerID).
			Int("consecutive_fails", s.ConsecutiveFails).
			Msg("Auto-disabling header due to consecutive failures")

		disabled := false
		ht.manager.UpdateHeader(platform, HeaderEntry{
			ID:      headerID,
			Enabled: &disabled,
		})
	}
}

// GetStats returns a snapshot of all tracked header stats.
func (ht *HealthTracker) GetStats() []EntryStats {
	ht.mu.Lock()
	defer ht.mu.Unlock()

	result := make([]EntryStats, 0, len(ht.stats))
	for _, s := range ht.stats {
		copy := *s
		result = append(result, copy)
	}
	return result
}

func (ht *HealthTracker) getOrCreate(key, platform, headerID string) *EntryStats {
	s, ok := ht.stats[key]
	if !ok {
		s = &EntryStats{
			Platform: platform,
			HeaderID: headerID,
		}
		ht.stats[key] = s
	}
	return s
}
