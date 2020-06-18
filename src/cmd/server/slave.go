package main

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)



func InitSlave(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	// lock
	err := registerSlave(conn, addr)
	if err != nil {
		return err
	}
	// lock
	getMaster(conn)
	pb.RegisterKvServer(grpcServer, &rpc.KvOp{})
	// todo: advanced: listen master
	return nil
}

func getSlaveList(conn *zk.Conn) (list []string, err error) {
	list, _, err = conn.Children("/slave")
	return list, err
}
func registerSlave(conn *zk.Conn, addr string) (err error) {
	_, err = conn.Create("/data/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}
func deregisterSlave(conn *zk.Conn, addr string) (err error) {
	_, err = conn.Create("/data/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

func getMaster(conn *zk.Conn) (string, error) {
	//list, _, err := conn.Children("/master")
	//return list, err
}
