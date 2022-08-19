package pool

import (
	"fmt"
	"testing"
	"time"
	"transaction-matching-engine/models"
)

func TestGTCPool(t *testing.T) {
	mpl := NewMatchPool()
	order := &models.Order{
		Id:            "1",
		UserId:        "1",
		Pair:          "BTC-USDT",
		Price:         "1000",
		Amount:        "100",
		Type:          "limit",
		Side:          "buy",
		TimeInForce:   "GTC",
		TimeUnixMilli: time.Now().UnixMilli(),
	}
	mpl.Input(order)
	order = &models.Order{
		Id:            "2",
		UserId:        "1",
		Pair:          "BTC-USDT",
		Price:         "2000",
		Amount:        "120",
		Type:          "limit",
		Side:          "buy",
		TimeInForce:   "GTC",
		TimeUnixMilli: time.Now().UnixMilli(),
	}
	mpl.Input(order)
	order = &models.Order{
		Id:            "3",
		UserId:        "2",
		Pair:          "BTC-USDT",
		Price:         "1500",
		Amount:        "100",
		Type:          "limit",
		Side:          "sell",
		TimeInForce:   "GTC",
		TimeUnixMilli: time.Now().UnixMilli(),
	}
	mpl.Input(order)
	order = &models.Order{
		Id:            "4",
		UserId:        "2",
		Pair:          "BTC-USDT",
		Price:         "900",
		Amount:        "30",
		Type:          "limit",
		Side:          "sell",
		TimeInForce:   "GTC",
		TimeUnixMilli: time.Now().UnixMilli(),
	}
	mpl.Input(order)
	ch := mpl.Output()
	for trade := range ch {
		fmt.Printf("成交:%+v\n", trade)
	}
}
