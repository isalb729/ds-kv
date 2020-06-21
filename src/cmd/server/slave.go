package main

import (
	"fmt"
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
		log.Println(err)
		return err
	}
	log.Printf("registered slave: %s\n", addr)
	err = utils.CreateDataDir(dataDir)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return err
	}
	if !exist {
		_, err := conn.Create("/data", nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Println(err)
			return err
		}
	}
	_, err = conn.Create("/data/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	// don't delete other but get them
	addrList, err := getAdjacent(conn, addr)
	log.Println("Adjacent servers: ", addrList)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func deregisterSlave(conn *zk.Conn, dataDir string) error {
	// lock the key
	// redistribute
	err := utils.DeleteDataDir(dataDir)
	return err
}

func getAdjacent(conn *zk.Conn, addr string) ([]string, error) {
	slaveList, labelList, err := getSlaveList(conn)
	if err != nil {
		return nil, err
	}
	i, j:= -1, -1
	for k, v := range slaveList {
		if v.Host == addr {
			i = k
			break
		}
	}
	if i == -1 {
		return nil, fmt.Errorf("label not found")
	}
	for k, v := range labelList {
		if v == slaveList[i].Label {
			j = k
			break
		}
	}
	if j == -1 {
		return nil, fmt.Errorf("label not found")
	}
	serverNum := len(labelList)
	left := (j - 1 + serverNum) % serverNum
	right := (j + 1) % serverNum
	var adjServer []string
	if left != j {
		h := -1
		for k, v := range slaveList {
			if v.Label == slaveList[left].Label {
				h = k
				break
			}
		}
		adjServer = append(adjServer, utils.Int2str(labelList[h]))
	}
	if right != j && right != left {
		h := -1
		for k, v := range slaveList {
			if v.Label == slaveList[right].Label {
				h = k
				break
			}
		}
		adjServer = append(adjServer, utils.Int2str(labelList[h]))
	}
	return adjServer, nil
}
