#!/usr/bin/env bash

set -x

IMAGE_DISTRO="$1"

docker ps -a | grep 'tcpmon[0-9]' | awk '{print $1}' | parallel docker stop
docker ps -a | grep 'tcpmon[0-9]' | awk '{print $1}' | parallel docker rm

mkdir -p data/tcpmon1/log
docker run --name tcpmon1 -d --net=tcpmon -v "$(pwd)/bin:$(pwd)/bin" -w "$(pwd)/bin" \
-v "$(pwd)/data/tcpmon1/log:/tmp/tcpmon/log" \
"tcpmon:runtime$IMAGE_DISTRO" ./tcpmon-linux start

mkdir -p data/tcpmon2/log
docker run --name tcpmon2 -d --net=tcpmon -v "$(pwd)/bin:$(pwd)/bin" -w "$(pwd)/bin" \
-v "$(pwd)/data/tcpmon2/log:/tmp/tcpmon/log" \
"tcpmon:runtime$IMAGE_DISTRO" ./tcpmon-linux start

mkdir -p data/tcpmon3/log
docker run --name tcpmon3 -d --net=tcpmon -v "$(pwd)/bin:$(pwd)/bin" -w "$(pwd)/bin" \
-v "$(pwd)/data/tcpmon3/log:/tmp/tcpmon/log" \
"tcpmon:runtime$IMAGE_DISTRO" ./tcpmon-linux start
