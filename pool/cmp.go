package pool

import "transaction-matching-engine/models"

/*
	排序优先级：
		价格
		时间 [ 早 -> 晚 ]
		数量 [ 少 -> 多 ]  //能保证池内数量按序消耗之后仍然保持 [少->多] 的顺序而不用重新排序
		id   [ 小 -> 大 ]  id肯定不相等，所以不会出现相等的key
*/

/*
	卖单，跳表价格升序
*/

type asksCmp struct {
	val    int             //decimal比较结果
	ak, bk *models.SortKey //排序key
}

//实现比较接口
func (asks *asksCmp) Compare(a, b interface{}) int {
	asks.ak, asks.bk = a.(*models.SortKey), b.(*models.SortKey)
	if asks.val = asks.ak.Price.Cmp(asks.bk.Price); asks.val != 0 { //价格降序
		return asks.val
	}
	if asks.ak.TimeUnixMilli < asks.bk.TimeUnixMilli {
		return -1
	} else if asks.ak.TimeUnixMilli > asks.bk.TimeUnixMilli {
		return 1
	}
	if asks.val = asks.ak.Amount.Cmp(asks.bk.Amount); asks.val != 0 {
		return asks.val
	}
	if asks.ak.Id < asks.bk.Id {
		return -1
	} else if asks.ak.Id > asks.bk.Id {
		return 1
	}
	return 0
}

/*
	买单，跳表价格降序
*/

type bidsCmp struct {
	val    int             //decimal比较结果
	ak, bk *models.SortKey //排序key
}

//实现比较接口
func (bids *bidsCmp) Compare(a, b interface{}) int {
	bids.ak, bids.bk = a.(*models.SortKey), b.(*models.SortKey)
	if bids.val = bids.bk.Price.Cmp(bids.ak.Price); bids.val != 0 { //价格降序
		return bids.val
	}
	if bids.ak.TimeUnixMilli < bids.bk.TimeUnixMilli {
		return -1
	} else if bids.ak.TimeUnixMilli > bids.bk.TimeUnixMilli {
		return 1
	}
	if bids.val = bids.ak.Amount.Cmp(bids.bk.Amount); bids.val != 0 {
		return bids.val
	}
	if bids.ak.Id < bids.bk.Id {
		return -1
	} else if bids.ak.Id > bids.bk.Id {
		return 1
	}
	return 0
}
