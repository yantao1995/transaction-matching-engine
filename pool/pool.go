package pool

import (
	"transaction-matching-engine/models"

	"github.com/shopspring/decimal"
	"github.com/yantao1995/ds/skiplist"
)

//订单池
type pool struct {
	sl *skiplist.SkipList
}

func NewPool(cmp skiplist.CompareAble) *pool {
	sl, _ := skiplist.New(cmp, skiplist.WithAllowTheSameKey(false))
	return &pool{
		sl: sl,
	}
}

//生成订单排序key
func (p *pool) GenerateSortKey(order *models.Order) *models.SortKey {
	price, _ := decimal.NewFromString(order.Price)
	amount, _ := decimal.NewFromString(order.Amount)
	return &models.SortKey{
		Price:         price,
		TimeUnixMilli: order.TimeUnixMilli,
		Amount:        amount,
		Id:            order.Id,
	}
}

//写入撮合池
func (p *pool) Insert(order *models.Order) {
	p.sl.Insert(p.GenerateSortKey(order), order)
}

//获取撮合池内订单数量
func (p *pool) GetOrderLength() int {
	return p.sl.GetLength()
}

//获取撮合池指定档位的数据
func (p *pool) GetDepthData(rk int) *models.Order {
	if rk <= p.sl.GetLength() {
		return p.sl.GetByRank(rk).(*models.Order)
	}
	return nil
}
