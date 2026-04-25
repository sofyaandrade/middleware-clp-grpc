package models

type RefreshTokenRequest struct {
	RefreshToken string `form:"RefreshToken" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
