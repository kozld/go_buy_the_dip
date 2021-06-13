package strategy

import (
	"fmt"
	"strconv"
	"time"

	"github.com/markcheno/go-talib"

	"buyTheDip/store"
)

var (
	prices []float64
	prevRsi float64
)

type FirstStrategy struct {
	Store store.Store
	Deposit float64
	TakeProfit float64
	RSIPeriod int
	Timeout float64
	Profit float64
}

func NewFirstStrategy(store store.Store, deposit float64, takeProfit float64, RSIPeriod int, timeout float64) *FirstStrategy {
	return &FirstStrategy{
		Store: store,
		Deposit: deposit,
		TakeProfit: takeProfit,
		RSIPeriod: RSIPeriod,
		Timeout: timeout,
	}
}

func (s *FirstStrategy) isSellByTakeProfit(buyPrice float64, sellPrice float64) bool {
	return (sellPrice - buyPrice) / buyPrice * 100 >= s.TakeProfit
}

func (s *FirstStrategy) isSellByTime(openTime time.Time, currentTime time.Time) bool {
	return currentTime.Sub(openTime).Hours() >= s.Timeout
}

func (s *FirstStrategy) TryBuy(balance float64, time time.Time, price float64) float64 {

	// fmt.Println("Try BUY")

	var rsi []float64

	assets, err := s.Store.Get()
	if err != nil {
		panic(err)
	}

	prices = append(prices, price)
	if len(prices) > s.RSIPeriod {
		rsi = talib.Rsi(prices, s.RSIPeriod)
		prices = prices[len(prices)-s.RSIPeriod:]
	}

	totalPrice := .0

	if len(rsi) > 0 {

		// fmt.Println("RSI:", rsi[len(rsi)-1])

		if prevRsi >= 30 && rsi[len(rsi)-1] < 30 {

			//fmt.Println("OVERSOLD !!!", rsi[len(rsi)-1])

			qty := fmt.Sprintf("%.f", s.Deposit / price / 4)

			_qty, err := strconv.ParseFloat(qty, 64)
			if err != nil {
				panic(err)
			}

			if balance >= _qty * price {
				//fmt.Printf("BUY %f x %f$\n", _qty, price)

				asset := store.Asset{
					Price: price,
					Time: time,
					Qty: _qty,
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

func (s *FirstStrategy) TrySell(time time.Time, price float64) float64 {

	// fmt.Println("Try SELL")

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
			//fmt.Printf("SELL BY TIME %f x %f$\n", asset.Qty, price)
			//fmt.Printf("PROFIT: %f\n", sellingPrice - purchasePrice)

			s.Profit += sellingPrice - purchasePrice
			totalPrice += sellingPrice
		} else if s.isSellByTakeProfit(asset.Price, price) {
			//fmt.Printf("SELL BY TAKE PROFIT %f x %f$\n", asset.Qty, price)
			//fmt.Printf("PROFIT: %f\n", sellingPrice - purchasePrice)

			s.Profit += sellingPrice - purchasePrice
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