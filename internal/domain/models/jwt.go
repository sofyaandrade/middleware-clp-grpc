package models

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	ID   uint   `json:"Id"`
	Role string `json:"Role"`
	jwt.StandardClaims
}

type JwtCustomRefreshClaims struct {
	ID uint `json:"Id"`
	jwt.StandardClaims
}
