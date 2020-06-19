go run server_main.go -cfg=../cfg.yaml -type=master -addr=:9666
go run client_main.go -op=sequential -addr=:9666
#go run server_main.go -cfg=cfg.yaml -type=master -addr=127.0.0.1:
zk/zk-1/bin/zkCli.sh -server=127.0.0.1:2181

