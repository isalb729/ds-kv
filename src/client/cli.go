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
		err, _ := cli.put(key, val)
		fmt.Printf("put key: %s val: %s err: %v\n", key, val, err)
	}

	get := func(key string, num int) {
		err, val := cli.get(key)
		fmt.Printf("%d get key: %s, err: %v, val: %s\n",num, key, err, val)
	}

	//del := func(key string) {
	//	err, _ := cli.del(key)
	//	fmt.Printf("del key: %s err: %v\n", key, err)
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
	//exit := make(chan int, 10)
	//go func() {
	//	put("os", "100")
	//	exit<-1
	//}()
	//go func() {
	//	put("os", "99")
	//	exit<-1
	//}()
	//go func() {
	//	put("os", "98")
	//	exit<-1
	//}()
	//go func() {
	//	put("os", "97")
	//	exit<-1
	//}()
	//go func() {
	//	get("os", 1)
	//	get("os", 1)
	//	get("os", 1)
	//	get("os", 1)
	//	get("os", 1)
	//	exit<-1
	//}()
	//go func() {
	//	get("os", 2)
	//	get("os", 2)
	//	get("os", 2)
	//	get("os", 2)
	//	get("os", 2)
	//	exit<-1
	//}()
	//go func() {
	//	get("os", 3)
	//	get("os", 3)
	//	get("os", 3)
	//	get("os", 3)
	//	exit<-1
	//}()
	//<-exit
	//<-exit
	//<-exit
	//<-exit
	//<-exit
	//<-exit
	//<-exit
	//dumpAll()

	for {
	put("os", "100")
	put("ds", "98")
	put("ca", "97")
	put("st", "96")
	get("os", 1)
	get("ds", 1)
	get("ca", 1)
	get("st", 1)
	dumpAll()}
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
