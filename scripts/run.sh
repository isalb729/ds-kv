cd src/cmd/server
# master
go run server_main.go -cfg=../cfg.yaml -type=master -addr=:9666
go run server_main.go -cfg=../cfg.yaml -type=master -addr=:9667
# data
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server1
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server2
go run server_main.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server3
# standby data
go run server_main.go master.go slave.go -cfg=../cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=data/sb1

cd src/cmd/client
# client
go run client_main.go test.go -op=sequential -addr=:9666,127.0.0.1:9667



