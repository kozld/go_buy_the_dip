package strategy

import (
	"fmt"
	"strconv"
	"time"

	"github.com/markcheno/go-talib"

	"buyTheDip/config"
	"buyTheDip/store"
)

var (
	prices  []float64
	prevRsi float64
)

type Strategy struct {
	Store  store.Store
	Config *config.BotConfig
}

func NewStrategy(store store.Store, cfg *config.BotConfig) *Strategy {
	return &Strategy{
		Store:  store,
		Config: cfg,
	}
}

func (s *Strategy) isSellByTakeProfit(buyPrice float64, sellPrice float64) bool {
	return (sellPrice-buyPrice)/buyPrice*100 >= s.Config.TakeProfit
}

func (s *Strategy) isSellByTime(openTime time.Time, currentTime time.Time) bool {
	return currentTime.Sub(openTime).Hours() >= s.Config.HoldTime
}

func (s *Strategy) TryBuy(balance float64, time time.Time, price float64, callback func(string) error) float64 {

	var rsi []float64

	assets, err := s.Store.Get()
	if err != nil {
		panic(err)
	}

	prices = append(prices, price)
	if len(prices) > s.Config.RsiPeriod {
		rsi = talib.Rsi(prices, s.Config.RsiPeriod)
		prices = prices[len(prices)-s.Config.RsiPeriod:]
	}

	totalPrice := .0

	if len(rsi) > 0 {
		fmt.Println("RSI -", rsi[len(rsi)-1])

		if prevRsi >= s.Config.RsiOversold && rsi[len(rsi)-1] < s.Config.RsiOversold {

			fmt.Println("!!! OVERSOLD !!!")
			fmt.Println("RSI -", rsi[len(rsi)-1])
			fmt.Println("PRICE -", price)

			qty := fmt.Sprintf("%.1f", s.Config.Deposit/price/4)

			_qty, err := strconv.ParseFloat(qty, 64)
			if err != nil {
				panic(err)
			}

			if balance >= _qty*price {
				fmt.Printf("BUY %f x %f$\n", _qty, price)

				if err := callback(qty); err != nil {
					fmt.Printf("error: %v", err)
					return 0
				}

				asset := store.Asset{
					Price: price,
					Time:  time,
					Qty:   _qty,
				}
				assets = append(assets, asset)
				totalPrice = _qty * price

				if err := s.Store.Set(assets); err != nil {
					panic(err)
				}
			}
		}

		prevRsi = rsi[len(rsi)-1]
	}

	return totalPrice
}

func (s *Strategy) TrySell(time time.Time, price float64, callback func(string) error) float64 {

	var _assets []store.Asset
	assets, err := s.Store.Get()

	if err != nil {
		panic(err)
	}

	totalPrice := .0
	for _, asset := range assets {
		purchasePrice := asset.Price * asset.Qty
		sellingPrice := price * asset.Qty

		if s.isSellByTime(asset.Time, time) {
			fmt.Printf("SELL BY TIME %f x %f$\n", asset.Qty, price)
			fmt.Printf("PROFIT - %f\n", sellingPrice-purchasePrice)

			qty := fmt.Sprintf("%.1f", asset.Qty)
			if err := callback(qty); err != nil {
				fmt.Printf("error: %v", err)
				return 0
			}

			// s.Profit += sellingPrice - purchasePrice
			totalPrice += sellingPrice
		} else if s.isSellByTakeProfit(asset.Price, price) {
			fmt.Printf("SELL BY TAKE PROFIT %f x %f$\n", asset.Qty, price)
			fmt.Printf("PROFIT - %f\n", sellingPrice-purchasePrice)

			qty := fmt.Sprintf("%.1f", asset.Qty)
			if err := callback(qty); err != nil {
				fmt.Printf("error: %v", err)
				return 0
			}

			// s.Profit += sellingPrice - purchasePrice
			totalPrice += sellingPrice
		} else {
			_assets = append(_assets, asset)
		}
	}

	if err := s.Store.Set(_assets); err != nil {
		panic(err)
	}

	return totalPrice
}
