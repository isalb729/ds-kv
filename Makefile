.ONESHELL:
.PHONY: help build zk
SHELL=/bin/bash
PROTO=src/rpc/proto
help:

clean:

test:

build:

proto:
	@PROTO_LIST=$$(ls ${PROTO}| grep -E "*.proto")
	@for proto in $${PROTO_LIST}; do \
		protoc --proto_path=${PROTO} --go_out=plugins=grpc:${PROTO}/.. $${proto}; \
	done
zk:
	@cd scripts
	@./zk.sh

