package main

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

/**
 *
 * @param zkList:
 * @return
 */
func connect(zkList []string) (*zk.Conn, error) {
	conn, _, err := zk.Connect(zkList, 10*time.Second)
	return conn, err

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
	zkList := []string{"localhost:12181"}
	conn, err := connect(zkList)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// ephemeral|sequential
	c, err := conn.Create("/go_serversaa", []byte("123"),  0, zk.WorldACL(zk.PermAll))
	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(c)
	children, _, err := conn.Get("/go_serversaa")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v \n", string(children))

}
