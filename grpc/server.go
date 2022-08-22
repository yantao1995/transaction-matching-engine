package grpc

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

func Run() {
	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor([]grpc_recovery.Option{
				grpc_recovery.WithRecoveryHandler(recoverInterceptor)}...),
			loggerInterceptor,
		),
	}

	server := grpc.NewServer(opts)
	lst, err := net.Listen("tcp", ":6666")
	if err != nil {
		panic("rpc listen err:" + err.Error())
	}

	RegisterMatchServiceServer(server, NewImplementedMatchServiceServer())

	fmt.Println("rpc running...")
	if err := server.Serve(lst); err != nil {
		panic("rpc Serve err:" + err.Error())
	}
	fmt.Println("rpc stopped.")
}
