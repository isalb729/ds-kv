package main

import (
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/utils"
	"google.golang.org/grpc"
	"log"
)
type Cfg struct {
	Zk     []string `yaml:"zk"`
	Master struct {
		Port int `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"master"`
}

func main() {
	var cfg Cfg
	err := utils.LoadConfig(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	masterAddr := fmt.Sprintf("%s:%d", cfg.Master.Host, cfg.Master.Port)
	conn, err := grpc.Dial(masterAddr)
	if err != nil {
		log.Fatalln(err)
	}
	client := rpc.NewAuthClient(conn)
	client.OnLogin()
	err = conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
