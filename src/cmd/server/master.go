package main

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/server"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func initMaster(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	err := registerMaster(conn, addr)
	if err != nil {
		return err
	}
	/* Get slave list and fill. */
	// TODO: LOCK
	slaveList, err := getSlaveList(conn)
	/* Listen slave list. */
	rpc.RegisterMasterServer(grpcServer, &server.Master{
		SlaveList: slaveList,
	})
	return nil
}

func registerMaster(conn *zk.Conn, addr string) (err error) {
	_, err = conn.Create("/master/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

func getSlaveList(conn *zk.Conn) (list []string, err error) {
	list, _, err = conn.Children("/slave")
	return list, err
}
