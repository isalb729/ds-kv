package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func Connect(zkList []string) (*zk.Conn, error) {
	conn, _, err := zk.Connect(zkList, 10*time.Second)
	return conn, err
}

/**
 *
 */
func RegisterServer() {

}

/**
 *
 */
func DeregisterServer() {

}


