package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func main()  {
	//Changes to that znode trigger the watch and then clear the watch.
	//All of the read operations in ZooKeeper - getData(), getChildren(), and exists() - have the option of setting a watch as a side effect.
	conn, _ := zookeeper.Connect([]string{"127.0.0.1:2181", "127.0.0.1:2281", "127.0.0.1:2381"})
	_, _, event, _ := conn.ExistsW("/a")
	for {
		select {
			/* after this event channel is closed */
			case e := <-event:
				fmt.Println(e.Type==zk.EventNodeChildrenChanged)
				_, _, event, _ = conn.ExistsW("/a")
		}
		time.Sleep(1 * time.Second)
	}

	conn.Close()
}

