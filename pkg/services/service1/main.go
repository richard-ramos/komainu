package service1

import (
	"context"
	"fmt"

	"github.com/richard-ramos/komainu/pkg/proto"
)

type Service1 struct {
	proto.ChatServiceServer

	ctx    context.Context
	cancel context.CancelFunc
}

func (s *Service1) Join(ctx context.Context, user *proto.User) (*proto.JoinResponse, error) {
	return &proto.JoinResponse{
		Msg:   "ABC",
		Error: 0,
	}, nil
}

func (s *Service1) SendMsg(ctx context.Context, msg *proto.ChatMessage) (*proto.Empty, error) {
	return &proto.Empty{}, nil
}

func (s *Service1) ReceiveMsg(empty *proto.Empty, stream proto.ChatService_ReceiveMsgServer) error {
	return nil
}

func (s *Service1) GetAllUsers(ctx context.Context, empty *proto.Empty) (*proto.UserList, error) {
	fmt.Println("Calling getAllUsers")

	fmt.Println("USERID:", ctx.Value("userID"))

	return &proto.UserList{
		Users: []*proto.User{
			{
				Id:   "A",
				Name: "BCD",
			},
		},
	}, nil
}

func New(ctx context.Context) (*Service1, error) {
	fmt.Println("Creating service1")
	ctx, cancel := context.WithCancel(ctx)

	return &Service1{
		ctx:    ctx,
		cancel: cancel,
	}, nil
}
