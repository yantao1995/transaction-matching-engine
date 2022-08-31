package models

/*
	业务系统需要校验参数的合法性

	GTC	成交为止 :订单会一直有效，直到被成交或者取消。
	IOC	无法立即成交的部分就撤销 : 订单在失效前会尽量多的成交。
	FOK	无法全部立即成交就撤销 : 如果无法全部成交，订单会失效。
*/

//订单挂单    订单撤销仅会传Id/Pair过来,会重复使用该结构，但会修改 TimeInForce = CANCEL
type Order struct {
	Id            string `json:"id"`              //订单唯一标识
	UserId        string `json:"user_id"`         //用户id
	Pair          string `json:"pair"`            //交易对  作为多个撮合池的唯一标识
	Price         string `json:"price"`           //价格 (市价单无price)
	Amount        string `json:"amount"`          //数量
	Type          string `json:"type"`            //订单类型  [limit,market]  limit 限价单  market市价单	[market市价单走IOC成交]
	Side          string `json:"side"`            //订单方向  [buy,sell]  buy 买单 sell 卖单
	TimeInForce   string `json:"time_in_force"`   //订单有效期 [GTC,IOC,FOK] 说明见注释  限价单必传，市价单不传 (默认GTC)
	TimeUnixMilli int64  `json:"time_unix_milli"` //下单时间戳 毫秒
}
