package middlewares

import (
	"middleware/internal/domain/security"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

func AccessPermissions(secretKey string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isPublicRoute(c.Request.URL.Path) {
			c.Next()
			return
		}

		t := SplitJWTToken(c)
		if len(t) != 2 || !strings.EqualFold(t[0], "Bearer") {
			c.Next()
			return
		}

		authToken := t[1]
		userID, usuarioPerfil, err := security.ExtractIdToken(authToken, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		ok, err := enforcer.Enforce(strings.ToUpper(usuarioPerfil), c.Request.URL.Path, c.Request.Method)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal Error")
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		c.Set("id", userID)
		c.Next()
	}
}

func SplitJWTToken(c *gin.Context) []string {
	authHeader := c.Request.Header.Get("Authorization")
	return strings.Fields(authHeader)
}

func isPublicRoute(path string) bool {
	normalizedPath := strings.TrimSuffix(path, "/")
	switch normalizedPath {
	case "/login", "/refresh":
		return true
	default:
		return false
	}
}
