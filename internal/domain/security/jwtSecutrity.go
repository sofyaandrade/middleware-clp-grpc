package security

import (
	"fmt"
	"time"

	"middleware/internal/domain/conversion"
	"middleware/internal/domain/models"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func expirationTimeFromHours(hours float64) time.Time {
	if hours <= 0 {
		return time.Now().Add(time.Minute)
	}
	return time.Now().Add(time.Duration(hours * float64(time.Hour)))
}

func CreateAcessToken(usuario *models.User) (accessToken string, err error) {
	expiry := float64(24)
	secretAccesKey := viper.GetString("ACCESS_TOKEN")
	emissor := "middleware"

	claims := &models.JwtCustomClaims{
		ID:   usuario.ID,
		Role: usuario.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeFromHours(expiry).Unix(),
			Issuer:    emissor,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secretAccesKey))
	if err != nil {
		return "", err
	}
	return t, err
}

func CreateRefreshToken(usuario *models.User) (refreshToken string, err error) {
	secretRefreshKey := viper.GetString("REFRESH_TOKEN")
	expiry := float64(128)

	claimsRefresh := &models.JwtCustomRefreshClaims{
		ID: usuario.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTimeFromHours(expiry).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secretRefreshKey))
	if err != nil {
		return "", err
	}
	return rt, err
}

func IsAuthorized(requestToken, secretKey string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIdToken(requestToken, secretRequestKey string) (string, string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(secretRequestKey), nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("token inválido")
	}

	id, ok := claims["Id"].(float64)
	if !ok {
		return "", "", fmt.Errorf("token inválido")
	}

	role, _ := claims["Role"].(string)
	return conversion.Float64ToString(id), role, nil
}
