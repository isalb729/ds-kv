package main

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func InitMaster(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	err := registerMaster(conn, addr)
	if err != nil {
		return err
	}
	/* Get slave list and fill. */
	// TODO: LOCK
	slaveList, err := getSlaveList(conn)
	/* Listen slave list. */
	pb.RegisterMetaServer(grpcServer, &rpc.Master{
		//SlaveList: slaveList,
	})
	return nil
}

func registerMaster(conn *zk.Conn, addr string) (err error) {
	_, err = conn.Create("/master/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

// TODO
func deregisterMaster(conn *zk.Conn, addr string) (err error) {
	return err
}
