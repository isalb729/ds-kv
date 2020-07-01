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
	"sync"
)

const (
	StoreLevel = 2
)

func InitSlave(grpcServer *grpc.Server, conn *zk.Conn, addr string, dataDir string) error {
	err := registerSlave(conn, addr, dataDir)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Data directory: %s\n", dataDir)
	log.Printf("Registered slave: %s\n", addr)
	err = utils.CreateDataDir(dataDir)
	if err != nil {
		log.Println(err)
		return err
	}
	list, _, event, err := conn.ChildrenW("/sb")
	if err != nil {
		return err
	}
	log.Println("StandByList: ", list)
	kvOp := rpc.KvOp{
		DataDir:    dataDir,
		StoreLevel: StoreLevel,
		RwLock:     map[string]*sync.RWMutex{},
		Sb:         list,
	}
	pb.RegisterDataServer(grpcServer, &kvOp)
	go func() {
		for {
			e := <-event
			list, _, event, err = conn.ChildrenW("/sb")
			if err != nil {
				panic(err)
			}
			if e.Type == zk.EventNodeChildrenChanged {
				kvOp.Sb = list
				log.Println("StandByList changed: ", list)
			}
		}
	}()
	return nil
}

func registerSlave(conn *zk.Conn, addr, dataDir string) (err error) {
	//test register lock
	//time.Sleep(5*time.Second)
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
		kvClient := pb.NewDataClient(kvConn)
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
	err = utils.WriteLocal(data, dataDir, StoreLevel)
	return err
}

func deregisterSlave(conn *zk.Conn, dataDir, addr string) error {
	// redistribute
	defer utils.DeleteDataDir(dataDir)
	done := make(chan bool)
	defer func() {
		<-done
	}()
	go func() {
		list, _, err := conn.Children("/sb")
		if err != nil {
			log.Println(err)
		}
		for _, sb := range list {
			conn, err := grpc.Dial(sb, grpc.WithInsecure())
			if err != nil {
				continue
			}
			sbClient := pb.NewDataStandByClient(conn)
			_, err = sbClient.DeregisterNotify(context.Background(), &pb.DeregisterNotifyRequest{
				Addr: addr,
			})
			if err != nil {
				log.Println(err)
			}
		}
		done <- true
	}()
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
			dataClient := pb.NewDataClient(conn)
			_, err = dataClient.Put(context.Background(), &pb.PutRequest{
				Key:   k,
				Value: v.(string),
			})
			if err != nil {
				return err
			}
		}
	}
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

func InitSlaveSb(grpcServer *grpc.Server, conn *zk.Conn, addr string, dataDir string) error {
	err := registerSlaveSb(conn, addr, dataDir)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Data directory: %s\n", dataDir)
	log.Printf("Registered slave standby: %s\n", addr)
	err = utils.CreateDataDir(dataDir)
	if err != nil {
		log.Println(err)
		return err
	}
	// set data direct
	sb := rpc.Sb{
		DataDir:    dataDir,
		StoreLevel: StoreLevel,
		Lock:       map[string]*sync.Mutex{},
		GLock:      sync.Mutex{},
		Working:    map[string]bool{},
	}
	pb.RegisterDataStandByServer(grpcServer, &sb)
	// listening
	list, _, event, err := conn.ChildrenW("/data")
	if err != nil {
		return err
	}
	for _, v := range list {
		sb.Working[v] = true
	}
	go func() {
		for {
			e := <-event
			list, _, event, err = conn.ChildrenW("/data")
			if err != nil {
				panic(err)
			}
			// only deregister notify is allowed to switch working from true to false
			if e.Type == zk.EventNodeChildrenChanged {
				for k, v := range sb.Working {
					if v && !utils.Include(list, k) {
						// something terrible happened
					}
				}
				for _, v := range list {
					sb.Working[v] = true
				}
			}
		}
	}()
	return nil
}

func registerSlaveSb(conn *zk.Conn, addr, dataDir string) (err error) {
	exist, _, err := conn.Exists("/sb")
	if err != nil {
		log.Println(err)
		return err
	}
	if !exist {
		_, err := conn.Create("/sb", nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			log.Println(err)
			return err
		}
	}
	_, err = conn.Create("/sb/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	_ = utils.DeleteDataDir(dataDir)
	return err
}
