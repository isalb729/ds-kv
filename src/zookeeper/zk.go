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
//
//func Create() {
//	var acls = zk.WorldACL(zk.PermAll)
//	p, err_create := conn.Create(path1, data1, zk.FlagEphemeral, acls)
//}
//
//func Get() {
//
//}
//func watch() {
//	for {
//		select {
//		case ch_event := <-child_ch:
//			{
//				fmt.Println("path:", ch_event.Path)
//				fmt.Println("type:", ch_event.Type.String())
//				fmt.Println("state:", ch_event.State.String())
//
//				if ch_event.Type == zk.EventNodeCreated {
//					fmt.Printf("has node[%s] detete\n", ch_event.Path)
//				} else if ch_event.Type == zk.EventNodeDeleted {
//					fmt.Printf("has new node[%d] create\n", ch_event.Path)
//				} else if ch_event.Type == zk.EventNodeDataChanged {
//					fmt.Printf("has node[%d] data changed", ch_event.Path)
//				}
//			}
//		}
//
//		time.Sleep(time.Millisecond * 10)
//	}
//		// watch the child events
//		children, s, child_ch, err := conn.ChildrenW(root_path)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//}
//
