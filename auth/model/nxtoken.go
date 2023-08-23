package model

type GenerateNXTokenRequest struct {
	Count       int    `json:"count"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ExpiresAt   string `json:"expires_at"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
}

type GenerateNXTokenResponse struct {
	Token string `json:"token"`
}
