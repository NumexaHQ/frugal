package model

// Dashboard responses here

// DashboardRequestStats
type GetUserRequestsStatsByProjectID map[string]UserUsageStats

type UserUsageStats struct {
	TotalRequest int     `json:"total_request"`
	TotalCost    float64 `json:"total_cost"`
}

type GetUserRequestsStatsByProjectIDResponse struct {
	Email        string  `json:"email"`
	TotalRequest int     `json:"total_request"`
	TotalCost    float64 `json:"cost"`
}
