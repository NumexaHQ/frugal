package model

import "time"

type AllRequestsTableResponse struct {
	ID             string    `json:"id"`
	StatusCode     int       `json:"status_code"`
	ProjectID      int32     `json:"project_id"`
	InitiatedAt    time.Time `json:"initiated_at"`
	Model          string    `json:"model"`
	Prompt         string    `json:"prompt"`
	Cost           float64   `json:"cost"`
	Provider       string    `json:"provider"`
	Latency        int64     `json:"latency"`
	IsCached       bool      `json:"is_cached"`
	IsCacheHit     bool      `json:"is_cache_hit"`
	CustomMetaData string    `json:"custom_metadata"`
}
