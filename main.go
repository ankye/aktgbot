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

	"github.com/tidwall/gjson"

	tb "./telebot"

	log "github.com/sirupsen/logrus"
)

const (
	BTC    = "BTC"
	BCC    = "BCH"
	LTC    = "LTC"
	ETH    = "ETH"
	BTCBCC = "BTCBCC"
	BTCLTC = "BTCLTC"
	BTCETH = "BTCETH"

	BITSTAMP = "Bitstamp"
	POLONIEX = "Poloniex"
	BITTREX  = "Bittrex"
	BITFINEX = "Bitfinex"
	BINANCE  = "Binance"
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
	BCC    *Market
	LTC    *Market
	ETH    *Market
	BTCBCC *Market
	BTCLTC *Market
	BTCETH *Market
}

//Subscription 订阅通知
type Subscription struct {
	Chat     *tb.Chat
	Trader   string
	Type     int
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
		BCC:    NewMarket(name, BCC, 0, 0),
		LTC:    NewMarket(name, LTC, 0, 0),
		ETH:    NewMarket(name, ETH, 0, 0),
		BTCBCC: NewMarket(name, BTCBCC, 0, 0),
		BTCLTC: NewMarket(name, BTCLTC, 0, 0),
		BTCETH: NewMarket(name, BTCETH, 0, 0),
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
	last := arr[6].Float()
	percent := arr[5].Float()
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
	account.BCC.Last = gjson.GetBytes(body, "USDT_BCH.last").Float()
	account.LTC.Last = gjson.GetBytes(body, "USDT_LTC.last").Float()
	account.ETH.Last = gjson.GetBytes(body, "USDT_ETH.last").Float()
	account.BTCBCC.Last = gjson.GetBytes(body, "BTC_BCH.last").Float()
	account.BTCLTC.Last = gjson.GetBytes(body, "BTC_LTC.last").Float()
	account.BTCETH.Last = gjson.GetBytes(body, "BTC_ETH.last").Float()

	account.BTC.PercentChange = gjson.GetBytes(body, "USDT_BTC.percentChange").Float()
	account.BCC.PercentChange = gjson.GetBytes(body, "USDT_BCH.percentChange").Float()
	account.LTC.PercentChange = gjson.GetBytes(body, "USDT_LTC.percentChange").Float()
	account.ETH.PercentChange = gjson.GetBytes(body, "USDT_ETH.percentChange").Float()
	account.BTCBCC.PercentChange = gjson.GetBytes(body, "BTC_BCH.percentChange").Float()
	account.BTCLTC.PercentChange = gjson.GetBytes(body, "BTC_LTC.percentChange").Float()
	account.BTCETH.PercentChange = gjson.GetBytes(body, "BTC_ETH.percentChange").Float()

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
var bchAlertValue float64
var btcAlertValue float64

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
						bchm := bitstamp("bchusd", BCC)
						if btcAlertValue > 0 && bchAlertValue > 0 {

							btcPercentChange := (btcAlertValue - btcm.Last) / btcAlertValue
							bchPercentChange := (bchAlertValue - bchm.Last) / bchAlertValue

							if btcPercentChange >= 0.07 || btcPercentChange <= -0.08 {

								msg := fmt.Sprintf("BTC价格跌幅 [%.2f]->[%.2f] [%.2f%%]", btcAlertValue, btcm.Last, btcPercentChange*100)
								if btcPercentChange > 0 {
									msg = fmt.Sprintf("BTC价格涨幅 [%.2f]->[%.2f] [%.2f%%]", btcAlertValue, btcm.Last, btcPercentChange*100)

								}
								bot.SendMessage(chat, msg, nil)
								btcAlertValue = btcm.Last

								doBTC(chat)
							}
							if bchPercentChange >= 0.07 || bchPercentChange <= -0.08 {

								msg := fmt.Sprintf("BCH价格跌幅 [%.2f]->[%.2f] [%.2f%%]", bchAlertValue, bchm.Last, bchPercentChange*100)
								if bchPercentChange > 0 {
									msg = fmt.Sprintf("BCH价格涨幅 [%.2f]->[%.2f] [%.2f%%]", bchAlertValue, bchm.Last, bchPercentChange*100)

								}
								bot.SendMessage(chat, msg, nil)
								bchAlertValue = bchm.Last
								doBCC(chat)
							}
						} else {
							bchAlertValue = bchm.Last
							btcAlertValue = btcm.Last
						}
					} else {
						if sub.Trader == BTC {
							//fmt.Printf("alert btc price to  %d", sub.Chat.ID)
							doBTC(chat)
						} else if sub.Trader == BCC {
							//fmt.Printf("alert bcc price to  %d", sub.Chat.ID)
							doBCC(chat)

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
	if HasNull(last1, last2, last3, last4, last5) {
		bot.SendMessage(chat, "查询失败，请重试", nil)
	} else {
		min := Minimum(last1, last2, last3, last4, last5)
		max := Maximum(last1, last2, last3, last4, last5)
		agiotage := max.Last - min.Last
		per := agiotage / min.Last * 100
		out := Output(last1, last2, last3, last4, last5)
		msg := fmt.Sprintf("BTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
		bot.SendMessage(chat, msg, nil)
	}
}

func doBCC(chat *tb.Chat) {
	account := poloniex()
	last1 := account.BCC
	last2 := bittrex("USDT-BCC", BCC)
	last3 := bitfinex("tBCHUSD", BCC)
	last4 := bitstamp("bchusd", BCC)
	last5 := Binance("BCCUSDT", BCC)
	if HasNull(last1, last2, last3, last4, last5) {
		bot.SendMessage(chat, "查询失败，请重试", nil)
	} else {
		min := Minimum(last1, last2, last3, last4, last5)
		max := Maximum(last1, last2, last3, last4, last5)
		agiotage := max.Last - min.Last
		per := agiotage / min.Last * 100
		out := Output(last1, last2, last3, last4, last5)
		msg := fmt.Sprintf("BCH \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
		bot.SendMessage(chat, msg, nil)
	}
}

func main() {

	subscriptionFile = "config/subscription.gob"
	loadSubscription(subscriptionFile)
	go alert()
	tempBot, err := tb.NewBot("429995834:AAH8T_JIn_tQ9fygYPCiOWppaLlO-buaEic")
	if err != nil {
		log.Fatalln(err)
	}
	bot = tempBot
	messages := make(chan tb.Message, 100)

	bot.Listen(messages, 10*time.Second)

	for message := range messages {
		//log.Infof("%v", message)
		arr := strings.Split(message.Text, "@")
		if arr[0] == "/alertbtc" {
			key := fmt.Sprintf("%s-%d", BTC, message.Chat.ID)
			str, _ := json.Marshal(message.Chat)
			chat := new(tb.Chat)
			_ = json.Unmarshal(str, chat)

			ns := NewSubscription(BTC, 1, 3600)
			ns.Chat = chat
			tgSubscription[key] = ns
			saveSubscription()
			bot.SendMessage(message.Chat, "订阅btc提醒成功,间隔1小时", nil)
		} else if arr[0] == "/alertbch" {
			key := fmt.Sprintf("%s-%d", BCC, message.Chat.ID)
			str, _ := json.Marshal(message.Chat)
			chat := new(tb.Chat)
			_ = json.Unmarshal(str, chat)
			ns := NewSubscription(BCC, 1, 3600)
			ns.Chat = chat
			tgSubscription[key] = ns
			saveSubscription()
			bot.SendMessage(message.Chat, "订阅bch提醒成功,间隔1小时", nil)
		} else if arr[0] == "/alertrange78" {
			key := fmt.Sprintf("%s%s-%d", BTC, BCC, message.Chat.ID)
			str, _ := json.Marshal(message.Chat)
			chat := new(tb.Chat)
			_ = json.Unmarshal(str, chat)
			ns := NewSubscription(BCC, 2, 600)
			ns.Chat = chat
			tgSubscription[key] = ns
			saveSubscription()
			btcm := bitstamp("btcusd", BTC)
			bchm := bitstamp("bchusd", BCC)
			btcAlertValue = btcm.Last
			bchAlertValue = bchm.Last

			bot.SendMessage(message.Chat, "订阅BCH,BTC行情大波动提醒成功，七上八下模式开启", nil)
		} else if arr[0] == "/dalertbtc" {
			key := fmt.Sprintf("%s-%d", BTC, message.Chat.ID)

			delete(tgSubscription, key)
			saveSubscription()
			bot.SendMessage(message.Chat, "取消订阅btc成功,不再提醒", nil)
		} else if arr[0] == "/dalertbch" {
			key := fmt.Sprintf("%s-%d", BCC, message.Chat.ID)
			delete(tgSubscription, key)
			saveSubscription()
			bot.SendMessage(message.Chat, "取消订阅bch成功,不再提醒", nil)

		} else if arr[0] == "/dalertrange78" {
			key := fmt.Sprintf("%s%s-%d", BTC, BCC, message.Chat.ID)
			delete(tgSubscription, key)
			saveSubscription()
			bchAlertValue = 0
			btcAlertValue = 0
			bot.SendMessage(message.Chat, "取消订阅BCH,BTC行情大波动提醒成功，七上八下模式关闭", nil)

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
			doBCC(chat)

		} else if arr[0] == "/ltc" {
			last1 := bitstamp("ltcusd", LTC)
			account := poloniex()
			last2 := account.LTC
			last3 := bittrex("USDT-LTC", LTC)
			last4 := bitfinex("tLTCUSD", LTC)
			last5 := Binance("LTCUSDT", LTC)
			if HasNull(last1, last2, last3, last4, last5) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				min := Minimum(last1, last2, last3, last4, last5)
				max := Maximum(last1, last2, last3, last4, last5)
				agiotage := max.Last - min.Last
				per := agiotage / min.Last * 100
				out := Output(last1, last2, last3, last4, last5)
				msg := fmt.Sprintf("LTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
				bot.SendMessage(message.Chat, msg, nil)
			}

		} else if arr[0] == "/eth" {
			last1 := bitstamp("ethusd", ETH)
			account := poloniex()
			last2 := account.ETH
			last3 := bittrex("USDT-ETH", ETH)
			last4 := bitfinex("tETHUSD", ETH)
			last5 := Binance("ETHUSDT", ETH)
			if HasNull(last1, last2, last3, last4, last5) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				min := Minimum(last1, last2, last3, last4, last5)
				max := Maximum(last1, last2, last3, last4, last5)
				agiotage := max.Last - min.Last
				per := agiotage / min.Last * 100
				out := Output(last1, last2, last3, last4, last5)
				msg := fmt.Sprintf("ETH \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
				bot.SendMessage(message.Chat, msg, nil)
			}

		} else if arr[0] == "/bchbtc" {

			account := poloniex()
			last1 := account.BTCBCC
			last2 := bittrex("BTC-BCC", BTCBCC)
			last3 := bitfinex("tBCHBTC", BTCBCC)
			last4 := bitstamp("bchbtc", BTCBCC)
			last5 := Binance("BCCBTC", BTCBCC)

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
			}

		} else if arr[0] == "/ltcbtc" {
			last1 := bitstamp("ltcbtc", BTCLTC)
			account := poloniex()
			last2 := account.BTCLTC
			last3 := bittrex("BTC-LTC", BTCLTC)
			last4 := bitfinex("tLTCBTC", BTCLTC)
			last5 := Binance("LTCBTC", BTCLTC)
			if HasNull(last1, last2, last3, last4, last5) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				min := Minimum(last1, last2, last3, last4, last5)
				max := Maximum(last1, last2, last3, last4, last5)
				agiotage := max.Last - min.Last
				per := agiotage / min.Last * 100
				out := Output(last1, last2, last3, last4, last5)
				msg := fmt.Sprintf("LTCBTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
				bot.SendMessage(message.Chat, msg, nil)
			}

		} else if arr[0] == "/ethbtc" {
			last1 := bitstamp("ethbtc", BTCETH)
			account := poloniex()
			last2 := account.BTCETH
			last3 := bittrex("BTC-ETH", BTCETH)
			last4 := bitfinex("tETHBTC", BTCETH)
			last5 := Binance("ETHBTC", BTCETH)
			if HasNull(last1, last2, last3, last4, last5) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				min := Minimum(last1, last2, last3, last4, last5)
				max := Maximum(last1, last2, last3, last4, last5)
				agiotage := max.Last - min.Last
				per := agiotage / min.Last * 100
				out := Output(last1, last2, last3, last4, last5)
				msg := fmt.Sprintf("ETHBTC \n%s\nmax: [%.2f] [%s]\nmin: [%.2f] [%s]\nagiotage:[%.2f][%.2f%%]", out, max.Last, max.Name, min.Last, min.Name, agiotage, per)
				bot.SendMessage(message.Chat, msg, nil)
			}

		} else if arr[0] == "/bitstamp" {
			last := bitstamp("btcusd", BTC)
			last2 := bitstamp("ltcusd", LTC)
			last3 := bitstamp("ethusd", ETH)
			last4 := bitstamp("ltcbtc", BTCLTC)
			last5 := bitstamp("ethbtc", BTCETH)
			last6 := bitstamp("bchusd", BCC)
			last7 := bitstamp("bchbtc", BTCBCC)
			if HasNull(last, last2, last3, last4, last5, last6, last7) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				out := Output2(last, last2, last3, last4, last5, last6, last7)
				msg := fmt.Sprintf("Bitstamp: \n%s", out)
				bot.SendMessage(message.Chat, msg, nil)
			}

		} else if arr[0] == "/poloniex" {
			account := poloniex()
			if account == nil {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				out := Output2(account.BTC, account.BCC, account.LTC, account.ETH, account.BTCBCC, account.BTCLTC, account.BTCETH)

				msg := fmt.Sprintf("Poloniex: \n%s\n", out)
				bot.SendMessage(message.Chat, msg, nil)
			}

		} else if arr[0] == "/bittrex" {

			last1 := bittrex("USDT-BTC", BTC)
			last2 := bittrex("USDT-BCC", BCC)
			last3 := bittrex("USDT-LTC", LTC)
			last4 := bittrex("USDT-ETH", ETH)
			last5 := bittrex("BTC-BCC", BTCBCC)
			last6 := bittrex("BTC-LTC", BTCLTC)
			last7 := bittrex("BTC-ETH", BTCETH)
			if HasNull(last1, last2, last3, last4, last5, last6, last7) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				out := Output2(last1, last2, last3, last4, last5, last6, last7)
				msg := fmt.Sprintf("Bittrex: \n%s\n", out)
				bot.SendMessage(message.Chat, msg, nil)
			}
		} else if arr[0] == "/bitfinex" {

			last1 := bitfinex("tBTCUSD", BTC)
			last2 := bitfinex("tBCHUSD", BCC)
			last3 := bitfinex("tLTCUSD", LTC)
			last4 := bitfinex("tETHUSD", ETH)
			last5 := bitfinex("tBCHBTC", BTCBCC)
			last6 := bitfinex("tLTCBTC", BTCLTC)
			last7 := bitfinex("tETHBTC", BTCETH)
			if HasNull(last1, last2, last3, last4, last5, last6, last7) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				out := Output2(last1, last2, last3, last4, last5, last6, last7)
				msg := fmt.Sprintf("Bitfinex: \n%s\n", out)
				bot.SendMessage(message.Chat, msg, nil)
			}
		} else if arr[0] == "/binance" {

			last1 := Binance("BTCUSDT", BTC)
			last2 := Binance("BCCUSDT", BCC)
			last3 := Binance("LTCUSDT", LTC)
			last4 := Binance("ETHUSDT", ETH)
			last5 := Binance("BCCBTC", BTCBCC)
			last6 := Binance("LTCBTC", BTCLTC)
			last7 := Binance("ETHBTC", BTCETH)
			if HasNull(last1, last2, last3, last4, last5, last6, last7) {
				bot.SendMessage(message.Chat, "查询失败，请重试", nil)
			} else {
				out := Output2(last1, last2, last3, last4, last5, last6, last7)
				msg := fmt.Sprintf("Binance: \n%s\n", out)
				bot.SendMessage(message.Chat, msg, nil)
			}
		} else {
			bot.SendMessage(message.Chat, "你等着，我等会找着了给你", nil)
		}
	}

}
