package openai

type ProviderOpenAI struct {
	Payload   Payload `json:"payload"`
	encrypted bool    `json:"encrypted"`
}

type Payload struct {
	Provider string     `json:"provider" validate:"required" enum:"openai"`
	Keys     OpenAIKeys `json:"keys" validate:"required"`
	Name     string     `json:"name" validate:"required"`
}

type OpenAIKeys struct {
	OpenAIOrg string `json:"openai_org" validate:"required"`
	OpenAIKey string `json:"openai_key"	validate:"required"`
}
