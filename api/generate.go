package proto

//go:generate protoc --go_out=../pkg/proto --go_opt=paths=source_relative --go-grpc_out=../pkg/proto --go-grpc_opt=paths=source_relative session_service.proto
//go:generate protoc --go_out=../pkg/proto --go_opt=paths=source_relative --go-grpc_out=../pkg/proto --go-grpc_opt=paths=source_relative main.proto
//go:generate protoc --go_out=../pkg/proto --go_opt=paths=source_relative --go-grpc_out=../pkg/proto --go-grpc_opt=paths=source_relative accounts_service.proto
