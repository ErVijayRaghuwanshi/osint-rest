package header

type HeaderConfig struct {
	Version   string                    `json:"version"`
	Platforms map[string]PlatformConfig `json:"platforms"`
}

type PlatformConfig struct {
	Enabled         bool          `json:"enabled"`
	Rotation        string        `json:"rotation"`
	Headers         []HeaderEntry `json:"headers"`
	CacheTTLSeconds int           `json:"cache_ttl_seconds,omitempty"`
}

type HeaderEntry struct {
	ID      string            `json:"id"`
	Headers map[string]string `json:"headers"`
	Enabled *bool             `json:"enabled,omitempty"`
}
