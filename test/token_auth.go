package main

import "context"

// TokenAuth is a JWT based credentials provider for gRPC
type TokenAuth string

// GetRequestMetadata implements credentials.PerRPCCredentials
func (t TokenAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	return map[string]string{"authorization": "Bearer " + string(t)}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (TokenAuth) RequireTransportSecurity() bool {
	return false
}
