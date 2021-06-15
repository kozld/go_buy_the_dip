package main

import (
	"fmt"

	"buyTheDip/backtest"
	"buyTheDip/config"
	"buyTheDip/store"
	"buyTheDip/strategy"
)

func main() {
	cfg := config.GetConfig()
	store := store.NewRedisStore()
	strategy := strategy.NewStrategy(store, &cfg)

	totalProfit := 0.0
	for day := 1; day <= 30; day++ {
		filename := fmt.Sprintf("backtest/fixtures/day_%d.csv", day)

		bot := backtest.NewBackTestBot(filename, &cfg, strategy)
		profit := bot.Start()

		totalProfit += profit
		fmt.Printf("Day %d\nDay profit: %f$\nTotal profit: %f$\n", day, profit, totalProfit)

		if profit < -30 {
			fmt.Println("Recalibrating RSI indicator...")
			rsiPeriod, timeFrame, holdTime := backtest.BruteForce(filename)

			cfg.RsiPeriod = rsiPeriod
			cfg.TimeFrame = float64(timeFrame)
			cfg.HoldTime = float64(holdTime)
		}
	}
}
