package engine

import (
	"fmt"
	"transaction-matching-engine/models"
	"transaction-matching-engine/pool"
)

type MatchEngine struct {
	pools map[string]*pool.MatchPool //每个交易对一个撮合池
}

func NewMatchEngine(pairs []string) *MatchEngine {
	me := &MatchEngine{
		pools: make(map[string]*pool.MatchPool),
	}
	for _, pair := range pairs {
		me.pools[pair] = pool.NewMatchPool()
	}
	me.subscribeTrade()
	return me
}

//订阅成交   //推送至消息队列供业务测消费 ---
func (me *MatchEngine) subscribeTrade() {
	for pair, mp := range me.pools {
		go func(pair string, mp *pool.MatchPool) {
			for trade := range mp.Output() {
				fmt.Println("pair:", pair, "\ttrade:", trade)
				/*
					TODO:消息推送
				*/
			}
		}(pair, mp)
	}
}

//挂单
func (me *MatchEngine) AddOrder(order *models.Order) {
	if pool, ok := me.pools[order.Pair]; ok {
		pool.Input(order)
		return
	}
	fmt.Println("[挂单]异常订单,id: ", order.Id, " 交易对: ", order.Pair)
}

//撤单
func (me *MatchEngine) CancelOrder(order *models.Order) {
	if pool, ok := me.pools[order.Pair]; ok {
		pool.Input(order)
		return
	}
	fmt.Println("[撤单]异常订单,id: ", order.Id, " 交易对: ", order.Pair)
}
