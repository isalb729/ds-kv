package main

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
)

const (
	StoreLevel = 2
)

func InitSlave(grpcServer *grpc.Server, conn *zk.Conn, addr string, dataDir string) error {
	// TODO: lock register
	err := registerSlave(conn, addr)
	if err != nil {
		return err
	}
	log.Printf("registered slave: %s\n", addr)
	err = utils.CreateDataDir(dataDir)
	if err != nil {
		return err
	}
	// set data direct
	pb.RegisterKvServer(grpcServer, &rpc.KvOp{
		DataDir:    dataDir,
		StoreLevel: StoreLevel,
	})
	return nil
}

func registerSlave(conn *zk.Conn, addr string) (err error) {
	exist, _, err := conn.Exists("/data")
	if err != nil {
		return err
	}
	if !exist {
		_, err := conn.Create("/data", nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}
	_, err = conn.Create("/data/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

// TODO
func deregisterSlave(conn *zk.Conn, dataDir string) error {

	err := utils.DeleteDataDir(dataDir)
	return err
}
