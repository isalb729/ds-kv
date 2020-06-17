.ONESHELL:
.PHONY: help build zk
SHELL=/bin/bash
help:

clean:

test:

build:

proto:
	@protoc --proto_path=src/rpc/proto --go_out=plugins=grpc:src/rpc server.proto;
zk:
	@cd scripts
	@./zk.sh