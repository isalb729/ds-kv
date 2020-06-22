package main

import (
	"context"
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
	err := registerSlave(conn, addr, dataDir)
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

func registerSlave(conn *zk.Conn, addr, dataDir string) (err error) {
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
	addrList, myMeta, err := getAdjacent(conn, addr)
	log.Println("Adjacent servers: ", addrList)
	if err != nil {
		log.Println(err)
		return err
	}
	data := map[string]interface{}{}
	for _, v := range addrList {
		kvConn, err := grpc.Dial(v.Host, grpc.WithInsecure())
		if err != nil {
			return err
		}
		kvClient := pb.NewKvClient(kvConn)
		rsp, err := kvClient.MoveData(context.Background(), &pb.MoveDataRequest{
			FromLabel: int32(v.Label),
			ToLabel:   int32(myMeta.Label),
		})
		if err != nil {
			log.Println(err)
			return err
		}
		for _, kv := range rsp.Kvs {
			data[kv.Key] = kv.Value
		}
	}
	err = writeLocal(data, dataDir)
	return err
}

func writeLocal(data map[string]interface{}, dataDir string) error {
	for k, v := range data {
		err, path := utils.GetPath(dataDir, k, StoreLevel)
		if err != nil {
			return err
		}
		data := map[string]interface{}{}
		err = utils.ReadMap(path, &data)
		if err != nil {
			return err
		}
		err = utils.AppendMap(path, map[string]interface{}{k: v})
		if err != nil {
			return err
		}
	}
	return nil
}

func deregisterSlave(conn *zk.Conn, dataDir, addr string) error {
	// TODO: lock the key
	// redistribute
	addrList, _, err := getAdjacent(conn, addr)
	paths, err := utils.ReadAllFiles(dataDir)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, path := range paths {
		data := map[string]interface{}{}
		err = utils.ReadMap(path, &data)
		if err != nil {
			log.Println(err)
			return err
		}
		for k, v := range data {
			var addr string
			if len(addrList) == 0 {
				log.Println("Warning: there are no server left, data will be lost.")
				return nil
			} else if len(addrList) == 1 {
				addr = addrList[0].Host
			} else if utils.ShouldBeMoved(k, int32(addrList[0].Label), int32(addrList[1].Label)) {
				addr = addrList[0].Host
			} else {
				addr = addrList[1].Host
			}
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			if err != nil {
				return err
			}
			kvClient := pb.NewKvClient(conn)
			_, err = kvClient.Put(context.Background(), &pb.PutRequest{
				Key:   k,
				Value: v.(string),
			})
			if err != nil {
				return err
			}
		}
	}
	err = utils.DeleteDataDir(dataDir)
	return err
}

// adjacent, my meta
func getAdjacent(conn *zk.Conn, addr string) ([]rpc.SlaveMeta, *rpc.SlaveMeta, error) {
	slaveList, labelList, err := getSlaveList(conn)
	if err != nil {
		return nil, nil, err
	}
	i, j := -1, -1
	for k, v := range slaveList {
		if v.Host == addr {
			i = k
			break
		}
	}
	if i == -1 {
		return nil, nil, fmt.Errorf("label not found")
	}
	for k, v := range labelList {
		if v == slaveList[i].Label {
			j = k
			break
		}
	}
	if j == -1 {
		return nil, nil, fmt.Errorf("label not found")
	}
	serverNum := len(labelList)
	left := (j - 1 + serverNum) % serverNum
	right := (j + 1) % serverNum
	var adjServer []rpc.SlaveMeta
	if left != j {
		h := -1
		for k, v := range slaveList {
			if v.Label == labelList[left] {
				h = k
				break
			}
		}
		adjServer = append(adjServer, slaveList[h])
	}
	if right != j && right != left {
		h := -1
		for k, v := range slaveList {
			if v.Label == labelList[right] {
				h = k
				break
			}
		}
		adjServer = append(adjServer, slaveList[h])
	}
	return adjServer, &slaveList[i], nil
}
