package main

import (
	"flag"
	"github.com/isalb729/ds-kv/src/client"
	"log"
	"strings"
)

func main() {
	// one master address
	addr := flag.String("addr", "", "master or slave listening address")
	option := flag.String("op", "func1", "the client program to run")
	flag.Parse()
	addrList := strings.Split(*addr, ",")
	for _, addr := range addrList {
		if (addr)[0] == ':' {
			addr = "127.0.0.1" + addr
		}
	}
	cli := client.Connect(addrList)
	var try func(*client.KvCli)
	switch *option {
	case "concurrent":
		try = client.Concurrent
	case "sequential":
		try = client.Sequential
	default:
		log.Fatalln("not implemented option")
	}
	try(cli)
}
