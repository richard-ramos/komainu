package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"time"

	"github.com/richard-ramos/komainu/pkg/certs"
	"github.com/richard-ramos/komainu/pkg/config"
	"github.com/richard-ramos/komainu/pkg/proto"
	"github.com/richard-ramos/komainu/pkg/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	go func() {
		services.App(&config.Config{}).Run()
	}()

	time.Sleep(1 * time.Second)

	certificateStr := certs.PublicTLSCertPEM()
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(certificateStr)) {
		panic("failed to add server CA's certificate")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	conn, err := grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	clientA := proto.NewSessionServiceClient(conn)
	result, err := clientA.Logout(context.Background(), &emptypb.Empty{})
	fmt.Println(result, err)
	fmt.Println("---------------------------------------------")

	client0 := proto.NewSessionServiceClient(conn)
	resultf, err := client0.Login(context.Background(), &proto.LoginRequest{
		UserId:   "abc",
		Password: "123",
	})
	fmt.Println(resultf, err)

	auth := TokenAuth(resultf.AccessToken)

	conn1, err := grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
		grpc.WithPerRPCCredentials(auth),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := proto.NewChatServiceClient(conn1)
	result2, err := client.GetAllUsers(context.Background(), &proto.Empty{})
	fmt.Println(result2, err)

	client5 := proto.NewAccountsServiceClient(conn1)
	result3, err := client5.Create(context.Background(), &proto.NewAccountRequest{
		Password: "ABC",
	})
	fmt.Println("====================================", result3, err)

	fmt.Println("LOGOUT")
	client4 := proto.NewSessionServiceClient(conn1)
	result5, err := client4.Logout(context.Background(), &emptypb.Empty{})
	fmt.Println(result5, err)
}
