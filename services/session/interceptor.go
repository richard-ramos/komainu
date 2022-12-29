package session

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager *JWTManager
}

func NewAuthInterceptor(jwtManager *JWTManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (a *AuthInterceptor) validateToken(tokens []string) (*UserClaims, error) {
	if len(tokens) < 1 || !strings.Contains(tokens[0], "Bearer ") {
		return nil, status.Errorf(codes.InvalidArgument, "missing Bearer token")
	}
	token := strings.Split(tokens[0], "Bearer ")[1]

	claims, err := a.jwtManager.Verify(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return claims, nil
}

func (a *AuthInterceptor) Interceptor() func(ctx context.Context) (context.Context, error) {
	return func(ctx context.Context) (context.Context, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		jwtToken, ok := md["authorization"]
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		claims, err := a.validateToken(jwtToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "valid token required.")
		}

		ctx = context.WithValue(ctx, "userID", claims.UserID)

		return ctx, nil
	}
}
