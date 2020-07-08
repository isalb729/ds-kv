cd src/cmd/server
# master
go run server_main.go -cfg=../cfg.yaml -type=master -addr=:9666
go run server_main.go -cfg=../cfg.yaml -type=master -addr=:9667
cd build
./server -cfg=cfg.yaml -type=master -addr=:9666
# data
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server1
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server2
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server3
cd build
./server -cfg=cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server1
# standby data
go run server_main.go -cfg=../cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=data/sb1
cd build
./server -cfg=cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=data/sb1

# client
cd src/cmd/client
# go run client_main.go test.go -op=sequential -addr=:9666,127.0.0.1:9667
cd build
./client -op=sequential -addr=:9666,127.0.0.1:9667

# shell client
cd src/cmd/client-shell
# go run shell.go -addr=:9666,127.0.0.1:9667
cd build
./shell -addr=:9666,127.0.0.1:9667
