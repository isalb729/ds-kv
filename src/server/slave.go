package server

import (
	"context"
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/isalb729/ds-kv/src/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"sync"
)

const (
	StoreLevel = 2
)

/**
 * Initialze and register.
 * @param grpcServer: to handle rpc
 * @param conn: zkclient
 * @param addr: master address
 * @param dataDir: data directory
 * @return err
 */
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
	// Listen standby nodes.
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
	// Listen for standby nodes changes and update.
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

/**
 * Register data node.
 * @param conn
 * @param addr: master address
 * @param dataDir: data directory
 * @return err
 */
func registerSlave(conn *zk.Conn, addr, dataDir string) (err error) {
	// Test register lock.
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
	// Move some data from adjacent to this server.
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

/**
 * Deregister data node.
 * @param conn
 * @param dataDir
 * @param addr
 * @return err
 */
func DeregisterSlave(conn *zk.Conn, dataDir, addr string) error {
	// Redistribute.
	defer utils.DeleteDataDir(dataDir)
	// Wait until notification finishes.
	done := make(chan bool)
	defer func() {
		<-done
	}()
	go func() {
		defer func() {
			done <- true
		}()
		list, _, err := conn.Children("/master")
		if err != nil {
			log.Println(err)
		}
		if len(list) == 0 {
			log.Println("No master node found")
			return
		}
		conn, err := grpc.Dial(list[0], grpc.WithInsecure())
		if err != nil {
			log.Println(err)
		}
		client := pb.NewMasterClient(conn)
		_, err = client.DeregisterNotify(context.Background(), &pb.DeregisterNotifyRequest{
			Addr: addr,
		})
		if err != nil {
			log.Println(err)
		}

	}()

	addrList, _, err := getAdjacent(conn, addr)
	// Move data to adjacent servers.
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

/**
 * @param conn
 * @param addr: my address
 * @return slaveMetaList: all slave metas
 * @return mymeta: mymeta after labeling
 * @return err
 */
func getAdjacent(conn *zk.Conn, addr string) ([]rpc.SlaveMeta, *rpc.SlaveMeta, error) {
	slaveList, labelList, err := getSlaveList(conn)
	if err != nil {
		return nil, nil, err
	}
	i, j := -1, -1
	// Find me.
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
	// Adjacent.
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

/**
 * Initialize standby data node.
 * @param grpcServer: handle rpc
 * @param conn
 * @param addr
 * @param dataDir
 * @param trans: a channel to tell whether it has been transformed to data server
 * @return err
 */
func InitSlaveSb(grpcServer *grpc.Server, conn *zk.Conn, addr string, dataDir string, trans chan bool) error {
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
	}
	pb.RegisterDataStandByServer(grpcServer, &sb)
	kvOp := rpc.KvOp{
		DataDir:    dataDir,
		StoreLevel: StoreLevel,
		RwLock:     map[string]*sync.RWMutex{},
		Sb:         nil,
	}
	pb.RegisterDataServer(grpcServer, &kvOp)
	go func() {
		for {
			// If deleted by master, then transform.
			_, _, event, err := conn.ExistsW("/sb/"+addr)
			if err != nil {
				panic(err)
			}
			e := <-event
			if e.Type == zk.EventNodeDeleted {
				log.Println("Transfer to data node")
				break
			}
		}
		trans <- true
		// Register as data server.
		relock, err := zookeeper.Lock(conn, "register")
		_, err = conn.Create("/data/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		_, _, _ = getSlaveList(conn)
		// Delete irrelevant to data.
		allData := make(map[string]interface{})
		files, err := utils.ReadAllFiles(dataDir)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			data := make(map[string]interface{})
			err = utils.ReadMap(file, &data)
			if err != nil {
				panic(err)
			}
			utils.MergeMap(&allData, data)
		}
		list, _, err := conn.Children("/master")
		if err != nil {
			panic(err)
		}
		if len(list) == 0 {
			panic(err)
		}
		mconn, err := grpc.Dial(list[0], grpc.WithInsecure())
		if err != nil {
			log.Println(err)
		}
		mc := pb.NewMasterClient(mconn)
		for k, _ := range allData {
			rsp, err := mc.GetSlave(context.Background(), &pb.GetSlaveRequest{
				Key:                  k,
			})
			if err != nil {
				log.Println(err)
				continue
			}
			if rsp.Addr != addr {
				delete(allData, k)
			}
		}
		err = utils.DeleteDataDir(dataDir)
		if err != nil {
			log.Println(err)
		}
		err = utils.WriteLocal(allData, dataDir, StoreLevel)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Data directory: %s\n", dataDir)
		log.Printf("Registered slave: %s\n", addr)
		list, _, event, err := conn.ChildrenW("/sb")
		if err != nil {
			return
		}
		log.Println("StandByList: ", list)
		_ = zookeeper.UnLock(conn, relock)
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
	}()
	return nil
}

/**
 * Register standby node.
 * @param conn
 * @param addr
 * @param dataDir
 * @return err
 */
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

