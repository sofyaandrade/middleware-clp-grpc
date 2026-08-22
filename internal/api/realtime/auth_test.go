package grpcserver

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestAuthenticateContextStoresUserClaims(t *testing.T) {
	const secretKey = "grpc-test-secret"

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		authorizationMetadataKey,
		"Bearer "+signedTestToken(t, secretKey, 7, "CONSUMIDOR"),
	))

	authenticatedContext, err := authenticateContext(ctx, secretKey)
	if err != nil {
		t.Fatalf("authenticateContext() error = %v", err)
	}

	userID, ok := UserIDFromContext(authenticatedContext)
	if !ok || userID != "7" {
		t.Fatalf("userID = %q, %v, want 7, true", userID, ok)
	}

	role, ok := UserRoleFromContext(authenticatedContext)
	if !ok || role != "CONSUMIDOR" {
		t.Fatalf("role = %q, %v, want CONSUMIDOR, true", role, ok)
	}
}

func TestUnaryAuthInterceptorRejectsMissingToken(t *testing.T) {
	interceptor := UnaryAuthInterceptor("grpc-test-secret")

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/middleware.realtime.v1.RealtimeTagCatalogService/GetTagsSnapshot",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("status.Code(err) = %v, want %v", status.Code(err), codes.Unauthenticated)
	}
}

func TestStreamAuthInterceptorInjectsAuthenticatedContext(t *testing.T) {
	const secretKey = "grpc-test-secret"

	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs(
		authorizationMetadataKey,
		"Bearer "+signedTestToken(t, secretKey, 9, "CONSUMIDOR"),
	))
	stream := &stubServerStream{ctx: ctx}

	interceptor := StreamAuthInterceptor(secretKey)
	err := interceptor(nil, stream, &grpc.StreamServerInfo{
		FullMethod:     "/middleware.realtime.v1.RealtimeTagService/StreamTagValues",
		IsServerStream: true,
		IsClientStream: false,
	}, func(srv interface{}, serverStream grpc.ServerStream) error {
		userID, ok := UserIDFromContext(serverStream.Context())
		if !ok || userID != "9" {
			t.Fatalf("userID = %q, %v, want 9, true", userID, ok)
		}

		role, ok := UserRoleFromContext(serverStream.Context())
		if !ok || role != "CONSUMIDOR" {
			t.Fatalf("role = %q, %v, want CONSUMIDOR, true", role, ok)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("interceptor() error = %v", err)
	}
}

func signedTestToken(t *testing.T, secretKey string, userID float64, role string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"Id":   userID,
		"Role": role,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		t.Fatalf("SignedString() error = %v", err)
	}
	return signedToken
}

type stubServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *stubServerStream) Context() context.Context {
	return s.ctx
}
