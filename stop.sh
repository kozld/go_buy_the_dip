#!/bin/bash

docker-compose down
docker rm -f redis
docker rm -f bot