#!/usr/bin/env bash

set -x

make build-linux

IMAGE_DISTRO="$1"

docker ps -a | grep 'tcpmon[0-9]' | awk '{print $1}' | parallel docker stop
docker ps -a | grep 'tcpmon[0-9]' | awk '{print $1}' | parallel docker rm

for IDX in $(seq 1 3); do
  mkdir -p "data/tcpmon$IDX/log"
  mkdir -p "data/tcpmon$IDX/db"
  docker run --name "tcpmon$IDX" -d --net=tcpmon \
  -p "$(echo "$IDX" | awk '{print 6788+$1}'):6789" \
  -p "$(echo "$IDX" | awk '{print 6879+$1}'):6790" \
  -p "$(echo "$IDX" | awk '{print 2344+$1}'):2345" \
  -w "$(pwd)/bin" \
  -v "$(pwd)/bin:$(pwd)/bin" \
  -v "$(pwd)/data/tcpmon$IDX/log:/tmp/tcpmon/log" \
  -v "$(pwd)/data/tcpmon$IDX/db:/tmp/tcpmon/db" \
  "tcpmon:runtime$IMAGE_DISTRO" \
  dlv exec --listen=:2345 --headless=true --api-version=2 --accept-multiclient --continue ./tcpmon start
done

sleep 1

docker exec -t tcpmon1 curl -X POST -d 192.168.228.3:6790 http://192.168.228.2:6789/members
docker exec -t tcpmon1 curl -X POST -d 192.168.228.4:6790 http://192.168.228.2:6789/members
