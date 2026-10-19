package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"middleware/internal/api/middlewares"
	"middleware/internal/infrastructure/database/migrations"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TestOperatorHTTPPermissions(t *testing.T) {
	policyModel, err := model.NewModelFromString(rbacModel)
	if err != nil {
		t.Fatal(err)
	}
	enforcer, err := casbin.NewEnforcer(policyModel)
	if err != nil {
		t.Fatal(err)
	}
	// Existing installations must receive the new permissions on startup too.
	if _, err := enforcer.AddPolicy("OPERADOR", "/clps/*", "(POST)|(PATCH)|(DELETE)"); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/clps", "/clps/*", "/tags", "/tags/*"} {
		if _, err := enforcer.AddPolicy("OPERADOR", path, "(GET)|(POST)|(PATCH)|(DELETE)"); err != nil {
			t.Fatal(err)
		}
	}
	migrations.InicializaPermissoesAcesso(enforcer)
	migrations.InicializaPermissoesAcesso(enforcer)

	const secret = "permissions-test-secret"
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middlewares.AccessPermissions(secret, enforcer), middlewares.JwtAuthMiddleware(secret))
	router.NoRoute(func(c *gin.Context) {
		// Password changes must use the authenticated user's ID.
		c.String(http.StatusOK, c.GetString("id"))
	})

	tests := []struct {
		role   string
		method string
		path   string
		status int
	}{
		{"OPERADOR", "GET", "/clps/", 200},
		{"OPERADOR", "GET", "/clps", 200},
		{"OPERADOR", "GET", "/clps/7/", 200},
		{"OPERADOR", "GET", "/clps/status/", 200},
		{"OPERADOR", "GET", "/clps/type/1/", 200},
		{"OPERADOR", "GET", "/tags/", 200},
		{"OPERADOR", "GET", "/tags/real-time/", 200},
		{"OPERADOR", "GET", "/type-clps/", 200},
		{"OPERADOR", "GET", "/type-tags/", 200},
		{"OPERADOR", "GET", "/type-operations/", 200},
		{"OPERADOR", "GET", "/swaps/", 200},
		{"OPERADOR", "PATCH", "/clps/7/", 200},
		{"OPERADOR", "PATCH", "/tags/7/", 200},
		{"OPERADOR", "GET", "/tags/7/", 200},
		{"OPERADOR", "POST", "/clps/", 403},
		{"OPERADOR", "POST", "/clps", 403},
		{"OPERADOR", "DELETE", "/clps/7/", 403},
		{"OPERADOR", "POST", "/tags/", 403},
		{"OPERADOR", "POST", "/tags", 403},
		{"OPERADOR", "DELETE", "/tags/7/", 403},
		{"OPERADOR", "POST", "/type-clps/", 403},
		{"OPERADOR", "DELETE", "/type-tags/", 403},
		{"OPERADOR", "PATCH", "/users/password/", 200},
		{"OPERADOR", "GET", "/users/", 403},
		{"OPERADOR", "GET", "/users/8/", 403},
		{"OPERADOR", "POST", "/users/", 403},
		{"OPERADOR", "PATCH", "/users/8/", 403},
		{"OPERADOR", "DELETE", "/users/8/", 403},
		{"OPERADOR", "GET", "/user-profiles/", 403},
		{"OPERADOR", "PATCH", "/users/password/extra/", 403},
		{"ADMINISTRADOR", "POST", "/users/", 200},
		{"ADMINISTRADOR", "POST", "/clps/", 200},
		{"ADMINISTRADOR", "DELETE", "/clps/7/", 200},
		{"ADMINISTRADOR", "POST", "/tags/", 200},
		{"ADMINISTRADOR", "DELETE", "/tags/7/", 200},
		{"ADMINISTRADOR", "GET", "/user-profiles/", 200},
		{"ADMINISTRADOR", "PATCH", "/users/8/", 200},
		{"CONSUMIDOR", "GET", "/clps/", 403},
	}
	for _, tt := range tests {
		t.Run(tt.role+" "+tt.method+" "+tt.path, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"Id": 7, "Role": tt.role, "exp": time.Now().Add(time.Hour).Unix(),
			})
			signed, err := token.SignedString([]byte(secret))
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(tt.method, tt.path, nil)
			request.Header.Set("Authorization", "Bearer "+signed)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.status, response.Body.String())
			}
			if tt.status == http.StatusOK && response.Body.String() != "7" {
				t.Fatalf("authenticated user = %q, want 7", response.Body.String())
			}
		})
	}
}
