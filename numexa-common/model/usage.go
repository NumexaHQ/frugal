package model

type UsageLimitSetting struct {
	RequestLimit []PlanRequestLimit `json:"request_limit"`
}

type PlanRequestLimit struct {
	PlanID string `json:"plan_id"`
	Name   string `json:"name"`
	Limit  int    `json:"limit"`
}
