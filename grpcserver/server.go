package grpcserver

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/richard-ramos/komainu/certs"
	"github.com/richard-ramos/komainu/proto"
	"github.com/richard-ramos/komainu/services/service1"
	"github.com/richard-ramos/komainu/services/session"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func NewGRPCServer(lc fx.Lifecycle, service1 *service1.Service1, jwtManager *session.JWTManager, sessionService *session.SessionServer) *grpc.Server {
	logger := zap.NewNop()

	grpcLogger := logger.Named("grpc")

	panicHandler := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic recover: %v", p)
	}

	recoverOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(panicHandler),
	}

	zapOpts := []grpc_zap.Option{}

	setGRPCLogger(grpcLogger)

	interceptor := session.NewAuthInterceptor(jwtManager)

	grpcOpts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(recoverOpts...),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(grpcLogger, zapOpts...),
			grpc_auth.UnaryServerInterceptor(interceptor.Interceptor()),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(recoverOpts...),
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(grpcLogger, zapOpts...),
			grpc_auth.StreamServerInterceptor(interceptor.Interceptor()),
		),
	}

	// Setting up SSL
	config := &tls.Config{
		Certificates: []tls.Certificate{*certs.GlobalCertificate()},
		ClientAuth:   tls.NoClientCert,
	}

	creds := credentials.NewTLS(config)
	grpcOpts = append(grpcOpts, grpc.Creds(creds))

	grpcServer := grpc.NewServer(grpcOpts...)

	proto.RegisterChatServiceServer(grpcServer, service1)
	proto.RegisterSessionServiceServer(grpcServer, sessionService)

	addr := ":8080"

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			fmt.Println("Starting GRPC server at", addr)
			go grpcServer.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			grpcServer.Stop()
			return nil
		},
	})

	return grpcServer
}
