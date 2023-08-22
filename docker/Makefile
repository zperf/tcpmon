.PHONY: all
all: builder runtime

.PHONY: builder
builder:
	$(MAKE) -C ubuntu builder

.PHONY: runtime
runtime:
	$(MAKE) -C el7 runtime
	$(MAKE) -C oe1 runtime
	$(MAKE) -C ubuntu runtime
