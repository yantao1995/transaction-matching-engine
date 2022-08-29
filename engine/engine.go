package engine

import (
	"fmt"
	"strings"
	"sync"
	"transaction-matching-engine/models"
	"transaction-matching-engine/pool"
)

var (
	once        = sync.Once{}
	matchEngine *MatchEngine
)

type MatchEngine struct {
	pools map[string]*pool.MatchPool //每个交易对一个撮合池
}

func GetMatchEngine(pairs []string) *MatchEngine {
	once.Do(func() {
		matchEngine = &MatchEngine{
			pools: make(map[string]*pool.MatchPool),
		}
		for _, pair := range pairs {
			if _, ok := matchEngine.pools[pair]; !ok {
				matchEngine.pools[pair] = pool.NewMatchPool()
			}
		}
		matchEngine.printPairs()
		matchEngine.subscribeTrade()
	})
	return matchEngine
}

//输出所有的交易对
func (me *MatchEngine) printPairs() {
	pairs := make([]string, 0, len(me.pools))
	for pair := range me.pools {
		pairs = append(pairs, pair)
	}
	fmt.Printf("撮合池内支持的交易对: [%s]\n", strings.Join(pairs, ","))
}

//订阅成交   //推送至消息队列供业务测消费 ---
func (me *MatchEngine) subscribeTrade() {
	for pair, mp := range me.pools {
		go func(pair string, mp *pool.MatchPool) {
			for trade := range mp.Output() {
				fmt.Printf("新的成交!交易对:%s\r\n详细信息:%+v\n", pair, trade)
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
