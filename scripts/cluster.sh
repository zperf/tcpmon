#!/usr/bin/env bash

set -x

IMAGE_DISTRO="$1"

docker ps -a | grep 'tcpmon[0-9]' | awk '{print $1}' | parallel docker stop
docker ps -a | grep 'tcpmon[0-9]' | awk '{print $1}' | parallel docker rm

mkdir -p data/tcpmon1/log
docker run --name tcpmon1 -d --net=tcpmon -p 6789:6789 -v "$(pwd)/bin:$(pwd)/bin" -w "$(pwd)/bin" \
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

docker exec -t tcpmon1 bash -c "curl -X POST -d 192.168.228.3:6790 http://192.168.228.2:6789/members"
docker exec -t tcpmon1 bash -c "curl -X POST -d 192.168.228.4:6790 http://192.168.228.2:6789/members"
