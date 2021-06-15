package backtest

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"buyTheDip/bot"
	"buyTheDip/config"
	"buyTheDip/store"
	"buyTheDip/strategy"
)

var (
	time1 time.Time
)

type BackTestBot struct {
	Filename string
	Strategy *strategy.Strategy
	Config   *config.BotConfig
	Balance  float64
}

func NewBackTestBot(filename string, cfg *config.BotConfig, strategy *strategy.Strategy) bot.Bot {
	time1 = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	return &BackTestBot{
		Filename: filename,
		Strategy: strategy,
		Config:   cfg,
		Balance:  cfg.Deposit,
	}
}

func (b *BackTestBot) CreateBuyOrder(price float64, qty string) error {
	return nil
}

func (b *BackTestBot) CreateSellOrder(price float64, qty string) error {
	return nil
}

func (b *BackTestBot) HandleCandle(time2 time.Time, price float64) error {

	var buyValue float64
	time3 := time1.Add(time.Minute * time.Duration(int64(b.Config.TimeFrame)))

	if time2.After(time3) {
		buyValue = b.Strategy.TryBuy(b.Balance, time2, price, b.CreateBuyOrder)
		b.Balance -= buyValue
		time1 = time2
	}

	sellValue := b.Strategy.TrySell(time2, price, b.CreateSellOrder)
	b.Balance += sellValue

	return nil
}

func (b *BackTestBot) Start() float64 {

	fmt.Printf("====================================\n")
	fmt.Printf("DEPOSIT     \t\t%f\n", b.Config.Deposit)
	fmt.Printf("RSI PERIOD  \t\t%d\n", b.Config.RsiPeriod)
	fmt.Printf("TAKE PROFIT \t\t%f\n", b.Config.TakeProfit)
	fmt.Printf("TIME FRAME  \t\t%f\n", b.Config.TimeFrame)
	fmt.Printf("HOLD TIME   \t\t%f\n", b.Config.HoldTime)
	fmt.Printf("====================================\n")

	file, err := os.Open(b.Filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		i, err := strconv.ParseInt(record[0], 10, 64)
		if err != nil {
			panic(err)
		}
		tm := time.Unix(i/1000, 0)

		price, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			panic(err)
		}

		b.HandleCandle(tm, price)
	}

	b.Strategy.Store.Set([]store.Asset{})

	return 0 // b.Strategy.Profit
}
