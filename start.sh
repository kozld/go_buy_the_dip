#!/bin/bash

./stop.sh

docker-compose build --no-cache
docker-compose up -d redis
sleep 15
docker-compose up -d bot