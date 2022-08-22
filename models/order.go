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

/*
	maker 订单 ： 盘口订单为maker订单;成交时以maker订单价格成交
	taker 订单 ： 与盘口订单成交部分为taker订单，无法成交部分转入盘口时，将由taker转为maker
*/

//成交/撤销 推送  分别至不同队列
type Trade struct {
	Id             string `json:"id"`               //成交id
	Pair           string `json:"pair"`             //交易对  作为多个撮合池的唯一标识
	TakerUserId    string `json:"taker_user_id"`    //taker用户id
	TakerOrderId   string `json:"taker_order_id"`   //taker订单id
	MakerUserId    string `json:"maker_user_id"`    //taker用户id
	MakerOrderId   string `json:"maker_order_id"`   //maker订单id
	Price          string `json:"price"`            //价格
	Amount         string `json:"amount"`           //数量
	TakerOrderType string `json:"taker_order_type"` //taker订单类型  [limit,market]  limit 限价单  market市价单
	MakerOrderType string `json:"maker_order_type"` //maker订单类型  [limit,market]  limit 限价单  market市价单
	TakerOrderSide string `json:"taker_order_side"` //taker订单方向  [buy,sell]  buy 买单 sell 卖单
	MakerOrderSide string `json:"maker_order_side"` //maker订单方向  [buy,sell]  buy 买单 sell 卖单
	TimeUnixMilli  int64  `json:"time_unix_milli"`  //成交时间戳 毫秒
	Type           string `json:"type"`             //订单的成交类型  [sell,buy,cancel]  (cancel数据均在taker，虽然包含maker撤单)
}
