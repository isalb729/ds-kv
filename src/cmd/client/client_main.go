package main

import (
	"flag"
	"github.com/isalb729/ds-kv/src/client"
	"log"
	"strings"
)

func main() {
	addr := flag.String("addr", "", "master or slave listening address")
	option := flag.String("op", "func1", "the client program to run")
	flag.Parse()
	// Master address list.
	addrList := strings.Split(*addr, ",")
	for _, addr := range addrList {
		if (addr)[0] == ':' {
			// Default host.
			addr = "127.0.0.1" + addr
		}
	}
	cli := client.Connect(addrList)
	var try func(*client.KvCli)
	// Some test cases.
	switch *option {
	case "concurrent":
		try = Concurrent
	case "sequential":
		try = Sequential
	case "crazyloop":
		try = CrazyLoop
	default:
		log.Fatalln("not implemented option")
	}
	try(cli)
}
