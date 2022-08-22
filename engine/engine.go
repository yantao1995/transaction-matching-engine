package engine

import (
	"transaction-matching-engine/models"
	"transaction-matching-engine/pool"
)

type MatchEngine struct {
	pools pool.MatchPool
}

func (meg *MatchEngine) Run() {

}

//挂单
func (meg *MatchEngine) AddOrder(order *models.Order) {

}

//撤单
func (meg *MatchEngine) CancelOrder(order *models.Order) {

}

//订阅成交
func (meg *MatchEngine) SubscribeTrade(order *models.Order) {

}
