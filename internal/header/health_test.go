package header

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func testLogger() zerolog.Logger {
	return zerolog.New(os.Stderr).With().Timestamp().Logger()
}

func TestHealthTracker_ReportSuccess(t *testing.T) {
	m := NewManager(testConfig())
	ht := NewHealthTracker(m, 3, testLogger())

	ht.ReportSuccess("test", "h1")
	ht.ReportSuccess("test", "h1")

	stats := ht.GetStats()
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat entry, got %d", len(stats))
	}
	if stats[0].Successes != 2 {
		t.Errorf("expected 2 successes, got %d", stats[0].Successes)
	}
	if stats[0].ConsecutiveFails != 0 {
		t.Errorf("expected 0 consecutive fails, got %d", stats[0].ConsecutiveFails)
	}
}

func TestHealthTracker_ReportFailure(t *testing.T) {
	m := NewManager(testConfig())
	ht := NewHealthTracker(m, 3, testLogger())

	ht.ReportFailure("test", "h1")
	ht.ReportFailure("test", "h1")

	stats := ht.GetStats()
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat entry, got %d", len(stats))
	}
	if stats[0].Failures != 2 {
		t.Errorf("expected 2 failures, got %d", stats[0].Failures)
	}
	if stats[0].ConsecutiveFails != 2 {
		t.Errorf("expected 2 consecutive fails, got %d", stats[0].ConsecutiveFails)
	}
	if stats[0].AutoDisabled {
		t.Error("should not be auto-disabled yet (threshold=3)")
	}
}

func TestHealthTracker_AutoDisable(t *testing.T) {
	m := NewManager(testConfig())
	ht := NewHealthTracker(m, 3, testLogger())

	// Trigger 3 consecutive failures → auto-disable
	ht.ReportFailure("test", "h1")
	ht.ReportFailure("test", "h1")
	ht.ReportFailure("test", "h1")

	stats := ht.GetStats()
	found := false
	for _, s := range stats {
		if s.HeaderID == "h1" {
			found = true
			if !s.AutoDisabled {
				t.Error("expected h1 to be auto-disabled after 3 consecutive failures")
			}
		}
	}
	if !found {
		t.Error("h1 stats not found")
	}

	// Verify the header was actually disabled in the manager
	cfg := m.GetConfig()
	for _, h := range cfg.Platforms["test"].Headers {
		if h.ID == "h1" {
			if h.Enabled == nil || *h.Enabled {
				t.Error("expected h1 to be disabled in manager config")
			}
		}
	}
}

func TestHealthTracker_SuccessResetsConsecutiveFails(t *testing.T) {
	m := NewManager(testConfig())
	ht := NewHealthTracker(m, 5, testLogger())

	ht.ReportFailure("test", "h1")
	ht.ReportFailure("test", "h1")
	ht.ReportSuccess("test", "h1")

	stats := ht.GetStats()
	for _, s := range stats {
		if s.HeaderID == "h1" {
			if s.ConsecutiveFails != 0 {
				t.Errorf("expected 0 consecutive fails after success, got %d", s.ConsecutiveFails)
			}
			if s.Failures != 2 {
				t.Errorf("expected 2 total failures, got %d", s.Failures)
			}
			if s.Successes != 1 {
				t.Errorf("expected 1 success, got %d", s.Successes)
			}
		}
	}
}
