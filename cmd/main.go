package main

import (
	"buyTheDip/bot"
	"flag"
)

func main() {

	// filename := flag.String("filename", "", "Path to CSV file with historical data. (Required)")
	deposit := flag.Float64("deposit", 1000, "Deposit amount. (Default: 1000)")
	period := flag.Int("period", 5, "RSI period. (Default: 14)")
	takeProfit := flag.Float64("takeProfit", 0.3, "Take profit percentage. (Default: 2.5)")
	timeFrame := flag.Float64("timeFrame", 1, "Time frame in minutes. (Default: 1)")
	timeout := flag.Float64("timeout", 5, "Timeout in hours. (Default: 5)")
	flag.Parse()

	//if *filename == "" {
	//	flag.PrintDefaults()
	//	os.Exit(1)
	//}

	bot := bot.NewBinanceBot(*deposit, *period, *takeProfit, *timeFrame, *timeout)
	bot.Start()

	//deposit := 1000.0
	//takeProfit := 0.3
	//
	//rsi := 5
	//timeFrame := 1
	//timeout := 5
	//
	//totalProfit := 0.0
	//for day := 1; day <= 30; day++ {
	//
	//	filename := fmt.Sprintf("backtest/fixtures/day_%d.csv", day)
	//
	//	bot := backtest.NewBackTestBot(filename, deposit, rsi, takeProfit, float64(timeFrame), float64(timeout))
	//	profit := bot.Start()
	//	totalProfit += profit
	//
	//	fmt.Printf("Day %d\nDay profit: %f$\nTotal profit: %f$\n", day, profit, totalProfit)

	//fmt.Println("Recalibration...")
	//rsi, timeFrame, timeout = utils.BruteForce(filename, deposit, takeProfit)

	//if profit < -30 {
	//	fmt.Println("Recalibration...")
	//	rsi, timeFrame, timeout = utils.BruteForce(filename, deposit, takeProfit)
	//}
	//}

}
