.PHONY: all
all: builder

.PHONY: builder
builder:
	docker build --target builder -f Dockerfile -t tcpmon:builder .

.PHONY: runtime
runtime:
	docker build --target runtime -f Dockerfile -t tcpmon:runtime .
