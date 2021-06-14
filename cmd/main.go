package main

import (
	"flag"

	"buyTheDip/bot"
)

func main() {

	deposit := flag.Float64("deposit", 1000, "Deposit amount. (Default: 1000)")
	period := flag.Int("period", 5, "RSI period. (Default: 14)")
	takeProfit := flag.Float64("takeProfit", 0.3, "Take profit percentage. (Default: 2.5)")
	timeFrame := flag.Float64("timeFrame", 1, "Time frame in minutes. (Default: 1)")
	timeout := flag.Float64("timeout", 5, "Timeout in hours. (Default: 5)")
	flag.Parse()

	botOptions := &bot.Options{
		Deposit:    *deposit,
		Period:     *period,
		TakeProfit: *takeProfit,
		TimeFrame:  *timeFrame,
		Timeout:    *timeout,
	}

	bot := bot.NewBinanceBot(botOptions)
	bot.Start()
}
