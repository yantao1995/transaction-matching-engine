package models

import "github.com/shopspring/decimal"

/*
	可自定义增加排序逻辑
	例如:需要增加大客户优先级，则可以修改 SortKey 的排序优先级
		增加字段 UserId
		使用map存储 [UserId: 大客户分数] 如： m = map[10001:100,10002:99]
		asks与bids中调用Compare,增加分值比较逻辑即可: m[userId1] <=> m[userId2]
*/

//跳表排序key  ,优先级分别为 price,TimeUnixMilli,amount,id
type SortKey struct {
	Price         decimal.Decimal //价格  根据 asks/bids 类型升降排序
	TimeUnixMilli int64           //下单时间  时间早的排前面
	Amount        decimal.Decimal //数量	数量大的排前面
	Id            string          //订单唯一标识   增加订单号比较，保证排序key不重复
}
