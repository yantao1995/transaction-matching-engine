package grpc

import (
	context "context"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type implementedMatchServiceServer struct {
}

func NewImplementedMatchServiceServer() *implementedMatchServiceServer {
	return &implementedMatchServiceServer{}
}

func (*implementedMatchServiceServer) AddOrder(context.Context, *AddOrderRequest) (*CommonResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method AddOrder not implemented")
}
func (*implementedMatchServiceServer) CancelOrder(context.Context, *CancelOrderRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
