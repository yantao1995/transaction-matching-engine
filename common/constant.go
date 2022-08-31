package common

import (
	"errors"
	"sync/atomic"
)

var (
	orderAtomicSource     int32 = 0                                               //订单原子自增量
	ServerCancelErr       error = errors.New("server done")                       //系统正在关闭
	OrderHandleTimeoutErr error = errors.New("order handle timeout")              //订单处理超时
	NotOpenMatchPoolErr   error = errors.New("this pair not open the match pool") //未配置交易撮合池
)

const (
	int32MaxValueRotation int32 = 1 << 20 //自增归零量    //1048576   fmt.Printf("%07d", 10) 补齐
	// 订单类型基础
	TypeOrderCancel string = "cancel"
	SideOrderBuy    string = "buy"
	SideOrderSell   string = "sell"
	TypeOrderLimit  string = "limit"
	TypeOrderMarket string = "market"
	TimeInForceGTC  string = "GTC"
	TimeInForceIOC  string = "IOC"
	TimeInForceFOK  string = "FOK"
)

//获取原子自增值
func GetAtomicIncrSeq() int {
	atomic.CompareAndSwapInt32(&orderAtomicSource, int32MaxValueRotation, 0)
	return int(atomic.AddInt32(&orderAtomicSource, 1))
}
