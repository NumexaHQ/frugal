package model

import "time"

type PromptDirectory struct {
	RequestID      string    `json:"id"`
	InitiatedAt    time.Time `json:"initiated_at"`
	IsCacheHit     bool      `json:"is_cache_hit"`
	IsCached       bool      `json:"is_cached"`
	Latency        int64     `json:"latency"`
	Model          string    `json:"model"`
	ProjectID      int32     `json:"project_id"`
	UserID         int32     `json:"user_id"`
	Prompt         string    `json:"prompt"`
	Provider       string    `json:"provider"`
	StatusCode     int64     `json:"status_code"`
	Cost           float64   `json:"cost"`
	CustomMetadata string    `json:"custom_metadata"`
	Score          int32     `json:"score"`
	Comment        string    `json:"comment"`
}
