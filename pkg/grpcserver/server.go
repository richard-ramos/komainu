package grpcserver

import (
	"context"
	"crypto/tls"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/richard-ramos/komainu/pkg/certs"
	"github.com/richard-ramos/komainu/pkg/config"
	"github.com/richard-ramos/komainu/pkg/proto"
	"github.com/richard-ramos/komainu/pkg/services/accounts"
	"github.com/richard-ramos/komainu/pkg/services/service1"
	"github.com/richard-ramos/komainu/pkg/services/session"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type GRPCServerParams struct {
	fx.In

	Cfg        *config.Config
	JWTManager *session.JWTManager

	Service1       *service1.Service1
	Accounts       *accounts.AccountsServer
	SessionService *session.SessionServer
}

func NewGRPCServer(lc fx.Lifecycle, logger *zap.Logger, p GRPCServerParams) *grpc.Server {
	logger = logger.Named("grpc")

	panicHandler := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic recover: %v", p)
	}

	recoverOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(panicHandler),
	}

	setGRPCLogger(logger)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{*certs.GlobalCertificate()},
		ClientAuth:   tls.NoClientCert,
	})

	interceptor := session.NewAuthInterceptor(p.JWTManager)

	grpcServer := grpc.NewServer(
		grpc.Creds(creds),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(recoverOpts...),
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_auth.UnaryServerInterceptor(interceptor.Interceptor()),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(recoverOpts...),
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.StreamServerInterceptor(logger),
			grpc_auth.StreamServerInterceptor(interceptor.Interceptor()),
		),
	)

	proto.RegisterChatServiceServer(grpcServer, p.Service1)
	proto.RegisterSessionServiceServer(grpcServer, p.SessionService)
	proto.RegisterAccountsServiceServer(grpcServer, p.Accounts)

	addr := ":8080"

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			logger.Info("Starting GRP server", zap.String("addr", addr))
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
