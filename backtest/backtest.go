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
	"buyTheDip/store"
	"buyTheDip/strategy"
)

var (
	time1 time.Time
)

type backTestBot struct {
	filename string
	strategy *strategy.FirstStrategy
	balance float64
	timeFrame float64
}

func NewBackTestBot(filename string, deposit float64, period int, takeProfit float64, timeFrame float64, timeout float64) bot.Bot {
	time1 = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	store := store.NewRedisStore()
	strategy := strategy.NewFirstStrategy(store, deposit, takeProfit, period, timeout)

	return &backTestBot{
		filename,
		strategy,
		strategy.Deposit,
		timeFrame,
	}
}

func (b *backTestBot) HandleCandle(time2 time.Time, price float64) error {

	//fmt.Println("HANDLE CANDLE")

	var buyValue float64
	time3 := time1.Add(time.Minute * time.Duration(int64(b.timeFrame)))
	if time2.After(time3) {
		buyValue = b.strategy.TryBuy(b.balance, time2, price)
		b.balance -= buyValue
		time1 = time2
	}

	sellValue := b.strategy.TrySell(time2, price)
	b.balance += sellValue

	return nil
}

func (b *backTestBot) Start() float64 {

	fmt.Printf("====================================\n")
	fmt.Printf("DEPOSIT     \t\t%f\n", b.strategy.Deposit)
	fmt.Printf("RSI PERIOD  \t\t%d\n", b.strategy.RSIPeriod)
	fmt.Printf("TAKE PROFIT \t\t%f\n", b.strategy.TakeProfit)
	fmt.Printf("TIME FRAME  \t\t%f\n", b.timeFrame)
	fmt.Printf("TIME OUT    \t\t%f\n", b.strategy.Timeout)
	fmt.Printf("====================================\n")

	file, err := os.Open(b.filename)
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
		tm := time.Unix(i / 1000, 0)

		price, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			panic(err)
		}

		b.HandleCandle(tm, price)
	}

	b.strategy.Store.Set([]store.Asset{})

	return b.strategy.Profit
}
