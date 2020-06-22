package main

import (
	"flag"
	"fmt"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/isalb729/ds-kv/src/zookeeper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Cfg struct {
	Zk     []string `yaml:"zk"`
}

func main() {
	// server address and type are passed through command line
	// default random port
	addr := flag.String("addr", "127.0.0.1:", "Master or slave listening address")
	tp := flag.String("type", "slave", "Master or slave")
	dataDir := flag.String("data", "", "Data directory; not end with /")
	// zk info is stored in config file
	var cfg Cfg
	err := utils.LoadConfig(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	if len(*addr) == 0 {
		log.Fatalln("Empty address")
	}
	if (*addr)[0] == ':' {
		*addr = "127.0.0.1" + *addr
	}
	if *dataDir != "" && (*dataDir)[len(*dataDir) - 1] == '/' {
		*dataDir = (*dataDir)[:len(*dataDir) - 1]
	}
	zkConn, err := zookeeper.Connect(cfg.Zk)
	if err != nil {
		log.Fatalln(err)
	}
	defer zkConn.Close()
	splitAddr := strings.Split(*addr, ":")
	host := splitAddr[0]
	lis, err := net.Listen("tcp", ":" + splitAddr[1])
	if err != nil {
		log.Fatalln(err)
	}
	splitAddr = strings.Split(lis.Addr().String(), ":")
	//fmt.Println(splitAddr)
	port := splitAddr[len(splitAddr) - 1]
	//fmt.Println(*addr)

	name := host + ":" + port
	log.Printf("Server name: %s\n", name)
	grpcServer := grpc.NewServer()
	switch *tp {
	case "master":
		err = InitMaster(grpcServer, zkConn, name)
		if err != nil {
			log.Fatalln("Fail to init master", err)
		}
	case "slave":
		err = InitSlave(grpcServer, zkConn, name, *dataDir)
		if err != nil {
			log.Fatalln("Fail to init slave", err)
		}
	default:
		log.Fatalln("Wrong type: only master or slave is supported")
	}
	log.Printf("Server run on addr: %s\n", lis.Addr())

	// run as a thread
	errChan := make(chan error)
	go func() {
		if err :=  grpcServer.Serve(lis); err != nil {
			errChan <- err
			log.Fatal("Fail to run the service", err)
			return
		}
	}()

	// listen for signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	log.Println("Shutting down the service...", <-errChan)

	/* TODO re/degister lock*/
	// stop the panicking
	r := recover()
	if r != nil {
		log.Println("Recovering from", r)
	}
	if *tp == "slave" {
		err = deregisterSlave(zkConn, *dataDir, name)
	} else if *tp == "master" {
		err = deregisterMaster(zkConn)
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Deregistered %s %s\n", *tp, name)
	}
	// TODO: unlock register
}
