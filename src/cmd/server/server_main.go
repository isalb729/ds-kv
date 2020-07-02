package main

import (
	"flag"
	"fmt"
	"github.com/isalb729/ds-kv/src/server"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/isalb729/ds-kv/src/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Cfg struct {
	Zk     []string `yaml:"zk"`
}

func main() {
	// Default random port.
	addr := flag.String("addr", "127.0.0.1:", "Master or slave listening address")
	tp := flag.String("type", "slave", "Master or slave")
	dataDir := flag.String("data", "", "Data directory; not end with /")
	// Zk info is stored in config file.
	var cfg Cfg
	err := utils.LoadConfig(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	if len(*addr) == 0 {
		log.Fatalln("Empty address")
	}
	// Default host.
	if (*addr)[0] == ':' {
		*addr = "127.0.0.1" + *addr
	}
	if *dataDir != "" && (*dataDir)[len(*dataDir) - 1] == '/' {
		*dataDir = (*dataDir)[:len(*dataDir) - 1]
	}
	// Connect to zk.
	zkConn, _, err := zk.Connect(cfg.Zk, 5 * time.Second)
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
	port := splitAddr[len(splitAddr) - 1]
	// My ip and port.
	name := host + ":" + port
	log.Printf("Server name: %s\n", name)
	grpcServer := grpc.NewServer()
	// Create lock znode.
	exist, _, err := zkConn.Exists("/slock")
	if err != nil {
		log.Fatalln(err)
	}
	if !exist {
		_, err = zkConn.Create("/slock", nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalln(err)
		}
	}
	exist, _, err = zkConn.Exists("/slock/register")
	if err != nil {
		log.Fatalln(err)
	}
	if !exist {
		_, err = zkConn.Create("/slock/register", nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalln(err)
		}
	}
	exist, _, err = zkConn.Exists("/slock/master")
	if err != nil {
		log.Fatalln(err)
	}
	if !exist {
		_, err = zkConn.Create("/slock/master", nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Fatalln(err)
		}
	}
	var masterLock string
	if *tp == "master" {
		log.Println("Trying to grab the lock")
		masterLock, err = zookeeper.Lock(zkConn, "master")
	}
	relock, err := zookeeper.Lock(zkConn, "register")
	if err != nil {
		log.Fatalln(err)
	}
	tr := false
	// Initialization and registration.
	switch *tp {
	case "master":
		err = server.InitMaster(grpcServer, zkConn, name)
	case "slave":
		err = server.InitSlave(grpcServer, zkConn, name, *dataDir)
	case "slave-sb":
		trans := make(chan bool)
		go func() {
			tr = <-trans
		}()
		err = server.InitSlaveSb(grpcServer, zkConn, name, *dataDir, trans)
	case "master-sb":
	default:
		err = fmt.Errorf("wrong type: only master or slave is supported")
	}
	if err != nil {
		log.Fatalln("Fail to init slave", err)
	}
	log.Printf("Server run on addr: %s\n", lis.Addr())
	err = zookeeper.UnLock(zkConn, relock)
	errChan := make(chan error)
	go func() {
		if err :=  grpcServer.Serve(lis); err != nil {
			errChan <- err
			log.Fatal("Fail to run the service", err)
			return
		}
	}()

	// Listen for interrupt signals.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Shutting down the service...", <-errChan)


	// Stop the panicking.
	r := recover()
	if r != nil {
		log.Println("Recovering from", r)
	}
	// Begin deregistration.
	delock, err := zookeeper.Lock(zkConn, "register")
	if err != nil {
		log.Println(err)
	}

	switch *tp {
	case "master":
		err = server.DeregisterMaster(zkConn)
		if err != nil {
			log.Println(err)
		}
		err = zookeeper.UnLock(zkConn, masterLock)
	case "slave":
		err = server.DeregisterSlave(zkConn, *dataDir, name)
	case "slave-sb":
		if tr {
			err = server.DeregisterSlave(zkConn, *dataDir, name)
		}
	case "master-sb":
	}
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Deregistered %s %s\n", *tp, name)
	}
	err = zookeeper.UnLock(zkConn, delock)
	if err != nil {
		log.Fatalln(err)
	}
}
