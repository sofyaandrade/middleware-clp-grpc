package models

type Login struct {
	Email    string
	Password string
}

type ResponseLogin struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}
