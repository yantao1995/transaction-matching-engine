package pool

import (
	"transaction-matching-engine/models"

	"github.com/shopspring/decimal"
	"github.com/yantao1995/ds/skiplist"
)

// 订单池
type pool struct {
	sl *skiplist.SkipList
}

func NewPool(cmp skiplist.CompareAble) *pool {
	sl, _ := skiplist.New(cmp, skiplist.WithAllowTheSameKey(false))
	return &pool{
		sl: sl,
	}
}

// 生成订单排序key
func (p *pool) generateSortKey(order *models.Order) *models.SortKey {
	price, _ := decimal.NewFromString(order.Price)
	//amount, _ := decimal.NewFromString(order.Amount)
	return &models.SortKey{
		Price:         price,
		TimeUnixMilli: order.TimeUnixMilli,
		//	Amount:        amount,
		Id: order.Id,
	}
}

// 写入订单池
func (p *pool) Insert(order *models.Order) {
	p.sl.Insert(p.generateSortKey(order), order)
}

// 获取订单池内订单深度
func (p *pool) GetOrderDepth() int {
	return p.sl.GetLength()
}

// 更新指定档位的数据
func (p *pool) UpdateDataByDepth(rk int, order *models.Order) bool {
	if rk <= p.GetOrderDepth() {
		return p.sl.UpdateByRank(rk, order)
	}
	return false
}

// 删除指定档位的数据
func (p *pool) DeleteByDepth(rk int) bool {
	if rk <= p.GetOrderDepth() {
		return p.sl.DeleteByRank(rk)
	}
	return false
}

// 删除指定订单
func (p *pool) DeleteByOrder(order *models.Order) bool {
	if p.GetOrderDepth() > 0 {
		return p.sl.DeleteBatchByKey(p.generateSortKey(order))
	}
	return false
}

// 获取订单池指定档位的数据
func (p *pool) GetDepthData(rk int) *models.Order {
	if rk <= p.GetOrderDepth() {
		return p.sl.GetByRank(rk).(*models.Order)
	}
	return nil
}

// 获取订单池内所有订单
func (p *pool) GetAllOrders() []*models.Order {
	all := p.sl.GetByRankRange(1, p.sl.GetLength())
	orders := make([]*models.Order, len(all))
	for k := range all {
		orders[k] = all[k].(*models.Order)
	}
	return orders
}
