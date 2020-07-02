package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/client"
	"log"
)

func Concurrent(cli *client.KvCli) {
	put := func(key, val string) {
		err, _ := cli.Put(key, val)
		fmt.Printf("put key: %s val: %s err: %v\n", key, val, err)
	}

	get := func(key string, num int) {
		err, val := cli.Get(key)
		fmt.Printf("%d get key: %s, err: %v, val: %s\n", num, key, err, val)
	}
	dumpAll := func() {
		fmt.Println("DUMPING ALL!!!")
		err, rsp := cli.DumpAll()
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
	exit := make(chan int, 10)
	go func() {
		put("os", "100")
		exit <- 1
	}()
	go func() {
		put("os", "99")
		exit <- 1
	}()
	go func() {
		put("os", "98")
		exit <- 1
	}()
	go func() {
		put("os", "97")
		exit <- 1
	}()
	go func() {
		get("os", 1)
		get("os", 1)
		get("os", 1)
		get("os", 1)
		get("os", 1)
		exit <- 1
	}()
	go func() {
		get("os", 2)
		get("os", 2)
		get("os", 2)
		get("os", 2)
		get("os", 2)
		exit <- 1
	}()
	go func() {
		get("os", 3)
		get("os", 3)
		get("os", 3)
		get("os", 3)
		exit <- 1
	}()
	<-exit
	<-exit
	<-exit
	<-exit
	<-exit
	<-exit
	<-exit
	dumpAll()
}

func Sequential(cli *client.KvCli) {
	put := func(key, val string) {
		err, _ := cli.Put(key, val)
		fmt.Printf("put key: %s val: %s err: %v\n", key, val, err)
	}

	get := func(key string) {
		err, val := cli.Get(key)
		fmt.Printf("get key: %s, err: %v, val: %s\n", key, err, val)
	}

	del := func(key string) {
		err, _ := cli.Del(key)
		fmt.Printf("del key: %s err: %v\n", key, err)
	}

	dumpAll := func() {
		fmt.Println("DUMPING ALL!!!")
		err, rsp := cli.DumpAll()
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

	for {
		get("os")
		del("ds")
		put("os", "100")
		put("ds", "98")
		put("ca", "97")
		put("st", "96")
		del("st")
		get("os")
		get("ds")
		get("ca")
		get("st")
		dumpAll()
	}
}

func CrazyLoop(cli *client.KvCli)  {
	put := func(key, val string) {
		err, _ := cli.Put(key, val)
		fmt.Printf("put err: %v\n", err)
	}

	get := func(key string) {
		err, _ := cli.Get(key)
		fmt.Printf("get err: %s\n", err)
	}

	del := func(key string) {
		err, _ := cli.Del(key)
		fmt.Printf("del err: %v\n", err)
	}

	for {
		get("os")
		del("ds")
		put("os", "100")
		put("ds", "98")
		put("ca", "97")
		put("st", "96")
		del("st")
		get("os")
		get("ds")
		get("ca")
		get("st")
	}
}