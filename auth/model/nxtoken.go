package model

import "time"

type GenerateNXTokenRequest struct {
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	UserID           int32                  `json:"user_id"`
	ProjectID        int32                  `json:"project_id"`
	NxaProviderKeyID int32                  `json:"nxa_provider_key_id"`
	Property         NXTokenPropertyRequest `json:"property"`
}

type NXTokenPropertyRequest struct {
	RateLimit        int32     `json:"rate_limit"`
	RateLimitPeriod  string    `json:"rate_limit_period"`
	EnforceCaching   bool      `json:"enforce_caching"`
	OverallCostLimit int32     `json:"overall_cost_limit"`
	AlertOnThreshold int32     `json:"alert_on_threshold"`
	ExpiresAt        time.Time `json:"expires_at"`
}

type GenerateNXTokenResponse struct {
	Token string `json:"token"`
}

/*
{
	"name": "prod-key",
	"description": "production key",
	"user_id": "1",
	"project_id": "1",
	"property": {
		"rate_limit": "1000",
		"rate_limit_period": "day",
		"enforce_caching": "true",
		"overall_cost_limit": "1000",
		"alert_on_threshold": "1000",
		"expires_at": "2021-01-01 00:00:00"
	}

}

*/
