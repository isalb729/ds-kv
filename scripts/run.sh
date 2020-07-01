cd src/cmd/server
go run server_main.go master.go slave.go -cfg=../cfg.yaml -type=master -addr=:9666

go run server_main.go master.go slave.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server1
go run server_main.go master.go slave.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server2
go run server_main.go master.go slave.go -cfg=../cfg.yaml -type=slave -addr=127.0.0.1: -data=data/server3

go run server_main.go master.go slave.go -cfg=../cfg.yaml -type=slave-sb -addr=127.0.0.1: -data=sb/server1


go run client_main.go -op=sequential -addr=:9666,127.0.0.1:9667


# zk client
zk/zk-1/bin/zkCli.sh -server=127.0.0.1:2181
# zk server
