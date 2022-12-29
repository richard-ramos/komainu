package grpcserver

import (
	"sync"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
)

var (
	grpcLoggerSet bool
	muGRPCLogger  sync.Mutex
)

func setGRPCLogger(l *zap.Logger) {
	muGRPCLogger.Lock()
	defer muGRPCLogger.Unlock()

	if !grpcLoggerSet {
		grpc_zap.ReplaceGrpcLoggerV2(l)
		grpcLoggerSet = true
	}
}
