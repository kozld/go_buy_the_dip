package main

import (
	"buyTheDip/bot"
	"buyTheDip/config"
	"buyTheDip/store"
	"buyTheDip/strategy"
)

func main() {

	cfg := config.GetConfig()
	store := store.NewRedisStore()
	strategy := strategy.NewStrategy(store, &cfg)

	bot := bot.NewBinanceBot(&cfg, strategy)
	bot.Start()
}
