package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/zookeeper"
	"time"
)

func main() {
	//Changes to that znode trigger the watch and then clear the watch.
	//All of the read operations in ZooKeeper - getData(), getChildren(), and exists() - have the option of setting a watch as a side effect.
	conn, _ := zookeeper.Connect([]string{"127.0.0.1:2181", "127.0.0.1:2281", "127.0.0.1:2381"})
	//_, _, event, _ := conn.ExistsW("/a")
	_, _, event, _ := conn.ChildrenW("/a")
	//create/delete /a/b  trigger node children changed, not include the children change
	//_, _, event, _ := conn.ExistsW("/a")
	//set/create/delete /a, not include the children change
	//acl := zk.WorldACL(zk.PermAll)
	//path, _ := conn.CreateProtectedEphemeralSequential("/a/c", []byte("wow"), acl)
	//fmt.Println(path)
	//path, _ = conn.CreateProtectedEphemeralSequential("/a/c", []byte("wow"), acl)
	//fmt.Println(path)
	//path, _ = conn.CreateProtectedEphemeralSequential("/a/c", []byte("wow"), acl)
	//
	//fmt.Println(path)
	////-1 matches any version
	////conn.Set("",[]byte(""), -1)
	//a,_,_:=conn.Children("/a")
	//
	//fmt.Println(a)
	//time.Sleep(20 * time.Second)
	for {
		select {
			/* after this event channel is closed */
			case e := <-event:
				fmt.Println(e)
				_, _, event, _ = conn.ChildrenW("/a")
		}
		time.Sleep(1 * time.Second)
	}
	conn.Close()
}
