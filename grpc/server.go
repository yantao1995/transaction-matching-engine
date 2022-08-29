package grpc

import (
	context "context"
	"fmt"
	"net"
	"transaction-matching-engine/common"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

type grpcServer struct {
	gs         *grpc.Server
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func Run(pairs []string) {
	common.ServerStatus.Add(1)
	defer common.ServerStatus.Done()

	opts := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor([]grpc_recovery.Option{
				grpc_recovery.WithRecoveryHandler(recoverInterceptor)}...),
			loggerInterceptor,
		),
	}

	ctx, cancelFunc := context.WithCancel(common.ServerStatus.Context())

	server := &grpcServer{
		gs:         grpc.NewServer(opts...),
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}

	lst, err := net.Listen("tcp", ":6666")
	if err != nil {
		panic("rpc listen err:" + err.Error())
	}

	RegisterMatchServiceServer(server.gs, NewImplementedMatchServiceServer(pairs))

	go gracefulStop(server)

	fmt.Println("rpc running...")
	if err := server.gs.Serve(lst); err != nil {
		panic("rpc Serve err:" + err.Error())
	}
	fmt.Println("rpc stopped.")
}

func gracefulStop(server *grpcServer) {
	<-server.ctx.Done()
	server.gs.GracefulStop()
}
