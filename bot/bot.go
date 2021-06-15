package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"

	"buyTheDip/config"
	"buyTheDip/strategy"
)

var (
	time1           time.Time
	wsKlineInterval = "1m"
	mainAsset       = "USDT"
)

type BinanceBot struct {
	Client   *binance.Client
	Strategy *strategy.Strategy
	Config   *config.BotConfig
}

func NewBinanceBot(cfg *config.BotConfig, strategy *strategy.Strategy) Bot {

	time1 = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	fmt.Println(cfg.TelegramBotToken)
	fmt.Println(cfg.TelegramChannelName)

	return &BinanceBot{
		Client:   binance.NewClient(cfg.BinanceApi, cfg.BinanceSecret),
		Strategy: strategy,
		Config:   cfg,
	}
}

func (b *BinanceBot) CreateBuyOrder(price float64, qty string) error {

	_, err := b.Client.NewCreateOrderService().Symbol(b.Config.Ticker).
		Side(binance.SideTypeBuy).Type(binance.OrderTypeMarket).
		Quantity(qty).Do(context.Background())
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	msg := fmt.Sprintf("BUY #%s: $%.2f x %s", b.Config.Ticker, price, qty)
	err = b.SendMessage(msg)
	if err != nil {
		log.Printf("error: %v", err)
	}

	return err
}

func (b *BinanceBot) CreateSellOrder(price float64, qty string) error {

	_, err := b.Client.NewCreateOrderService().Symbol(b.Config.Ticker).
		Side(binance.SideTypeSell).Type(binance.OrderTypeMarket).
		Quantity(qty).Do(context.Background())
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}

	msg := fmt.Sprintf("SELL #%s: $%.2f x %s", b.Config.Ticker, price, qty)
	err = b.SendMessage(msg)
	if err != nil {
		log.Printf("error: %v", err)
	}

	return err
}

func (b *BinanceBot) GetBalance(name string) (float64, error) {

	res, err := b.Client.NewGetAccountService().Do(context.Background())
	if err != nil {
		fmt.Printf("error: %v", err)
		return 0, err
	}

	for _, asset := range res.Balances {
		if asset.Asset == name {
			value, err := strconv.ParseFloat(asset.Free, 64)
			if err != nil {
				return 0, err
			}

			return value, nil
		}
	}

	return 0, err
}

func (b *BinanceBot) HandleCandle(time2 time.Time, price float64) error {

	time3 := time1.Add(time.Minute * time.Duration(int64(b.Config.TimeFrame)))
	if time2.After(time3) {

		balance, err := b.GetBalance(mainAsset)
		if err != nil {
			fmt.Printf("error: %v", err)
			return err
		}
		_ = b.Strategy.TryBuy(balance, time2, price, b.CreateBuyOrder)

		time1 = time2
	}

	_ = b.Strategy.TrySell(time2, price, b.CreateSellOrder)

	return nil
}

func (b *BinanceBot) Start() float64 {

	fmt.Printf("====================================\n")
	fmt.Printf("DEPOSIT     \t\t%f\n", b.Config.Deposit)
	fmt.Printf("RSI PERIOD  \t\t%d\n", b.Config.RsiPeriod)
	fmt.Printf("TAKE PROFIT \t\t%f\n", b.Config.TakeProfit)
	fmt.Printf("TIME FRAME  \t\t%f\n", b.Config.TimeFrame)
	fmt.Printf("HOLD TIME   \t\t%f\n", b.Config.HoldTime)
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

	doneC, _, err := binance.WsKlineServe(b.Config.Ticker, wsKlineInterval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	<-doneC

	return 0 // b.Strategy.Profit
}
