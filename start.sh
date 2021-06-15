#!/bin/bash

./stop.sh

source .env

docker build --no-cache -t go_binance_bot .

docker-compose -f redis/docker-compose.yaml up -d
sleep 15

docker run --name bot -e API_KEY=${API_KEY} -e API_SECRET=${API_SECRET} -e REDIS_PASSWORD=${REDIS_PASSWORD} --network=host -d go_binance_bot