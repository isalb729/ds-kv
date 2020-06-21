package client

import (
	"fmt"
)

//import "github.com/isalb729/ds-kv/src/rpc"

func Concurrent(cli *KvCli) {

}


func Sequential(cli *KvCli) {
	put := func (key, val string) {
		fmt.Printf("put key: %s val: %s\n", key, val)
		fmt.Println(cli.put(key, val))
	}

	get := func(key string) {
		fmt.Printf("get key: %s\n", key)
		fmt.Println(cli.get(key))
	}

	del := func(key string) {
		fmt.Printf("del key: %s\n", key)
		fmt.Println(cli.del(key))
	}

	put("os", "100")
	put("ds", "98")
	put("ca", "97")
	put("st", "96")
	get("os")
	put("os", "99")
	get("os")
	del("os")
	get("ds")
	get("ds")
	get("st")
	get("os")
	//go func() {
	//	put("os", "100")
	//	put("os", "100")
	//	put("os", "100")
	//	put("os", "100")
	//	put("os", "100")
	//	put("os", "100")
	//}()
	//go func() {
	//	put("ds", "100")
	//	put("ds", "100")
	//	put("ds", "100")
	//	put("ds", "100")
	//	put("ds", "100")
	//	put("ds", "100")
	//	put("ds", "100")
	//	put("ds", "100")
	//}()
	//time.Sleep(20 * time.Second)
}
