package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func getConnect() {
	zkList := []string{"localhost:12181"}
	conn, _, err := zk.Connect(zkList, 10*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	c, err := conn.Create("/go_servers", nil, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(c)
	c, err = conn.Create("/testadaadsasdsaw", nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Print(c)
	children, _, err := conn.Children("/go_servers")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v \n", children)
}

func RegistServer(conn *zk.Conn, host string) (err error) {
	_, err = conn.Create("/go_servers/"+host, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return
}

func GetServerList(conn *zk.Conn) (list []string, err error) {
	list, _, err = conn.Children("/go_servers")
	return
}

func main() {
	getConnect()
}
