package services

import (
	"context"
	"fmt"

	"github.com/richard-ramos/komainu/grpcserver"
	"github.com/richard-ramos/komainu/services/service1"
	"github.com/richard-ramos/komainu/services/session"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func runServer(server *grpc.Server) {
	fmt.Println("RUNNING!")
}

func App() *fx.App {
	ctx := func() context.Context {
		return context.Background()
	}

	return fx.New(
		fx.Provide(
			fx.Annotate(ctx, fx.As(new(context.Context))),
			service1.New,
			session.NewJWTManager,
			session.NewSessionServer,
			grpcserver.NewGRPCServer,
		),
		fx.Invoke(runServer),
		// fx.NopLogger, // Disable FX logging
	)
}
