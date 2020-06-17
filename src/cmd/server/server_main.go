package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/server/zookeeper"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Cfg struct {
	Zk     []string `yaml:"zk"`
	Master struct {
		Port int `yaml:"port"`
	} `yaml:"master"`
}

func main() {
	var cfg Cfg
	err := utils.LoadConfig(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	zkConn, err := zookeeper.Connect(cfg.Zk)
	if err != nil {
		log.Fatalln(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Master.Port))
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	rpc.RegisterAuthServer(grpcServer, &zk.AuthServer{})
	grpcServer.Serve(lis)
	//go starServer("127.0.0.1:8897")
	//go starServer("127.0.0.1:8898")
	//go starServer("127.0.0.1:8899")
	//
	//a := make(chan bool, 1)
	//<-a
}

//func starServer(port string) {
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", port)
//	fmt.Println(tcpAddr)
//	checkError(err)
//
//	listener, err := net.ListenTCP("tcp", tcpAddr)
//	checkError(err)
//
//	conn, err := example.GetConnect()
//	if err != nil {
//		fmt.Printf(" connect zk error: %s ", err)
//	}
//
//	defer conn.Close()
//	err = example.RegistServer(conn, port)
//	if err != nil {
//		fmt.Printf(" regist node error: %s ", err)
//	}
//
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "Error: %s", err)
//			continue
//		}
//		go handleCient(conn, port)
//	}
//
//	fmt.Println("aaaaaa")
//}
//
//func handleCient(conn net.Conn, port string) {
//	defer conn.Close()
//
//	daytime := time.Now().String()
//	conn.Write([]byte(port + ": " + daytime))
//}
//
//func RegistServer(conn *zk.Conn, host string) (err error) {
//	_, err = conn.Create("/go_servers/"+host, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
//	return
//}
//
//func GetServerList(conn *zk.Conn) (list []string, err error) {
//	list, _, err = conn.Children("/go_servers")
//	return
//}
//
//func main() {
//	zkList := []string{"localhost:12181"}
//	conn, err := connect(zkList)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer conn.Close()
//	// ephemeral|sequential
//	c, err := conn.Create("/go_serversaa", []byte("123"),  0, zk.WorldACL(zk.PermAll))
//	if err != nil {
//		fmt.Print(err)
//		return
//	}
//	fmt.Println(c)
//	children, _, err := conn.Get("/go_serversaa")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Printf("%v \n", string(children))
//
//}
//
