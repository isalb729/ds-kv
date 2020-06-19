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
	"syscall"
)

type Cfg struct {
	Zk     []string `yaml:"zk"`
}

func main() {
	// server address and type are passed through command line
	// default random port
	addr := flag.String("addr", "127.0.0.1:", "master or slave listening address")
	tp := flag.String("type", "slave", "master or slave")
	dataDir := flag.String("data", "", "data directory")
	if len(*addr) == 0 {
		log.Fatalln("empty address")
	}
	if (*addr)[0] == ':' {
		*addr = "127.0.0.1" + *addr
	}

	// zk info is stored in config file
	var cfg Cfg
	err := utils.LoadConfig(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	zkConn, err := zookeeper.Connect(cfg.Zk)
	if err != nil {
		log.Fatalln(err)
	}
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	switch *tp {
	case "master":
		err = InitMaster(grpcServer, zkConn, lis.Addr().String())
		if err != nil {
			log.Fatalln("fail to init master", err)
		}
	case "slave":
		err = InitSlave(grpcServer, zkConn, lis.Addr().String(), *dataDir)
		if err != nil {
			log.Fatalln("fail to init slave", err)
		}
	default:
		log.Fatalln("wrong type: only master or slave is supported")
	}
	log.Printf("server run on addr: %s\n", lis.Addr())

	// run as a thread
	errChan := make(chan error)
	go func() {
		if err :=  grpcServer.Serve(lis); err != nil {
			errChan <- err
			log.Fatal("fail to run the service", err)
			return
		}
	}()

	// listen for signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	log.Println("shutting down the service...", <-errChan)
	/* TODO: Deregister. */
	/* acquire re/degister lock*/
	/* delete node */
	/* delete data directory*/
	zkConn.Close()
}

