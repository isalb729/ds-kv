package main

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/server"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func initMaster(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	/* Register. */
	/* Get slave list and fill. */
	/* Listen slave list. */
	rpc.RegisterMetaServer(grpcServer, &server.Master{})
	return nil
}

//func RegistServer(conn *zk.Conn, host string) (err error) {
//	_, err = conn.Create("/go_servers/"+host, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
//	return
//}
//
//func GetServerList(conn *zk.Conn) (list []string, err error) {
//	list, _, err = conn.Children("/go_servers")
//	return
//}
//