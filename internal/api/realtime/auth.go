package grpcserver

import (
	"context"
	"fmt"
	"middleware/internal/domain/security"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	authorizationMetadataKey = "authorization"
	authSchemeBearer         = "bearer"
	contextUserIDKey         = contextKey("grpc.user.id")
	contextUserRoleKey       = contextKey("grpc.user.role")
)

func UnaryAuthInterceptor(secretKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		authenticatedContext, err := authenticateContext(ctx, secretKey)
		if err != nil {
			return nil, annotateUnauthenticatedError(info.FullMethod, err)
		}

		return handler(authenticatedContext, req)
	}
}

func StreamAuthInterceptor(secretKey string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		authenticatedContext, err := authenticateContext(stream.Context(), secretKey)
		if err != nil {
			return annotateUnauthenticatedError(info.FullMethod, err)
		}

		wrappedStream := &authenticatedServerStream{
			ServerStream: stream,
			ctx:          authenticatedContext,
		}
		return handler(srv, wrappedStream)
	}
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(contextUserIDKey).(string)
	return userID, ok
}

func UserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(contextUserRoleKey).(string)
	return role, ok
}

type authenticatedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *authenticatedServerStream) Context() context.Context {
	return s.ctx
}

func authenticateContext(ctx context.Context, secretKey string) (context.Context, error) {
	token, err := bearerTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	authorized, err := security.IsAuthorized(token, secretKey)
	if err != nil || !authorized {
		if err == nil {
			err = fmt.Errorf("token nao autorizado")
		}
		return nil, err
	}

	userID, role, err := security.ExtractIdToken(token, secretKey)
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, contextUserIDKey, userID)
	ctx = context.WithValue(ctx, contextUserRoleKey, role)
	return ctx, nil
}

func bearerTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("metadata de autenticacao ausente")
	}

	values := md.Get(authorizationMetadataKey)
	if len(values) == 0 {
		return "", fmt.Errorf("cabecalho authorization ausente")
	}

	parts := strings.Fields(values[0])
	if len(parts) != 2 || !strings.EqualFold(parts[0], authSchemeBearer) {
		return "", fmt.Errorf("authorization deve usar Bearer token")
	}

	if parts[1] == "" {
		return "", fmt.Errorf("token JWT ausente")
	}

	return parts[1], nil
}

func annotateUnauthenticatedError(method string, err error) error {
	return status.Errorf(codes.Unauthenticated, "falha ao autenticar chamada %s: %v", method, err)
}
