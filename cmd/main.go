package main

import (
	"buyTheDip/backtest"
	"flag"
	"fmt"
	"os"
)

func main() {

	filename := flag.String("filename", "", "Path to CSV file with historical data. (Required)")
	deposit := flag.Float64("deposit", 1000, "Deposit amount. (Default: 1000)")
	// period := flag.Int("period", 5, "RSI period. (Default: 14)")
	takeProfit := flag.Float64("takeProfit", 0.3, "Take profit percentage. (Default: 2.5)")
	// timeFrame := flag.Float64("timeFrame", 1, "Time frame in minutes. (Default: 1)")
	// timeout := flag.Float64("timeout", 5, "Timeout in hours. (Default: 5)")
	flag.Parse()

	if *filename == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	rsiStart := 5
	rsiEnd := 5
	frameStart := 1
	frameEnd := 1
	timeoutStart := 5
	timeoutEnd := 5
	timeoutStep := 1

	bestProfit := -9999.0
	bestRsi := rsiStart
	bestFrame := frameStart
	bestTimeout := timeoutStart

	for rsi := rsiStart; rsi <= rsiEnd; rsi++ {
		for frame := frameStart; frame <= frameEnd; frame++ {
			for timeout := timeoutStart; timeout <= timeoutEnd; timeout += timeoutStep {

				bot := backtest.NewBackTestBot(*filename, *deposit, rsi, *takeProfit, float64(frame), float64(timeout))
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

				//time.Sleep(1 * time.Second)
			}
		}
	}

	fmt.Printf("====================================\n")
	fmt.Printf("BEST PROFIT     \t\t%f\n", bestProfit)
	fmt.Printf("BEST RSI  	   \t\t%d\n", bestRsi)
	fmt.Printf("BEST FRAME  	   \t\t%d\n", bestFrame)
	fmt.Printf("BEST TIMEOUT    \t\t%d\n", bestTimeout)
	fmt.Printf("====================================\n")
}
