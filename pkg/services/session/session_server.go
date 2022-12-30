package session

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/richard-ramos/komainu/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SessionServer struct {
	proto.SessionServiceServer
	sync.Mutex

	jwtManager *JWTManager
	loggedIn   bool
}

var ErrLoggedIn = errors.New("already logged in")
var ErrNotLoggedIn = errors.New("not logged in")

func NewSessionServer(jwtManager *JWTManager) *SessionServer {
	return &SessionServer{
		jwtManager: jwtManager,
	}
}

func (server *SessionServer) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	server.Lock()
	defer server.Unlock()

	if fullMethodName == "/SessionService/Login" && server.loggedIn {
		return nil, ErrLoggedIn
	}

	if fullMethodName == "/SessionService/Logout" {
		if !server.loggedIn {
			return nil, ErrNotLoggedIn
		}

		interceptor := NewAuthInterceptor(server.jwtManager)
		interceptFn := interceptor.Interceptor()
		return interceptFn(ctx)
	}

	return ctx, nil
}

func (server *SessionServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	server.Lock()
	defer server.Unlock()

	// TODO: verify user exists
	// return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	// TODO: verify user password
	// return nil, status.Errorf(codes.NotFound, "incorrect username/password")

	token, err := server.jwtManager.Generate(req.UserId)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}

	// TODO: DB login

	server.loggedIn = true

	res := &proto.LoginResponse{AccessToken: token}
	return res, nil
}

func (server *SessionServer) Logout(context.Context, *empty.Empty) (*empty.Empty, error) {
	server.Lock()
	defer server.Unlock()

	// TODO: DB logout

	server.loggedIn = false

	return &emptypb.Empty{}, nil
}
