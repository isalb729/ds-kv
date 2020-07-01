package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

func main() {
	conn, _, _ := zk.Connect([]string{"localhost:2181"}, 5*time.Second)
	p, _, _ := conn.Children("/sb")
	log.Println(p)

}
