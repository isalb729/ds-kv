.ONESHELL:
.PHONY: help build zk
SHELL=/bin/bash
PROTO=src/rpc/proto
help:

clean:
	@rm -rf build
build:
	@go build -o build/server src/cmd/server/server_main.go
	@go build -o build/client src/cmd/client/client_main.go src/cmd/client/test.go
	@go build -o build/shell src/cmd/client-shell/shell.go
	@cp src/cmd/cfg.yaml build/
proto:
	@PROTO_LIST=$$(ls ${PROTO}| grep -E "*.proto")
	@for proto in $${PROTO_LIST}; do \
		protoc --proto_path=${PROTO} --go_out=plugins=grpc:${PROTO}/../pb $${proto}; \
	done
zk:
	@cd scripts
	@./zk.sh
zkcli:
	@zk/zk-1/bin/zkCli.sh -server=127.0.0.1:2181
