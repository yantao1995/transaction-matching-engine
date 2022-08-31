package pool

import (
	"sync"
	"time"
	"transaction-matching-engine/models"
)

/*
	快照惰性更新，有请求时，判断上次更新时间，如果时间差大于当前计算出的最大时间差，则更新

	快照更新时间 计算公式
	update Millisecond =  snapshotUpdateTimeMillisecondBase +  ( len(orders) / increaseFactor )
	每 increaseFactor 个订单，增加更新时间1毫秒

	快照结构, 返回u方便业务对量化帐号进行筛选
	[[u,p,a]]
	[[用户id,价格,数量]]
*/

//盘口深度快照
type deepSnapshot struct {
	queryLock                          *sync.Mutex //保证更新深度和挂单不冲突
	bidsDeepSnapshot, asksDeepSnapshot [][3]string //订单深度快照
	lastUpdateTimeMilli                int64       //上次更新时间 毫秒
	snapshotUpdateTimeMillisecondBase  int64       //订单快照更新基础时间
	increaseFactor                     int64       //更新时间递增因子
	hasNewOrder                        bool        //这段时间有新订单进入 如果一段时间无新订单则不更新
}

func newDeepSnapshot() *deepSnapshot {
	return &deepSnapshot{
		queryLock:                         &sync.Mutex{},
		bidsDeepSnapshot:                  [][3]string{},
		asksDeepSnapshot:                  [][3]string{},
		lastUpdateTimeMilli:               0,
		snapshotUpdateTimeMillisecondBase: 200,
		increaseFactor:                    10000,
		hasNewOrder:                       false,
	}
}

//根据上次更新时间判断当前是否需要更新
func (d *deepSnapshot) IsNeedUpdate() bool {
	return d.hasNewOrder &&
		time.Now().UnixMilli()-d.lastUpdateTimeMilli > d.snapshotUpdateTimeMillisecondBase+
			int64(len(d.bidsDeepSnapshot)+len(d.asksDeepSnapshot))/d.increaseFactor
}

//更新快照
func (d *deepSnapshot) Update(bids, asks []*models.Order) {
	d.bidsDeepSnapshot = d.bidsDeepSnapshot[:0]
	for k := range bids {
		d.bidsDeepSnapshot = append(d.bidsDeepSnapshot, [3]string{bids[k].UserId, bids[k].Price, bids[k].Amount})
	}
	d.asksDeepSnapshot = d.asksDeepSnapshot[:0]
	for k := range asks {
		d.asksDeepSnapshot = append(d.asksDeepSnapshot, [3]string{asks[k].UserId, asks[k].Price, asks[k].Amount})
	}
	d.lastUpdateTimeMilli = time.Now().UnixMilli()
	d.hasNewOrder = false
}

//获取快照
func (d *deepSnapshot) GetSnapshot() ([][3]string, [][3]string) {
	return d.bidsDeepSnapshot, d.asksDeepSnapshot
}
