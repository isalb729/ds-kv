package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func initSlave(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	//rpc.RegisterMetaServer(grpcServer, &server.Master{})
	return nil
}
