package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

// https://zookeeper.apache.org/doc/r3.3.5/zookeeperProgrammers.html
func Connect(zkList []string) (*zk.Conn, error) {
	conn, _, err := zk.Connect(zkList, 10*time.Second)
	return conn, err
}
// https://zookeeper.apache.org/doc/r3.3.5/recipes.html
func lock() {

}