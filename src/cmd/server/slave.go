package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func initSlave(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	err := registerSlave(conn, addr)
	if err != nil {
		return err
	}
	 getMasterList(conn)
	//rpc.RegisterMetaServer(grpcServer, &server.Master{})
	return nil
}

func registerSlave(conn *zk.Conn, addr string) (err error) {
	_, err = conn.Create("/master/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

func getMasterList(conn *zk.Conn) (list []string, err error) {
	list, _, err = conn.Children("/master")
	return list, err
}
