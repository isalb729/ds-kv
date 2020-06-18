package main

import (
	"flag"
	"github.com/isalb729/ds-kv/src/client"
	"github.com/isalb729/ds-kv/src/rpc"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// one master address
	// TODO: multiple
	addr := flag.String("addr", "", "master or slave listening address")
	option := flag.String("op", "func1", "the client program to run")
	flag.Parse()
	if (*addr)[0] == ':' {
		*addr = "127.0.0.1" + *addr
	}
	conn, err := grpc.Dial(*addr)
	if err != nil {
		log.Fatalln(err)
	}
	var try func(rpc.MasterClient)
	switch *option {
	case "concurrent":
		try = client.Concurrent
	case "sequential":
		try = client.Sequential
	default:
		log.Fatalln("not implemented option")
	}
	cli := rpc.NewMasterClient(conn)
	try(cli)
	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
