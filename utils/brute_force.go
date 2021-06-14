package utils

import (
	"fmt"

	"buyTheDip/backtest"
)

func BruteForce(filename string, deposit float64, takeProfit float64) (int, int, int) {
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

	for rsi := rsiStart; rsi <= rsiEnd; rsi++ {
		for frame := frameStart; frame <= frameEnd; frame++ {
			for timeout := timeoutStart; timeout <= timeoutEnd; timeout += timeoutStep {

				bot := backtest.NewBackTestBot(filename, deposit, rsi, takeProfit, float64(frame), float64(timeout))
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
