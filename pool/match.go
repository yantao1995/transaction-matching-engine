package pool

import (
	"context"
	"fmt"
	"time"
	"transaction-matching-engine/common"
	"transaction-matching-engine/models"

	"github.com/shopspring/decimal"
)

//撮合池
type MatchPool struct {
	orderMap     map[string]*models.Order //存储订单池内的订单,撤销订单时
	orderChan    chan *models.Order       //订单输入
	bids, asks   *pool                    //买卖盘订单池
	tradeChan    chan *models.Trade       //订单成交
	ctx          context.Context
	cancelFunc   context.CancelFunc
	deepSnapshot *deepSnapshot
}

func NewMatchPool() *MatchPool {
	ctx, cancelFunc := context.WithCancel(common.ServerStatus.Context())
	matchPool := &MatchPool{
		orderMap:     make(map[string]*models.Order),
		orderChan:    make(chan *models.Order, 10000),
		bids:         NewPool(&bidsCmp{}),
		asks:         NewPool(&asksCmp{}),
		tradeChan:    make(chan *models.Trade, 10000),
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		deepSnapshot: newDeepSnapshot(),
	}
	go matchPool.listenSignal()
	go matchPool.run()
	return matchPool
}

/*
	业务相关  单线程模式
*/

//接收退出信号
func (m *MatchPool) listenSignal() {
	<-m.ctx.Done()
	close(m.orderChan)
}

//运行
func (m *MatchPool) run() {
	common.ServerStatus.Add(1)
	defer common.ServerStatus.Done()
	for orderPtr := range m.orderChan {
		m.deepSnapshot.queryLock.Lock()
		m.deepSnapshot.hasNewOrder = true
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
		m.deepSnapshot.queryLock.Unlock()
	}
}

//查询深度
func (m *MatchPool) QueryDeep(pair string) ([][3]string, [][3]string) {
	if m.deepSnapshot.IsNeedUpdate() {
		m.deepSnapshot.queryLock.Lock()
		if m.deepSnapshot.IsNeedUpdate() { //可能已经被上一个请求更新了
			m.deepSnapshot.Update(m.GetOrders())
		}
		m.deepSnapshot.queryLock.Unlock()
	}
	return m.deepSnapshot.GetSnapshot()
}

//订单输入  异步
func (m *MatchPool) Input(order *models.Order) error {
	select {
	case <-m.ctx.Done():
		return common.ServerCancelErr
	default:
		select {
		case m.orderChan <- order:
			return nil
		case <-time.After(time.Second):
			return common.OrderHandleTimeoutErr
		}
	}
}

//成交输出  异步
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

	//计算无法成交的数量
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

	if inputAmount.Cmp(decimal.Zero) == 1 {
		order.Amount = inputAmount.String()      //重新生成数量
		if order.Type == common.TypeOrderLimit { //限价单进入撮合池
			self.Insert(order)
			m.orderMap[order.Id] = order
		} else { //市价单撤单
			m.handleTrade(order, nil, nil, common.TypeOrderCancel, 1)
		}
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

//处理成交/撤销 trade     pl为对手盘订单池
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
	common.ServerStatus.Add(1)
	m.tradeChan <- trade
}

/*
	持久化相关
*/

//获取订单池   (bids,asks)
func (m *MatchPool) GetOrders() ([]*models.Order, []*models.Order) {
	return m.bids.GetAllOrders(), m.asks.GetAllOrders()
}

//写入订单池
func (m *MatchPool) SetOrders(bids, asks []*models.Order) {
	for k := range bids {
		m.bids.Insert(bids[k])
		m.orderMap[bids[k].Id] = bids[k]
	}
	for k := range asks {
		m.asks.Insert(asks[k])
		m.orderMap[asks[k].Id] = asks[k]
	}
}
