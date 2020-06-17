package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/samuel/go-zookeeper/zk"
	"log"
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

func watch() {
	fmt.Printf("ZKOperateWatchTest\n")
	// callback is a function
	option := zk.WithEventCallback(callback)
	var hosts = []string{"localhost:2181"}
	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	var path1 = "/zk_test_go1"
	var data1 = []byte("zk_test_go1_data1")
	exist, s, _, err := conn.ExistsW(path1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("path[%s] exist[%t]\n", path1, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))

	// try create
	var acls = zk.WorldACL(zk.PermAll)
	p, err_create := conn.Create(path1, data1, zk.FlagEphemeral, acls)
	if err_create != nil {
		fmt.Println(err_create)
		return
	}
	fmt.Printf("created path[%s]\n", p)

	time.Sleep(time.Second * 2)

	exist, s, _, err = conn.ExistsW(path1)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("path[%s] exist[%t] after create\n", path1, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))

	// delete
	conn.Delete(path1, s.Version)

	exist, s, _, err = conn.ExistsW(path1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("path[%s] exist[%t] after delete\n", path1, exist)
	fmt.Printf("state:\n")
	fmt.Printf("%s\n", ZkStateStringFormat(s))
}

func watchChildren() {
	fmt.Printf("ZkChildWatchTest")

	var hosts = []string{"localhost:2181"}
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// try create root path
	var root_path = "/test_root"

	// check root path exist
	exist, _, err := conn.Exists(root_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !exist {
		fmt.Printf("try create root path: %s\n", root_path)
		var acls = zk.WorldACL(zk.PermAll)
		p, err := conn.Create(root_path, []byte("root_value"), 0, acls)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("root_path: %s create\n", p)
	}

	// try create child node
	cur_time := time.Now().Unix()
	ch_path := fmt.Sprintf("%s/ch_%d", root_path, cur_time)
	var acls = zk.WorldACL(zk.PermAll)
	p, err := conn.Create(ch_path, []byte("ch_value"), zk.FlagEphemeral, acls)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("ch_path: %s create\n", p)

	// watch the child events
	children, s, child_ch, err := conn.ChildrenW(root_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("root_path[%s] child_count[%d]\n", root_path, len(children))
	for idx, ch := range children {
		fmt.Printf("%d, %s\n", idx, ch)
	}

	fmt.Printf("watch children result state[%s]\n", ZkStateString(s))

	for {
		select {
		case ch_event := <-child_ch:
			{
				fmt.Println("path:", ch_event.Path)
				fmt.Println("type:", ch_event.Type.String())
				fmt.Println("state:", ch_event.State.String())

				if ch_event.Type == zk.EventNodeCreated {
					fmt.Printf("has node[%s] detete\n", ch_event.Path)
				} else if ch_event.Type == zk.EventNodeDeleted {
					fmt.Printf("has new node[%d] create\n", ch_event.Path)
				} else if ch_event.Type == zk.EventNodeDataChanged {
					fmt.Printf("has node[%d] data changed", ch_event.Path)
				}
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}






type Cfg struct {
	Zk     []string `yaml:"zk"`
	Master struct {
		Port int `yaml:"port"`
	} `yaml:"master"`
}

// todo: remove, change package to zookeeper
func main() {
	var cfg Cfg
	err := utils.LoadConfig(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := Connect(cfg.Zk)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	path := "/zk_test"
	data := []byte("hello")
	flags := int32(0)
	var acls = zk.WorldACL(zk.PermAll)
	p, err := conn.Create(path, data, flags, acls)
	// and p is the path(connsider the case of empherral)

	v, s, err := conn.Get(path)
	// v is the value, s is the status



	//// exist
	//exist, s, err := conn.Exists(path)

	//// update
	//var new_data = []byte("zk_test_new_value")
	//s, err = conn.Set(path, new_data, s.Version)


	//// get
	//v, s, err = conn.Get(path)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Printf("new value of path[%s]=[%s].\n", path, v)
	//fmt.Printf("state:\n")
	//fmt.Printf("%s\n", ZkStateStringFormat(s))
	//
	//// delete
	//err = conn.Delete(path, s.Version)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}


}


