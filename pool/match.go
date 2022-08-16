package pool

import "transaction-matching-engine/models"

type MatchPool struct {
	Bids, Asks *pool
}

// func NewMatchPool() *MatchPool {

// }

//FOK订单
func (m *MatchPool) FOKOrder(order *models.Order) {

}

//IOC订单
func (m *MatchPool) IOCOrder(order *models.Order) {

}

//GTC订单
func (m *MatchPool) GTCOrder(order *models.Order) {

}
