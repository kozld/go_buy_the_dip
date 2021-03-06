package backtest

import (
	"buyTheDip/config"
	"buyTheDip/store"
	"buyTheDip/strategy"
	"fmt"
)

func BruteForce(filename string) (int, int, int) {
	rsiStart := 2
	rsiEnd := 10
	frameStart := 1
	frameEnd := 5
	timeoutStart := 4
	timeoutEnd := 24
	timeoutStep := 4

	bestProfit := -9999.0
	bestRsi := rsiStart
	bestFrame := frameStart
	bestTimeout := timeoutStart

	cfg := config.GetConfig()
	store := store.NewRedisStore()
	strategy := strategy.NewStrategy(store, &cfg)

	for rsi := rsiStart; rsi <= rsiEnd; rsi++ {
		for frame := frameStart; frame <= frameEnd; frame++ {
			for timeout := timeoutStart; timeout <= timeoutEnd; timeout += timeoutStep {

				bot := NewBackTestBot(filename, &cfg, strategy)
				profit := bot.Start()

				fmt.Printf("[PROFIT] %f\n", profit)
				fmt.Printf("[RSI] %d\n", rsi)
				fmt.Printf("[FRAME] %d\n", frame)
				fmt.Printf("[TIMEOUT] %d\n", timeout)

				if profit > bestProfit {
					bestProfit = profit
					bestRsi = rsi
					bestFrame = frame
					bestTimeout = timeout
				}
			}
		}
	}

	fmt.Printf("====================================\n")
	fmt.Printf("BEST PROFIT     \t\t%f\n", bestProfit)
	fmt.Printf("BEST RSI  	   \t\t%d\n", bestRsi)
	fmt.Printf("BEST FRAME  	   \t\t%d\n", bestFrame)
	fmt.Printf("BEST TIMEOUT    \t\t%d\n", bestTimeout)
	fmt.Printf("====================================\n")

	return bestRsi, bestFrame, bestTimeout
}
