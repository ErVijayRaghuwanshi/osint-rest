package header

import (
	"testing"
)

func testConfig() *HeaderConfig {
	t1 := true
	f := false
	return &HeaderConfig{
		Version: "1.0",
		Platforms: map[string]PlatformConfig{
			"test": {
				Enabled:  true,
				Rotation: "round_robin",
				Headers: []HeaderEntry{
					{ID: "h1", Headers: map[string]string{"User-Agent": "ua1"}, Enabled: &t1},
					{ID: "h2", Headers: map[string]string{"User-Agent": "ua2"}, Enabled: &t1},
					{ID: "h3", Headers: map[string]string{"User-Agent": "ua3"}, Enabled: &f},
				},
			},
			"disabled": {
				Enabled:  false,
				Rotation: "round_robin",
				Headers: []HeaderEntry{
					{ID: "d1", Headers: map[string]string{"User-Agent": "disabled"}},
				},
			},
		},
	}
}

func TestNewManager(t *testing.T) {
	m := NewManager(testConfig())
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
	if len(m.counters) != 2 {
		t.Fatalf("expected 2 counters, got %d", len(m.counters))
	}
}

func TestGetHeaders_RoundRobin(t *testing.T) {
	m := NewManager(testConfig())

	// h3 is disabled, so we should only cycle through h1 and h2
	ids := make([]string, 4)
	for i := 0; i < 4; i++ {
		_, id, err := m.GetHeaders("test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		ids[i] = id
	}

	// Should alternate between h1 and h2 (h3 is disabled)
	for _, id := range ids {
		if id != "h1" && id != "h2" {
			t.Errorf("expected h1 or h2, got %s", id)
		}
	}

	// Verify round-robin: ids should alternate
	if ids[0] == ids[1] {
		t.Errorf("expected different IDs for consecutive calls, got %s and %s", ids[0], ids[1])
	}
}

func TestGetHeaders_DisabledPlatform(t *testing.T) {
	m := NewManager(testConfig())

	_, _, err := m.GetHeaders("disabled")
	if err == nil {
		t.Fatal("expected error for disabled platform")
	}
}

func TestGetHeaders_UnknownPlatform(t *testing.T) {
	m := NewManager(testConfig())

	_, _, err := m.GetHeaders("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown platform")
	}
}

func TestGetHeaders_CopiesHeaders(t *testing.T) {
	m := NewManager(testConfig())

	headers, _, err := m.GetHeaders("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Mutate the returned map
	headers["User-Agent"] = "mutated"

	// Original should be unchanged
	headers2, _, _ := m.GetHeaders("test")
	if headers2["User-Agent"] == "mutated" {
		t.Error("GetHeaders should return a copy, not a reference")
	}
}

func TestReload(t *testing.T) {
	m := NewManager(testConfig())

	newCfg := &HeaderConfig{
		Version: "2.0",
		Platforms: map[string]PlatformConfig{
			"test": {
				Enabled:  true,
				Rotation: "round_robin",
				Headers: []HeaderEntry{
					{ID: "new1", Headers: map[string]string{"X": "1"}},
				},
			},
			"newplatform": {
				Enabled:  true,
				Rotation: "round_robin",
				Headers: []HeaderEntry{
					{ID: "np1", Headers: map[string]string{"Y": "2"}},
				},
			},
		},
	}

	m.Reload(newCfg)

	// Old platform should use new headers
	headers, id, err := m.GetHeaders("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "new1" {
		t.Errorf("expected new1, got %s", id)
	}
	if headers["X"] != "1" {
		t.Errorf("expected X=1, got %s", headers["X"])
	}

	// New platform should work
	_, id, err = m.GetHeaders("newplatform")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "np1" {
		t.Errorf("expected np1, got %s", id)
	}
}

func TestGetCacheTTL(t *testing.T) {
	cfg := testConfig()
	cfg.Platforms["test"] = PlatformConfig{
		Enabled:         true,
		Rotation:        "round_robin",
		Headers:         cfg.Platforms["test"].Headers,
		CacheTTLSeconds: 300,
	}

	m := NewManager(cfg)

	ttl := m.GetCacheTTL("test")
	if ttl != 300 {
		t.Errorf("expected 300, got %d", ttl)
	}

	ttl = m.GetCacheTTL("nonexistent")
	if ttl != 0 {
		t.Errorf("expected 0 for unknown platform, got %d", ttl)
	}
}

func TestAddHeader(t *testing.T) {
	m := NewManager(testConfig())

	err := m.AddHeader("test", HeaderEntry{ID: "h4", Headers: map[string]string{"Z": "4"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := m.GetConfig()
	if len(cfg.Platforms["test"].Headers) != 4 {
		t.Errorf("expected 4 headers, got %d", len(cfg.Platforms["test"].Headers))
	}

	err = m.AddHeader("nonexistent", HeaderEntry{ID: "x"})
	if err == nil {
		t.Error("expected error for unknown platform")
	}
}

func TestUpdateHeader(t *testing.T) {
	m := NewManager(testConfig())

	err := m.UpdateHeader("test", HeaderEntry{ID: "h1", Headers: map[string]string{"User-Agent": "updated"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := m.GetConfig()
	found := false
	for _, h := range cfg.Platforms["test"].Headers {
		if h.ID == "h1" {
			if h.Headers["User-Agent"] != "updated" {
				t.Errorf("expected updated, got %s", h.Headers["User-Agent"])
			}
			found = true
		}
	}
	if !found {
		t.Error("h1 not found after update")
	}

	err = m.UpdateHeader("test", HeaderEntry{ID: "nonexistent"})
	if err == nil {
		t.Error("expected error for unknown header ID")
	}
}

func TestDeleteHeader(t *testing.T) {
	m := NewManager(testConfig())

	err := m.DeleteHeader("test", "h1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfg := m.GetConfig()
	if len(cfg.Platforms["test"].Headers) != 2 {
		t.Errorf("expected 2 headers after delete, got %d", len(cfg.Platforms["test"].Headers))
	}

	err = m.DeleteHeader("test", "nonexistent")
	if err == nil {
		t.Error("expected error for unknown header ID")
	}
}

func TestToggleHeader(t *testing.T) {
	m := NewManager(testConfig())

	// h1 starts enabled
	enabled, err := m.ToggleHeader("test", "h1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enabled {
		t.Error("expected h1 to be disabled after toggle")
	}

	// Toggle back
	enabled, err = m.ToggleHeader("test", "h1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !enabled {
		t.Error("expected h1 to be enabled after second toggle")
	}

	_, err = m.ToggleHeader("test", "nonexistent")
	if err == nil {
		t.Error("expected error for unknown header ID")
	}
}
