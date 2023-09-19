package model

type ProviderKeys struct {
	Name      string            `json:"name"`
	Provider  string            `json:"provider" validate:"required" enum:"openai"`
	Keys      map[string]string `json:"keys"`
	ProjectId int32             `json:"project_id"`
}
