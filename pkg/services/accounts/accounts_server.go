package accounts

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"path"

	"github.com/richard-ramos/komainu/pkg/accounts"
	"github.com/richard-ramos/komainu/pkg/config"
	"github.com/richard-ramos/komainu/pkg/persistence/sqlcipher"
	"github.com/richard-ramos/komainu/pkg/persistence/sqlite"
	"github.com/richard-ramos/komainu/pkg/proto"
	"github.com/richard-ramos/komainu/pkg/services/session"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountsServer struct {
	proto.AccountsServiceServer

	jwtManager *session.JWTManager
	cfg        *config.Config
	logger     *zap.Logger

	accountsDB *accounts.Persistence
}

func NewAccountsServer(cfg *config.Config, logger *zap.Logger, jwtManager *session.JWTManager) (*AccountsServer, error) {
	db, err := sqlite.Initialize(path.Join(cfg.PathDataDir(), "accounts.db"))
	if err != nil {
		return nil, err
	}

	return &AccountsServer{
		jwtManager: jwtManager,
		cfg:        cfg,
		logger:     logger.Named("accounts"),
		accountsDB: accounts.NewPersistence(db),
	}, nil
}

func (server *AccountsServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if fullMethodName == "/AccountsService/Create" {
		return ctx, nil
	}

	interceptor := session.NewAuthInterceptor(server.jwtManager)
	interceptFn := interceptor.Interceptor()
	return interceptFn(ctx)
}

func (server *AccountsServer) Create(ctx context.Context, req *proto.NewAccountRequest) (*proto.NewAccountResponse, error) {
	acc, err := accounts.NewAccount("", nil)
	if err != nil {
		server.logger.Error("generating account", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "cannot generate account")
	}

	logger := server.logger.With(zap.String("accountID", acc.ID))

	err = server.accountsDB.SaveAccount(acc)
	if err != nil {
		logger.Error("saving new account", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "cannot save new account")
	}

	// TODO: validate password?

	passwdBytes := sha256.Sum256([]byte(req.Password))
	password := hex.EncodeToString(passwdBytes[:])

	db, err := sqlcipher.Initialize(path.Join(server.cfg.PathDataDir(), acc.ID+".db"), password, acc.KDFIterations)
	if err != nil {
		// TODO: server.accountsDB.DeleteAccount(acc)
		server.logger.Error("initializing sqlcipher db", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "cannot generate account")
	}
	db.Close()

	logger.Info("account generated")

	return &proto.NewAccountResponse{
		UserId: acc.ID,
	}, nil
}
