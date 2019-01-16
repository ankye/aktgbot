package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gonethopper/libs/config"
	log "github.com/gonethopper/libs/logs"
	"github.com/gonethopper/libs/utils"
	"github.com/tidwall/gjson"
	tb "tg.robot/telebot"
)

const (
	BTC    = "BTC"
	BCH    = "BCH"
	LTC    = "LTC"
	ETH    = "ETH"
	BCHBTC = "BCHBTC"
	LTCBTC = "LTCBTC"
	ETHBTC = "ETHBTC"

	BITSTAMP = "Bitstamp"
	POLONIEX = "Poloniex"
	BITTREX  = "Bittrex"
	BITFINEX = "Bitfinex"
	BINANCE  = "Binance"
	COINEX   = "CoinEx"
)

//Market market struct
type Market struct {
	Name          string
	Trader        string
	Last          float64
	PercentChange float64
}

//Account account info
type Account struct {
	BTC    *Market
	BCH    *Market
	LTC    *Market
	ETH    *Market
	BCHBTC *Market
	LTCBTC *Market
	ETHBTC *Market
}

//Subscription 订阅通知
type Subscription struct {
	Chat     *tb.Chat
	Trader   string
	Type     int
	BCHPrice float64
	BTCPrice float64
	Duration int
	LastTime int
}

//LocalMilliscond LocalMilliscond
func LocalMilliscond() int64 {
	return time.Now().UnixNano() / 1e6
}

//LocalSecond LocalSecond
func LocalSecond() int {
	return int(LocalMilliscond() / 1000)
}

//NewSubscription create NewSubscription
func NewSubscription(trader string, t int, duration int) *Subscription {
	return &Subscription{
		Chat:     nil,
		Trader:   trader,
		Type:     t,
		Duration: duration,
		BTCPrice: 0,
		BCHPrice: 0,
		LastTime: LocalSecond(),
	}
}

//NewMarket create new Market  data
func NewMarket(name string, trader string, last float64, percentChange float64) *Market {

	return &Market{
		Name:          name,
		Trader:        trader,
		Last:          last,
		PercentChange: percentChange,
	}
}

//NewAccount create new Account  data
func NewAccount(name string) *Account {
	return &Account{
		BTC:    NewMarket(name, BTC, 0, 0),
		BCH:    NewMarket(name, BCH, 0, 0),
		LTC:    NewMarket(name, LTC, 0, 0),
		ETH:    NewMarket(name, ETH, 0, 0),
		BCHBTC: NewMarket(name, BCHBTC, 0, 0),
		LTCBTC: NewMarket(name, LTCBTC, 0, 0),
		ETHBTC: NewMarket(name, ETHBTC, 0, 0),
	}
}

//Output output string
func Output(rest ...*Market) string {
	str := ""
	for _, v := range rest {
		if v.Last > 10 {
			str = fmt.Sprintf("%s%s [%.2f] %.2f%%\n", str, v.Name, v.Last, v.PercentChange*100)
		} else {
			str = fmt.Sprintf("%s%s [%.4f] %.2f%%\n", str, v.Name, v.Last, v.PercentChange*100)
		}

	}
	return str
}

//Output2 output string
func Output2(rest ...*Market) string {
	str := ""
	for _, v := range rest {
		if v.Last > 10 {
			str = fmt.Sprintf("%s%s [%.2f] %.2f%%\n", str, v.Trader, v.Last, v.PercentChange*100)
		} else {
			str = fmt.Sprintf("%s%s [%.4f] %.2f%%\n", str, v.Trader, v.Last, v.PercentChange*100)
		}

	}
	return str
}

//Minimum Minimum
func Minimum(first *Market, rest ...*Market) *Market {
	minimum := first

	for _, v := range rest {
		if v.Last < minimum.Last {
			minimum = v
		}
	}
	return minimum
}

//Maximum Maximum
func Maximum(first *Market, rest ...*Market) *Market {
	maximum := first

	for _, v := range rest {
		if v.Last > maximum.Last {
			maximum = v
		}
	}
	return maximum
}

//HasNull HasNull
func HasNull(rest ...*Market) bool {

	for _, v := range rest {
		if v == nil {
			return true
		}
	}
	return false
}
func coinex(market string, trader string) *Market {
	url := "https://api.coinex.com/v1/market/ticker?market=" + market
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}
	last := gjson.GetBytes(body, "data.ticker.last").Float()

	open := gjson.GetBytes(body, "data.ticker.open").Float()

	percentChange := (last - open) / open

	return NewMarket(COINEX, trader, last, percentChange)
}
func bitfinex(market string, trader string) *Market {
	url := "https://api.bitfinex.com/v2/ticker/" + market
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}

	arr := gjson.ParseBytes(body).Array()
	if arr == nil {
		return nil
	}
	last := 0.0
	percent := 0.0
	if len(arr) > 6 {
		last = arr[6].Float()
		percent = arr[5].Float()
	}
	return NewMarket(BITFINEX, trader, last, percent)

}
func poloniex() *Account {
	url := "https://poloniex.com/public?command=returnTicker"
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}
	account := NewAccount(POLONIEX)
	account.BTC.Last = gjson.GetBytes(body, "USDT_BTC.last").Float()
	account.BCH.Last = gjson.GetBytes(body, "USDC_BCHABC.last").Float()
	account.LTC.Last = gjson.GetBytes(body, "USDT_LTC.last").Float()
	account.ETH.Last = gjson.GetBytes(body, "USDT_ETH.last").Float()
	account.BCHBTC.Last = gjson.GetBytes(body, "BTC_BCHABC.last").Float()
	account.LTCBTC.Last = gjson.GetBytes(body, "BTC_LTC.last").Float()
	account.ETHBTC.Last = gjson.GetBytes(body, "BTC_ETH.last").Float()

	account.BTC.PercentChange = gjson.GetBytes(body, "USDT_BTC.percentChange").Float()
	account.BCH.PercentChange = gjson.GetBytes(body, "USDC_BCHABC.percentChange").Float()
	account.LTC.PercentChange = gjson.GetBytes(body, "USDT_LTC.percentChange").Float()
	account.ETH.PercentChange = gjson.GetBytes(body, "USDT_ETH.percentChange").Float()
	account.BCHBTC.PercentChange = gjson.GetBytes(body, "BTC_BCHABC.percentChange").Float()
	account.LTCBTC.PercentChange = gjson.GetBytes(body, "BTC_LTC.percentChange").Float()
	account.ETHBTC.PercentChange = gjson.GetBytes(body, "BTC_ETH.percentChange").Float()

	return account
}

func bittrex(market string, trader string) *Market {
	resp, err := http.Get("https://bittrex.com/api/v1.1/public/getmarketsummary?market=" + market)
	if err != nil {
		// handle error
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}

	result := gjson.GetBytes(body, "result").Array()
	row := result[0].Map()

	last := row["Last"].Float()
	prev := row["PrevDay"].Float()

	percentChange := (last - prev) / prev

	return NewMarket(BITTREX, trader, last, percentChange)

}

func bitstamp(market string, trader string) *Market {
	resp, err := http.Get("https://www.bitstamp.net/api/v2/ticker/" + market + "/")
	if err != nil {
		// handle error
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}
	last := gjson.GetBytes(body, "last").Float()

	open := gjson.GetBytes(body, "open").Float()

	percentChange := (last - open) / open

	return NewMarket(BITSTAMP, trader, last, percentChange)

}

//Binance 币安价格查询
func Binance(market string, trader string) *Market {
	resp, err := http.Get("https://api.binance.com/api/v1/ticker/24hr?symbol=" + market)
	if err != nil {
		// handle error
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return nil
	}
	last := gjson.GetBytes(body, "lastPrice").Float()

	open := gjson.GetBytes(body, "openPrice").Float()

	percentChange := (last - open) / open

	return NewMarket(BINANCE, trader, last, percentChange)
}

var tgSubscription map[string]*Subscription
var subscriptionFile string
var bot *tb.Bot

func saveSubscription() {

	file, _ := os.OpenFile(subscriptionFile, os.O_RDWR|os.O_CREATE, 0777)
	defer file.Close()
	enc := gob.NewEncoder(file)
	if err := enc.Encode(tgSubscription); err != nil {
		fmt.Println(err)
	}

}
func loadSubscription(name string) {

	file, _ := os.OpenFile(subscriptionFile, os.O_RDWR|os.O_CREATE, 0777)
	defer file.Close()

	dec := gob.NewDecoder(file)
	err2 := dec.Decode(&tgSubscription)
	if err2 != nil {
		fmt.Println(err2)
		fmt.Println(tgSubscription)
		if tgSubscription == nil {
			tgSubscription = make(map[string]*Subscription)
			saveSubscription()
		}
		return
	}

}
func alert() {
	for {

		currentTime := LocalSecond()
		for k, sub := range tgSubscription {
			if currentTime-sub.LastTime > sub.Duration {
				chat := new(tb.Chat)
				if sub.Chat != nil {
					str, _ := json.Marshal(sub.Chat)
					_ = json.Unmarshal(str, chat)
					if sub.Type == 2 {
						btcm := bitstamp("btcusd", BTC)
						bchm := bitstamp("bchusd", BCH)
						if HasNull(bchm, btcm) {
							continue
						}
						if sub.BTCPrice > 0 && sub.BCHPrice > 0 {

							btcPercentChange := (btcm.Last - sub.BTCPrice) / btcm.Last
							bchPercentChange := (bchm.Last - sub.BCHPrice) / bchm.Last

							if btcPercentChange >= 0.07 || btcPercentChange <= -0.08 {

								msg := fmt.Sprintf("BTC价格跌幅 [%.2f]->[%.2f] [%.2f%%]", sub.BTCPrice, btcm.Last, btcPercentChange*100)
								if btcPercentChange > 0 {
									msg = fmt.Sprintf("BTC价格涨幅 [%.2f]->[%.2f] [%.2f%%]", sub.BTCPrice, btcm.Last, btcPercentChange*100)

								}
								bot.SendMessage(chat, msg, nil)
								sub.BTCPrice = btcm.Last
								saveSubscription()

								doBTC(chat)
							}
							if bchPercentChange >= 0.07 || bchPercentChange <= -0.08 {

								msg := fmt.Sprintf("BCH价格跌幅 [%.2f]->[%.2f] [%.2f%%]", sub.BCHPrice, bchm.Last, bchPercentChange*100)
								if bchPercentChange > 0 {
									msg = fmt.Sprintf("BCH价格涨幅 [%.2f]->[%.2f] [%.2f%%]", sub.BCHPrice, bchm.Last, bchPercentChange*100)

								}
								log.Info(msg)
								bot.SendMessage(chat, msg, nil)
								sub.BCHPrice = bchm.Last
								saveSubscription()
								doBCH(chat)
							}
						} else {
							sub.BCHPrice = bchm.Last
							sub.BTCPrice = btcm.Last
							saveSubscription()
							msg := fmt.Sprintf("订阅BCH,BTC行情大波动提醒成功，七上八下模式开启 BTC %.2f BCH %.2f", sub.BTCPrice, sub.BCHPrice)
							log.Info(msg)
							bot.SendMessage(chat, msg, nil)
						}
					} else {
						if sub.Trader == BTC {
							//fmt.Printf("alert btc price to  %d", sub.Chat.ID)
							doBTC(chat)
						} else if sub.Trader == BCH {
							//fmt.Printf("alert BCH price to  %d", sub.Chat.ID)
							doBCH(chat)

						}
					}
				}

				if tgSubscription[k] != nil {
					tgSubscription[k].LastTime = currentTime
				}
			}
		}
		select {

		case <-time.After(time.Millisecond * 1000):
			//	fmt.Printf("tick\n")
			continue
		}
	}
}
func doBTC(chat *tb.Chat) {
	last1 := bitstamp("btcusd", BTC)
	account := poloniex()
	last2 := account.BTC
	last3 := bittrex("USDT-BTC", BTC)
	last4 := bitfinex("tBTCUSD", BTC)
	last5 := Binance("BTCUSDT", BTC)
	last6 := coinex("BTCUSDT", BTC)
	if HasNull(last1, last2, last3, last4, last5, last6) {
		msg := "查询失败，请重试"
		log.Info(msg)
		bot.SendMessage(chat, msg, nil)
	} else {
		min := Minimum(last1, last2, last3, last4, last5, last6)
		max := Maximum(last1, last2, last3, last4, last5, last6)
		agiotage := max.Last - min.Last
		per := agiotage / min.Last * 100
		out := Output(last1, last2, last3, last4, last5, last6)
		msg := fmt.Sprintf("BTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
		log.Info(msg)
		bot.SendMessage(chat, msg, nil)
	}
}

func doBCH(chat *tb.Chat) {
	account := poloniex()
	last1 := account.BCH
	last2 := bittrex("USDT-BCH", BCH)
	last3 := bitfinex("tBABUSD", BCH)
	last4 := bitstamp("bchusd", BCH)
	last5 := Binance("BCHABCUSDT", BCH)
	last6 := coinex("BCHUSDT", BCH)
	if HasNull(last1, last2, last3, last4, last5, last6) {
		msg := "查询失败，请重试"
		log.Info(msg)
		bot.SendMessage(chat, msg, nil)
	} else {
		min := Minimum(last1, last2, last3, last4, last5, last6)
		max := Maximum(last1, last2, last3, last4, last5, last6)
		agiotage := max.Last - min.Last
		per := agiotage / min.Last * 100
		out := Output(last1, last2, last3, last4, last5, last6)
		msg := fmt.Sprintf("BCH \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
		log.Info(msg)
		bot.SendMessage(chat, msg, nil)
	}
}
func web() {
	r := gin.Default()
	r.GET("/coinex", func(c *gin.Context) {
		text := ""
		last1 := coinex("BTCUSDT", BTC)
		last2 := coinex("BCHUSDT", BCH)
		last3 := coinex("LTCUSDT", LTC)
		last4 := coinex("ETHUSDT", ETH)
		last5 := coinex("BTCBCH", BCHBTC)
		last6 := coinex("LTCBTC", LTCBTC)
		last7 := coinex("ETHBTC", ETHBTC)
		last8 := coinex("CETUSDT", "CETUSDT")
		if HasNull(last1, last2, last3, last4, last5, last6, last7, last8) {
			text = "查询失败，请重试"
		} else {
			out := Output2(last1, last2, last3, last4, last5, last6, last7, last8)
			msg := fmt.Sprintf("Coinex: \n%s\n", out)
			text = msg
		}
		c.String(http.StatusOK, text)
	})
	r.Run("0.0.0.0:9999") // listen and serve on 0.0.0.0:8080
}
func main() {

	c := NewConfig()

	err := config.ParseConfig(c, "./conf/bot.yml")
	if err != nil {
		return
	}

	logConfig, err := utils.LoadLogConfig("./conf/log.yml")
	if err != nil {
		log.Error("load log config failed.", err)
		return
	}
	c.Log = logConfig

	subscriptionFile = "config/subscription.gob"
	loadSubscription(subscriptionFile)

	tempBot, err := tb.NewBot(c.App.Botkey)
	if err != nil {
		log.Error(err)
	}
	bot = tempBot
	go alert()
	messages := make(chan tb.Message, 100)

	bot.Listen(messages, 10*time.Second)
	go web()
	for message := range messages {
		log.Info("%v", message.Text)
		res := strings.Split(message.Text, "@")
		if len(res) > 0 && len(res[0]) > 0 {
			arr := strings.Split(res[0], " ")
			if arr[0] == "/alertbtc" {
				key := fmt.Sprintf("%s-%d", BTC, message.Chat.ID)
				str, _ := json.Marshal(message.Chat)
				chat := new(tb.Chat)
				_ = json.Unmarshal(str, chat)

				ns := NewSubscription(BTC, 1, 3600)
				ns.Chat = chat
				tgSubscription[key] = ns
				saveSubscription()
				msg := "订阅btc提醒成功,间隔1小时"
				log.Info(msg)
				bot.SendMessage(message.Chat, msg, nil)
			} else if arr[0] == "/alertbch" {
				key := fmt.Sprintf("%s-%d", BCH, message.Chat.ID)
				str, _ := json.Marshal(message.Chat)
				chat := new(tb.Chat)
				_ = json.Unmarshal(str, chat)
				ns := NewSubscription(BCH, 1, 3600)
				ns.Chat = chat
				tgSubscription[key] = ns
				saveSubscription()

				msg := "订阅bch提醒成功,间隔1小时"
				log.Info(msg)
				bot.SendMessage(message.Chat, msg, nil)

			} else if arr[0] == "/alertrange78" {
				key := fmt.Sprintf("%s%s-%d", BTC, BCH, message.Chat.ID)
				str, _ := json.Marshal(message.Chat)
				chat := new(tb.Chat)
				_ = json.Unmarshal(str, chat)
				ns := NewSubscription(BCH, 2, 600)
				ns.Chat = chat
				tgSubscription[key] = ns

				btcm := bitstamp("btcusd", BTC)
				bchm := bitstamp("bchusd", BCH)
				ns.BTCPrice = btcm.Last
				ns.BCHPrice = bchm.Last
				saveSubscription()

				msg := fmt.Sprintf("订阅BCH,BTC行情大波动提醒成功，七上八下模式开启 BTC %.2f BCH %.2f", btcm.Last, bchm.Last)
				log.Info(msg)
				bot.SendMessage(message.Chat, msg, nil)
			} else if arr[0] == "/dalertbtc" {
				key := fmt.Sprintf("%s-%d", BTC, message.Chat.ID)

				delete(tgSubscription, key)
				saveSubscription()
				bot.SendMessage(message.Chat, "取消订阅btc成功,不再提醒", nil)
			} else if arr[0] == "/dalertbch" {
				key := fmt.Sprintf("%s-%d", BCH, message.Chat.ID)
				delete(tgSubscription, key)
				saveSubscription()
				bot.SendMessage(message.Chat, "取消订阅bch成功,不再提醒", nil)

			} else if arr[0] == "/dalertrange78" {
				key := fmt.Sprintf("%s%s-%d", BTC, BCH, message.Chat.ID)
				delete(tgSubscription, key)
				saveSubscription()

				msg := "取消订阅BCH,BTC行情大波动提醒成功，七上八下模式关闭"
				log.Info(msg)
				bot.SendMessage(message.Chat, msg, nil)

			} else if arr[0] == "/hi" {
				bot.SendMessage(message.Chat, "Hello, "+message.Sender.FirstName+" ! \ndonated bch adress : 32LSbGXhDjUie578wGFPVUhK2M7boNcTsB", nil)
			} else if arr[0] == "/btc" {
				str, _ := json.Marshal(message.Chat)
				chat := new(tb.Chat)
				_ = json.Unmarshal(str, chat)
				doBTC(chat)

			} else if arr[0] == "/bch" {
				str, _ := json.Marshal(message.Chat)
				chat := new(tb.Chat)
				_ = json.Unmarshal(str, chat)
				doBCH(chat)

			} else if arr[0] == "/ltc" {
				last1 := bitstamp("ltcusd", LTC)
				account := poloniex()
				last2 := account.LTC
				last3 := bittrex("USDT-LTC", LTC)
				last4 := bitfinex("tLTCUSD", LTC)
				last5 := Binance("LTCUSDT", LTC)
				last6 := coinex("LTCUSDT", LTC)
				if HasNull(last1, last2, last3, last4, last5, last6) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					min := Minimum(last1, last2, last3, last4, last5, last6)
					max := Maximum(last1, last2, last3, last4, last5, last6)
					agiotage := max.Last - min.Last
					per := agiotage / min.Last * 100
					out := Output(last1, last2, last3, last4, last5, last6)
					msg := fmt.Sprintf("LTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}

			} else if arr[0] == "/eth" {
				last1 := bitstamp("ethusd", ETH)
				account := poloniex()
				last2 := account.ETH
				last3 := bittrex("USDT-ETH", ETH)
				last4 := bitfinex("tETHUSD", ETH)
				last5 := Binance("ETHUSDT", ETH)
				last6 := coinex("ETHUSDT", ETH)
				if HasNull(last1, last2, last3, last4, last5, last6) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					min := Minimum(last1, last2, last3, last4, last5, last6)
					max := Maximum(last1, last2, last3, last4, last5, last6)
					agiotage := max.Last - min.Last
					per := agiotage / min.Last * 100
					out := Output(last1, last2, last3, last4, last5, last6)
					msg := fmt.Sprintf("ETH \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}

			} else if arr[0] == "/bchbtc" {

				account := poloniex()
				last1 := account.BCHBTC
				last2 := bittrex("BTC-BCH", BCHBTC)
				last3 := bitfinex("tBABBTC", BCHBTC)
				last4 := bitstamp("bchbtc", BCHBTC)
				last5 := Binance("BCHABCBTC", BCHBTC)

				if HasNull(last1, last2, last3, last4, last5) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					min := Minimum(last1, last2, last3, last4, last5)
					max := Maximum(last1, last2, last3, last4, last5)
					agiotage := max.Last - min.Last
					per := agiotage / min.Last * 100
					out := Output(last1, last2, last3, last4, last5)
					msg := fmt.Sprintf("BCHBTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}

			} else if arr[0] == "/ltcbtc" {
				last1 := bitstamp("ltcbtc", LTCBTC)
				account := poloniex()
				last2 := account.LTCBTC
				last3 := bittrex("BTC-LTC", LTCBTC)
				last4 := bitfinex("tLTCBTC", LTCBTC)
				last5 := Binance("LTCBTC", LTCBTC)
				last6 := coinex("LTCBTC", LTCBTC)
				if HasNull(last1, last2, last3, last4, last5, last6) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					min := Minimum(last1, last2, last3, last4, last5, last6)
					max := Maximum(last1, last2, last3, last4, last5, last6)
					agiotage := max.Last - min.Last
					per := agiotage / min.Last * 100
					out := Output(last1, last2, last3, last4, last5, last6)
					msg := fmt.Sprintf("LTCBTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}

			} else if arr[0] == "/ethbtc" {
				last1 := bitstamp("ethbtc", ETHBTC)
				account := poloniex()
				last2 := account.ETHBTC
				last3 := bittrex("BTC-ETH", ETHBTC)
				last4 := bitfinex("tETHBTC", ETHBTC)
				last5 := Binance("ETHBTC", ETHBTC)
				last6 := coinex("ETHBTC", ETHBTC)
				if HasNull(last1, last2, last3, last4, last5, last6) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					min := Minimum(last1, last2, last3, last4, last5, last6)
					max := Maximum(last1, last2, last3, last4, last5, last6)
					agiotage := max.Last - min.Last
					per := agiotage / min.Last * 100
					out := Output(last1, last2, last3, last4, last5, last6)
					msg := fmt.Sprintf("ETHBTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}

			} else if arr[0] == "/bitstamp" {
				last := bitstamp("btcusd", BTC)
				last2 := bitstamp("ltcusd", LTC)
				last3 := bitstamp("ethusd", ETH)
				last4 := bitstamp("ltcbtc", LTCBTC)
				last5 := bitstamp("ethbtc", ETHBTC)
				last6 := bitstamp("bchusd", BCH)
				last7 := bitstamp("bchbtc", BCHBTC)
				if HasNull(last, last2, last3, last4, last5, last6, last7) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					out := Output2(last, last2, last3, last4, last5, last6, last7)
					msg := fmt.Sprintf("Bitstamp: \n%s", out)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}

			} else if arr[0] == "/poloniex" {
				account := poloniex()
				if account == nil {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					out := Output2(account.BTC, account.BCH, account.LTC, account.ETH, account.BCHBTC, account.LTCBTC, account.ETHBTC)

					msg := fmt.Sprintf("Poloniex: \n%s\n", out)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}

			} else if arr[0] == "/bittrex" {

				last1 := bittrex("USDT-BTC", BTC)
				last2 := bittrex("USDT-BCH", BCH)
				last3 := bittrex("USDT-LTC", LTC)
				last4 := bittrex("USDT-ETH", ETH)
				last5 := bittrex("BTC-BCH", BCHBTC)
				last6 := bittrex("BTC-LTC", LTCBTC)
				last7 := bittrex("BTC-ETH", ETHBTC)
				if HasNull(last1, last2, last3, last4, last5, last6, last7) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					out := Output2(last1, last2, last3, last4, last5, last6, last7)
					msg := fmt.Sprintf("Bittrex: \n%s\n", out)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}
			} else if arr[0] == "/bitfinex" {

				last1 := bitfinex("tBTCUSD", BTC)
				last2 := bitfinex("tBABUSD", BCH)
				last3 := bitfinex("tLTCUSD", LTC)
				last4 := bitfinex("tETHUSD", ETH)
				last5 := bitfinex("tBABBTC", BCHBTC)
				last6 := bitfinex("tLTCBTC", LTCBTC)
				last7 := bitfinex("tETHBTC", ETHBTC)
				if HasNull(last1, last2, last3, last4, last5, last6, last7) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					out := Output2(last1, last2, last3, last4, last5, last6, last7)
					msg := fmt.Sprintf("Bitfinex: \n%s\n", out)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}
			} else if arr[0] == "/binance" {

				last1 := Binance("BTCUSDT", BTC)
				last2 := Binance("BCHABCUSDT", BCH)
				last3 := Binance("LTCUSDT", LTC)
				last4 := Binance("ETHUSDT", ETH)
				last5 := Binance("BCHABCBTC", BCHBTC)
				last6 := Binance("LTCBTC", LTCBTC)
				last7 := Binance("ETHBTC", ETHBTC)
				if HasNull(last1, last2, last3, last4, last5, last6, last7) {
					bot.SendMessage(message.Chat, "查询失败，请重试", nil)
				} else {
					out := Output2(last1, last2, last3, last4, last5, last6, last7)
					msg := fmt.Sprintf("Binance: \n%s\n", out)
					bot.SendMessage(message.Chat, msg, nil)
					log.Info(msg)
				}
			} else if arr[0] == "/coinex" {
				if len(arr) > 1 {
					last1 := coinex(arr[1], arr[1])
					if HasNull(last1) {
						bot.SendMessage(message.Chat, "查询失败，请重试", nil)
					} else {
						out := Output2(last1)
						msg := fmt.Sprintf("Coinex: \n%s\n", out)
						bot.SendMessage(message.Chat, msg, nil)
						log.Info(msg)
					}
				} else {

					last1 := coinex("BTCUSDT", BTC)
					last2 := coinex("BCHUSDT", BCH)
					last3 := coinex("LTCUSDT", LTC)
					last4 := coinex("ETHUSDT", ETH)
					last5 := coinex("BTCBCH", BCHBTC)
					last6 := coinex("LTCBTC", LTCBTC)
					last7 := coinex("ETHBTC", ETHBTC)
					last8 := coinex("CETUSDT", "CETUSDT")
					if HasNull(last1, last2, last3, last4, last5, last6, last7, last8) {
						bot.SendMessage(message.Chat, "查询失败，请重试", nil)
					} else {
						out := Output2(last1, last2, last3, last4, last5, last6, last7, last8)
						msg := fmt.Sprintf("Coinex: \n%s\n", out)
						bot.SendMessage(message.Chat, msg, nil)
						log.Info(msg)
					}
				}
			} else {
				bot.SendMessage(message.Chat, "你等着，我等会找着了给你", nil)
			}
		}
	}

}
