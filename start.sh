#!/bin/bash

./stop.sh

source .env

docker build --no-cache -t go_binance_bot .

docker-compose -f redis/docker-compose.yaml up -d
sleep 15

docker run --name bot -e BINANCE_API=${BINANCE_API} -e BINANCE_SECRET=${BINANCE_SECRET} -e REDIS_PASSWORD=${REDIS_PASSWORD} -e TICKER=${TICKER} -e TELEGRAM_BOT_TOKEN=${TELEGRAM_BOT_TOKEN} -e TELEGRAM_CHANNEL_NAME=${TELEGRAM_CHANNEL_NAME} --network=host -d go_binance_bot