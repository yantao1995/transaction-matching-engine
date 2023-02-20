package grpc

import (
	context "context"
	"encoding/json"
	"strings"
	"time"
	"transaction-matching-engine/common"
	"transaction-matching-engine/engine"
	"transaction-matching-engine/models"
)

// 交易对不区分大小写
type implementedMatchServiceServer struct {
	me      *engine.MatchEngine
	timeout time.Duration
}

// 参数由业务侧校验
func NewImplementedMatchServiceServer(pairs []string) *implementedMatchServiceServer {
	return &implementedMatchServiceServer{
		me:      engine.GetMatchEngine(pairs),
		timeout: time.Second,
	}
}

func (im *implementedMatchServiceServer) handleErr(err error) *CommonResponse {
	resp := &CommonResponse{}
	if err != nil {
		switch err {
		case common.ServerCancelErr, common.OrderHandleTimeoutErr:
			resp.Code = 500
		default:
			resp.Code = 400
		}
		resp.Msg = err.Error()
	}
	return resp
}

func (im *implementedMatchServiceServer) AddOrder(ctx context.Context, req *AddOrderRequest) (*CommonResponse, error) {
	order := &models.Order{
		Id:            req.GetId(),
		UserId:        req.GetUserId(),
		Pair:          strings.ToUpper(req.GetPair()),
		Price:         req.GetPrice(),
		Amount:        req.GetAmount(),
		Type:          req.GetType(),
		Side:          req.GetSide(),
		TimeInForce:   req.GetTimeInForce(),
		TimeUnixMilli: req.GetTimeUnixMilli(),
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, im.timeout)
	defer cancel()
	return im.handleErr(im.me.AddOrder(ctxTimeout, order)), nil
}

func (im *implementedMatchServiceServer) CancelOrder(ctx context.Context, req *CancelOrderRequest) (*CommonResponse, error) {
	order := &models.Order{
		Id:   req.GetId(),
		Pair: strings.ToUpper(req.GetPair()),
	}
	ctxTimeout, cancel := context.WithTimeout(ctx, im.timeout)
	defer cancel()
	return im.handleErr(im.me.AddOrder(ctxTimeout, order)), nil
}

func (im *implementedMatchServiceServer) QueryDeep(ctx context.Context, req *QueryDeepRequest) (*CommonResponse, error) {
	req.Pair = strings.ToUpper(req.GetPair())
	bids, asks, err := im.me.QueryDeep(req.GetPair())
	resp := im.handleErr(err)
	if err == nil {
		data := models.Deep{
			Pair:          req.GetPair(),
			TimeUnixMilli: time.Now().UnixMilli(),
			Bids:          bids,
			Asks:          asks,
		}
		resp.Data, _ = json.Marshal(data)
	}
	return resp, nil
}
