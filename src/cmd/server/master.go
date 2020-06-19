package main

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
)

func InitMaster(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	// register lock
	err := registerMaster(conn, addr)
	if err != nil {
		return err
	}
	/* Get slave list and fill. */
	// register lock
	slaveList, err := getSlaveList(conn)
	// Listen slave list.
	// add label
	pb.RegisterMetaServer(grpcServer, &rpc.Master{
		//SlaveList: slaveList,
	})
	return nil
}

func registerMaster(conn *zk.Conn, addr string) (err error) {
	exist, _, err := conn.Exists("/master")
	if err != nil {
		return err
	}
	if !exist {
		_, err := conn.Create("/master", nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	_, err = conn.Create("/master/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}


// TODO: replica master
func deregisterMaster(conn *zk.Conn, addr string) (err error) {
	return err
}

func getSlaveList(conn *zk.Conn) (list []string, err error) {
	list, _, err = conn.Children("/slave")
	return list, err
}