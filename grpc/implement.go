package grpc

import (
	context "context"
	"transaction-matching-engine/engine"

	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type implementedMatchServiceServer struct {
	me *engine.MatchEngine
}

func NewImplementedMatchServiceServer(pairs []string) *implementedMatchServiceServer {
	return &implementedMatchServiceServer{
		me: engine.NewMatchEngine(pairs),
	}
}

func (*implementedMatchServiceServer) AddOrder(context.Context, *AddOrderRequest) (*CommonResponse, error) {

	return nil, status.Errorf(codes.Unimplemented, "method AddOrder not implemented")
}
func (*implementedMatchServiceServer) CancelOrder(context.Context, *CancelOrderRequest) (*CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelOrder not implemented")
}
