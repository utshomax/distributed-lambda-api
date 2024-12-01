package config

import "encoding/json"

type RequestConfig struct {
	URL              string            `json:"url"`
	RequestCount     int               `json:"requestCount,omitempty"`
	BatchSize        int               `json:"batch,omitempty"`
	Method           string            `json:"method,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	DNSServers       []string          `json:"dns,omitempty"`
	DisableKeepAlive bool              `json:"disableKeepAlive,omitempty"`
	DisableDNSCache  bool              `json:"disableDNSCache,omitempty"`
}

func Parse(data []byte) (RequestConfig, error) {
	var cfg RequestConfig
	err := json.Unmarshal(data, &cfg)

	// Set defaults
	if cfg.RequestCount == 0 {
		cfg.RequestCount = 100
	}
	if cfg.BatchSize == 0 {
		cfg.BatchSize = cfg.RequestCount
	}
	if cfg.Method == "" {
		cfg.Method = "GET"
	}
	// Keep-alive and DNS cache are enabled by default (false)

	return cfg, err
}
