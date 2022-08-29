package grpc

import (
	context "context"
	"transaction-matching-engine/engine"
	"transaction-matching-engine/models"
)

type implementedMatchServiceServer struct {
	me *engine.MatchEngine
}

func NewImplementedMatchServiceServer(pairs []string) *implementedMatchServiceServer {
	return &implementedMatchServiceServer{
		me: engine.GetMatchEngine(pairs),
	}
}

func (im *implementedMatchServiceServer) AddOrder(ctx context.Context, req *AddOrderRequest) (*CommonResponse, error) {
	//参数由业务侧校验
	order := &models.Order{
		Id:            req.GetId(),
		UserId:        req.GetUserId(),
		Pair:          req.GetPair(),
		Price:         req.GetPrice(),
		Amount:        req.GetAmount(),
		Type:          req.GetType(),
		Side:          req.GetSide(),
		TimeInForce:   req.GetTimeInForce(),
		TimeUnixMilli: req.GetTimeUnixMilli(),
	}
	im.me.AddOrder(order)
	return &CommonResponse{}, nil
}

func (im *implementedMatchServiceServer) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CommonResponse, error) {
	order := &models.Order{
		Id:   req.GetId(),
		Pair: req.GetPair(),
	}
	im.me.CancelOrder(order)
	return &CommonResponse{}, nil
}
