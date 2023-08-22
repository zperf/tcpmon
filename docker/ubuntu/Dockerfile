FROM ubuntu:jammy AS builder
RUN apt update && apt install -y protobuf-compiler protoc-gen-go

FROM ubuntu:jammy AS runtime
RUN sed -i 's@//.*archive.ubuntu.com@//mirrors.ustc.edu.cn@g' /etc/apt/sources.list && \
apt update -y && apt install -y iproute2 net-tools curl which jq
