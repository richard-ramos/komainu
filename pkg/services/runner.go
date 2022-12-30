package services

import (
	"context"
	"fmt"

	"github.com/richard-ramos/komainu/pkg/config"
	"github.com/richard-ramos/komainu/pkg/grpcserver"
	"github.com/richard-ramos/komainu/pkg/services/accounts"
	"github.com/richard-ramos/komainu/pkg/services/service1"
	"github.com/richard-ramos/komainu/pkg/services/session"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func runServer(server *grpc.Server) {
	fmt.Println("RUNNING!")
}

func App(cfg *config.Config) *fx.App {
	ctx := func() context.Context {
		return context.Background()
	}

	return fx.New(
		fx.Supply(cfg),
		fx.Provide(
			fx.Annotate(ctx, fx.As(new(context.Context))),
			service1.New,
			session.NewJWTManager,
			session.NewSessionServer,
			accounts.NewAccountsServer,
			grpcserver.NewGRPCServer,
			zap.NewProduction,
		),
		fx.Invoke(runServer),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)
}
