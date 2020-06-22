package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/utils"
)

func main() {


	data := map[string]interface{}{"ds":"99"}
	utils.WriteMap("cmd/server/data/server2/0/0", data)
	utils.ReadMap("cmd/server/data/server2/0/0", &data)
	fmt.Println(data)
}


//func main() {
//
//	//acl := zk.WorldACL(zk.PermAll)
//	//path, _ := conn.CreateProtectedEphemeralSequential("/a/c", []byte("wow"), acl)
//	//fmt.Println(path)
//	//path, _ = conn.CreateProtectedEphemeralSequential("/a/c", []byte("wow"), acl)
//	//fmt.Println(path)
//	//path, _ = conn.CreateProtectedEphemeralSequential("/a/c", []byte("wow"), acl)
//
//}
