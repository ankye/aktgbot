package main

import (
	"fmt"
	"testing"
)

func TestBinance(t *testing.T) {
	m := Binance("BTCUSDT", BTC)
	fmt.Println(m.Name)
}
