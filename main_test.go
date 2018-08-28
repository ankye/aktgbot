package main

import (
	"fmt"
	"testing"
)

func TestBinance(t *testing.T) {
	m := Binance("BTCUSDT", BTC)
	fmt.Println(m.Name)
}
func TestCoinex(t *testing.T) {
	m := coinex("BCHUSDT", BCC)
	fmt.Println(m.Last)
}
