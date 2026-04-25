package middlewares

import (
	"middleware/internal/domain/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := SplitJWTToken(c)
		if len(t) == 2 {
			authToken := t[1]
			authorized, err := security.IsAuthorized(authToken, secretKey)
			if authorized {
				userID, _, err := security.ExtractIdToken(authToken, secretKey)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, "a")
					return
				}
				c.Set("id", userID)
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, "abortd jwt "+err.Error())
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauthorized")
	}
}
