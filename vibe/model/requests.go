package model

type AllRequestsTableResponse struct {
	ID             string  `json:"id"`
	StatusCode     int     `json:"status_code"`
	ProjectID      int32   `json:"project_id"`
	InitiatedAt    string  `json:"initiated_at"`
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt"`
	Cost           float64 `json:"cost"`
	Provider       string  `json:"provider"`
	IsCached       bool    `json:"is_cached"`
	IsCacheHit     bool    `json:"is_cache_hit"`
	CustomMetaData string  `json:"custom_meta_data"`
}
