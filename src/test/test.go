package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func main() {
	conn, _, _ := zk.Connect([]string{"127.0.0.1:2181"}, 5*time.Second)
	ch := make(chan bool, 3)
	p1, err := zookeeper.Lock(conn, "register")
	fmt.Println(p1, err)
	time.Sleep(5 * time.Second)
	zookeeper.UnLock(conn, p1)
	ch <- true
	go func() {
		p2, err := zookeeper.Lock(conn, "register")
		fmt.Println(p2, err)
		time.Sleep(1 * time.Second)
		zookeeper.UnLock(conn, p2)
		ch <- true
	}()
	go func() {
		p3, err := zookeeper.Lock(conn, "register")
		fmt.Println(p3, err)
		time.Sleep(3 * time.Second)
		zookeeper.UnLock(conn, p3)
		ch <- true
	}()
	<-ch
	<-ch
	<-ch

	//path, _ = conn.CreateProtectedEphemeralSequential("/lock/" + "1", nil, zk.WorldACL(zk.PermAll))
	//fmt.Println(path)
	//time.Sleep(10*time.Second)
}
