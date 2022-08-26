package pool

import (
	"context"
	"fmt"
	"time"
	"transaction-matching-engine/common"
	"transaction-matching-engine/models"

	"github.com/shopspring/decimal"
)

type MatchPool struct {
	orderMap   map[string]*models.Order //存储订单池内的订单,撤销订单时
	orderChan  chan *models.Order       //订单输入
	bids, asks *pool                    //买卖盘撮合池
	tradeChan  chan *models.Trade       //订单成交
	ctx        context.Context
}

func NewMatchPool() *MatchPool {
	matchPool := &MatchPool{
		orderMap:  make(map[string]*models.Order),
		orderChan: make(chan *models.Order, 10000),
		bids:      NewPool(&bidsCmp{}),
		asks:      NewPool(&asksCmp{}),
		tradeChan: make(chan *models.Trade, 10000),
		ctx:       context.Background(),
	}
	go matchPool.run()
	return matchPool
}

func (m *MatchPool) run() {
	for orderPtr := range m.orderChan {
		if orderPtr.Type == common.TypeOrderMarket { //市价单走IOC挂单
			orderPtr.TimeInForce = common.TimeInForceIOC
		}
		switch orderPtr.TimeInForce {
		case common.TimeInForceFOK:
			m.orderFOK(orderPtr)
		case common.TimeInForceIOC:
			m.orderIOC(orderPtr)
		case common.TimeInForceGTC:
			m.orderGTC(orderPtr)
		default: //默认Cancel
			m.orderCancel(orderPtr)
		}
	}
}

//订单输入
func (m *MatchPool) Input(order *models.Order) error {
	select {
	case <-m.ctx.Done():
		return common.ServerCancelErr
	case m.orderChan <- order:
		return nil
	case <-time.After(time.Second):
		return common.OrderHandleTimeoutErr
	}
}

//成交输出
func (m *MatchPool) Output() <-chan *models.Trade {
	return m.tradeChan
}

//FOK	无法全部立即成交就撤销 : 如果无法全部成交，订单会失效。
func (m *MatchPool) orderFOK(order *models.Order) {
	rival := m.bids
	if order.Side == common.SideOrderBuy {
		rival = m.asks
	}
	needAmount, _ := decimal.NewFromString(order.Amount)
	canDealAmount := m.getCanDealAmount(rival, order, order.Type)

	//撤单
	if needAmount.Cmp(canDealAmount) == 1 { // need > can
		m.handleTrade(order, nil, nil, common.TypeOrderCancel, 1)
		return
	}

	//成交
	for ; needAmount.Cmp(decimal.Zero) == 1; needAmount, _ = decimal.NewFromString(order.Amount) {
		m.handleTrade(order, rival.GetDepthData(1), rival, order.Side, 1)
	}
}

//IOC	无法立即成交的部分就撤销 : 订单在失效前会尽量多的成交。
func (m *MatchPool) orderIOC(order *models.Order) {
	rival := m.bids
	if order.Side == common.SideOrderBuy {
		rival = m.asks
	}
	needAmount, _ := decimal.NewFromString(order.Amount)
	canDealAmount := m.getCanDealAmount(rival, order, order.Type)
	cancelAmount := decimal.Zero

	//计算需要撤销的数量
	if needAmount.Cmp(canDealAmount) == 1 { // need > can
		cancelAmount = needAmount.Sub(canDealAmount)
	}

	//计算需要成交的数量并成交
	if canDealAmount.Cmp(decimal.Zero) == 1 { //能成交的数量大于0

		if canDealAmount.Cmp(needAmount) == -1 { //重新生成数量
			order.Amount = canDealAmount.String()
			needAmount, _ = decimal.NewFromString(order.Amount)
		}

		for ; needAmount.Cmp(decimal.Zero) == 1; needAmount, _ = decimal.NewFromString(order.Amount) {
			m.handleTrade(order, rival.GetDepthData(1), rival, order.Side, 1)
		}
	}

	//撤单
	if cancelAmount.Cmp(decimal.Zero) == 1 {
		order.Amount = cancelAmount.String()
		m.handleTrade(order, nil, nil, common.TypeOrderCancel, 1)
	}
}

//GTC	成交为止 :订单会一直有效，直到被成交或者取消。
func (m *MatchPool) orderGTC(order *models.Order) {
	self, rival := m.asks, m.bids
	if order.Side == common.SideOrderBuy {
		self, rival = m.bids, m.asks
	}
	needAmount, _ := decimal.NewFromString(order.Amount)
	canDealAmount := m.getCanDealAmount(rival, order, order.Type)
	inputAmount := decimal.Zero

	//计算需要进入撮合池数量
	if needAmount.Cmp(canDealAmount) == 1 { // need >= can
		inputAmount = needAmount.Sub(canDealAmount)
	}

	//计算需要成交的数量并成交
	if canDealAmount.Cmp(decimal.Zero) == 1 { //能成交的数量 > 0

		if canDealAmount.Cmp(needAmount) == -1 { //重新生成数量
			order.Amount = canDealAmount.String()
			needAmount, _ = decimal.NewFromString(order.Amount)
		}

		for ; needAmount.Cmp(decimal.Zero) == 1; needAmount, _ = decimal.NewFromString(order.Amount) {
			m.handleTrade(order, rival.GetDepthData(1), rival, order.Side, 1)
		}
	}

	//进入撮合池
	if inputAmount.Cmp(decimal.Zero) == 1 {
		order.Amount = inputAmount.String() //重新生成数量
		self.Insert(order)
		m.orderMap[order.Id] = order
	}
}

//Cancel订单
func (m *MatchPool) orderCancel(order *models.Order) {
	order, ok := m.orderMap[order.Id]
	if ok {
		self := m.asks
		if order.Side == common.SideOrderBuy {
			self = m.bids
		}
		m.handleTrade(order, nil, self, common.TypeOrderCancel, 1)
		delete(m.orderMap, order.Id)
		self.DeleteByOrder(order)
	}
}

//获取能成交的数量   side,price为taker的状态
func (m *MatchPool) getCanDealAmount(p *pool, order *models.Order, orderType string) decimal.Decimal {
	canDeal := decimal.NewFromInt(0)
	orderPrice, _ := decimal.NewFromString(order.Price)
	orderAmount, _ := decimal.NewFromString(order.Amount)
	for rk := 1; rk <= p.GetOrderDepth() && orderAmount.Cmp(canDeal) == 1; rk++ {
		data := p.GetDepthData(rk)
		currentPrice, _ := decimal.NewFromString(data.Price)
		currentAmount, _ := decimal.NewFromString(data.Amount)
		if orderType == common.TypeOrderLimit && //限价单需要判断价格	不满足的条件
			((order.Side == common.SideOrderBuy && orderPrice.Cmp(currentPrice) == -1) || //买价 小于 卖一档价
				(order.Side == common.SideOrderSell && orderPrice.Cmp(currentPrice) == 1)) { //卖价 大于 买一档价
			return canDeal
		}
		canDeal = canDeal.Add(currentAmount)
	}
	return canDeal
}

//处理成交/撤销 trade     pl为对手盘撮合池
func (m *MatchPool) handleTrade(taker, maker *models.Order, pl *pool, tradeType string, rk int) {
	nowUnixMilli := time.Now().UnixMilli()
	trade := &models.Trade{
		Id:             fmt.Sprintf("%d%07d", nowUnixMilli, common.GetAtomicIncrSeq()),
		Pair:           taker.Pair,
		TakerUserId:    taker.UserId,
		TakerOrderId:   taker.Id,
		Price:          taker.Pair,
		Amount:         taker.Amount,
		TakerOrderType: taker.Type,
		TakerOrderSide: taker.Side,
		TimeUnixMilli:  nowUnixMilli,
		Type:           tradeType,
	}
	if tradeType != common.TypeOrderCancel {
		trade.Price = maker.Price
		trade.MakerUserId = maker.UserId
		trade.MakerOrderId = maker.Id
		trade.MakerOrderType = maker.Type
		trade.MakerOrderSide = maker.Side
		//计算amount
		takerAmountDecimal, _ := decimal.NewFromString(taker.Amount)
		makerAmountDecimal, _ := decimal.NewFromString(maker.Amount)
		if takerAmountDecimal.Cmp(makerAmountDecimal) > -1 { // taker >= maker
			taker.Amount = takerAmountDecimal.Sub(makerAmountDecimal).String()
			trade.Amount = maker.Amount
			delete(m.orderMap, maker.Id)
			pl.DeleteByDepth(rk)
		} else { // taker < maker
			taker.Amount = "0"
			maker.Amount = makerAmountDecimal.Sub(takerAmountDecimal).String()
			m.orderMap[maker.Id] = maker
			pl.UpdateDataByDepth(rk, maker)
		}
	}
	m.tradeChan <- trade
}
