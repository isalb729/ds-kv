package client

import (
	"fmt"
	"log"
)

//import "github.com/isalb729/ds-kv/src/rpc"

func Concurrent(cli *KvCli) {

}


func Sequential(cli *KvCli) {
	put := func (key, val string) {
		fmt.Printf("put key: %s val: %s\n", key, val)
		fmt.Println(cli.put(key, val))
	}

	//get := func(key string) {
	//	fmt.Printf("get key: %s\n", key)
	//	fmt.Println(cli.get(key))
	//}

	//del := func(key string) {
	//	fmt.Printf("del key: %s\n", key)
	//	fmt.Println(cli.del(key))
	//}

	dumpAll := func() {
		fmt.Println("DUMPING ALL!!!")
		err, rsp := cli.dumpAll()
		if err != nil {
			log.Println(err)
			return
		}
		for _, data := range rsp {
			fmt.Printf("-------Data server %s with label %d-------\n", data.Host, data.Label)
			for _, kvl := range data.Kvls {
				fmt.Printf("    key: %s value: %s label: %d\n", kvl.Key, kvl.Value, kvl.Label)
			}
			fmt.Println()
		}
	}
	put("os", "100")
	put("ds", "98")
	put("ca", "97")
	put("st", "96")
	dumpAll()
	//get("os")
	//get("ds")
	//get("ca")
	//get("st")
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
