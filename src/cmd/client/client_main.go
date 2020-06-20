package main

import (
	"flag"
	"github.com/isalb729/ds-kv/src/client"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	// one master address
	// TODO: multiple master
	addr := flag.String("addr", "", "master or slave listening address")
	option := flag.String("op", "func1", "the client program to run")
	flag.Parse()
	if (*addr)[0] == ':' {
		*addr = "127.0.0.1" + *addr
	}
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	var try func(*client.KvCli)
	switch *option {
	case "concurrent":
		try = client.Concurrent
	case "sequential":
		try = client.Sequential
	default:
		log.Fatalln("not implemented option")
	}
	masterCli := pb.NewMetaClient(conn)
	try(&client.KvCli{Mc: masterCli})
	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
