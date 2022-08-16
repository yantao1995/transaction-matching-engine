package common

import "sync/atomic"

var (
	orderAtomicSource int32 = 0 //订单原子自增量
)

const (
	int32MaxValueRotation int32 = 2 << 16 //自增归零量
)

//获取原子自增值
func GetAtomicIncrSeq() int {
	atomic.CompareAndSwapInt32(&orderAtomicSource, int32MaxValueRotation, 0)
	return int(atomic.AddInt32(&orderAtomicSource, 1))
}
