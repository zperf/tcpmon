#!/usr/bin/env bash

mkdir -p data/tcpmon1/log
mkdir -p data/tcpmon2/log

docker run --name tcpmon1 -d --net=tcpmon -v "$(pwd)/bin:$(pwd)/bin" -w "$(pwd)/bin" \
-v "$(pwd)/data/tcpmon1/log:/tmp/tcpmon/log" tcpmon:runtime ./tcpmon-linux start

docker run --name tcpmon2 -d --net=tcpmon -v "$(pwd)/bin:$(pwd)/bin" -w "$(pwd)/bin" \
-v "$(pwd)/data/tcpmon2/log:/tmp/tcpmon/log" tcpmon:runtime ./tcpmon-linux start
