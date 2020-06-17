package main

import (
	"flag"
	"fmt"
	"github.com/isalb729/ds-kv/src/client"
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/utils"
	"google.golang.org/grpc"
	"log"
)



func main() {
	option := flag.String("op", "func1", "the client program to run")
	var cfg Cfg
	err := utils.LoadConfig(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	var a func(rpc.AuthClient)
	switch *option {
	case "concurrent":
		a := client.F1
	case "sequence":
	default:
	}
	masterAddr := fmt.Sprintf("%s:%d", cfg.Master.Host, cfg.Master.Port)
	conn, err := grpc.Dial(masterAddr)
	if err != nil {
		log.Fatalln(err)
	}
	cli := rpc.NewAuthClient(conn)
	a(cli)
	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
