package bot

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"

	"buyTheDip/store"
	"buyTheDip/strategy"
)

var (
	time1 time.Time

	tradePair       = "DOGEUSDT"
	wsKlineInterval = "1m"

	apiKey    = os.Getenv("API_KEY")
	secretKey = os.Getenv("API_SECRET")
)

type binanceBot struct {
	client    *binance.Client
	strategy  *strategy.FirstStrategy
	balance   float64
	timeFrame float64
}

func NewBinanceBot(deposit float64, period int, takeProfit float64, timeFrame float64, timeout float64) Bot {
	time1 = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	client := binance.NewClient(apiKey, secretKey)
	store := store.NewRedisStore()
	strategy := strategy.NewFirstStrategy(store, deposit, takeProfit, period, timeout)

	return &binanceBot{
		client,
		strategy,
		strategy.Deposit,
		timeFrame,
	}
}

func (b *binanceBot) CreateBuyOrder(qty string) error {
	_, err := b.client.NewCreateOrderService().Symbol(tradePair).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(qty).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (b *binanceBot) CreateSellOrder(qty string) error {
	_, err := b.client.NewCreateOrderService().Symbol(tradePair).
		Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).
		TimeInForce(binance.TimeInForceTypeGTC).Quantity(qty).Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func (b *binanceBot) HandleCandle(time2 time.Time, price float64) error {

	var buyValue float64
	time3 := time1.Add(time.Minute * time.Duration(int64(b.timeFrame)))
	if time2.After(time3) {
		buyValue = b.strategy.TryBuy(b.balance, time2, price, b.CreateBuyOrder)
		b.balance -= buyValue
		time1 = time2
	}

	sellValue := b.strategy.TrySell(time2, price, b.CreateSellOrder)
	b.balance += sellValue

	return nil
}

func (b *binanceBot) Start() float64 {

	fmt.Printf("====================================\n")
	fmt.Printf("DEPOSIT     \t\t%f\n", b.strategy.Deposit)
	fmt.Printf("RSI PERIOD  \t\t%d\n", b.strategy.RSIPeriod)
	fmt.Printf("TAKE PROFIT \t\t%f\n", b.strategy.TakeProfit)
	fmt.Printf("TIME FRAME  \t\t%f\n", b.timeFrame)
	fmt.Printf("TIME OUT    \t\t%f\n", b.strategy.Timeout)
	fmt.Printf("====================================\n")

	wsKlineHandler := func(event *binance.WsKlineEvent) {
		time := time.Unix(event.Time/1000, 0)
		price, err := strconv.ParseFloat(event.Kline.Close, 64)
		if err != nil {
			panic(err)
		}

		_ = b.HandleCandle(time, price)
	}

	errHandler := func(err error) {
		fmt.Println(err)
	}

	doneC, _, err := binance.WsKlineServe(tradePair, wsKlineInterval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	<-doneC

	return b.strategy.Profit
}
