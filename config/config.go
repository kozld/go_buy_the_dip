package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type BotConfig struct {

	// ticker example: "DOGEUSDT"
	Ticker string `env:"TICKER,required"`
	// amount of the deposit in the main currency
	Deposit float64 `env:"DEPOSIT" envDefault:"1000"`
	// take profit in percentage
	TakeProfit float64 `env:"TAKE_PROFIT" envDefault:"2.5"`
	// time frame in minutes
	TimeFrame float64 `env:"TIME_FRAME" envDefault:"1"`
	// hold time in hours
	HoldTime float64 `env:"HOLD_TIME" envDefault:"5"`

	// RSI indicator
	RsiPeriod     int     `env:"RSI_PERIOD" envDefault:"5"`
	RsiOverbought float64 `env:"RSI_OVERBOUGHT" envDefault:"70"`
	RsiOversold   float64 `env:"RSI_OVERSOLD" envDefault:"30"`

	// telegram notifications
	TelegramBotToken    string `env:"TELEGRAM_BOT_TOKEN,required"`
	TelegramChannelName string `env:"TELEGRAM_CHANNEL_NAME,required"`

	// binance
	BinanceApi    string `env:"BINANCE_API,required"`
	BinanceSecret string `env:"BINANCE_SECRET,required"`
}

func GetConfig() BotConfig {
	cfg := &BotConfig{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal("Cannot parse initial ENV vars: ", err)
	}
	return *cfg
}
