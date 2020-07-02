package server

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"sort"
)

/**
 * Initialze and register.
 * @param grpcServer: to handle rpc
 * @param conn: zkclient
 * @param addr: master address
 * @return err
 */
func InitMaster(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {

	log.Println("Now I'm the master.")
	// Register.
	err := registerMaster(conn, addr)
	if err != nil {
		return err
	}
	log.Printf("Registered master: %s\n", addr)
	// Get slave list and add label.
	slaveList, _, err := getSlaveList(conn)
	if err != nil {
		return err
	}
	log.Println("Get slave list: ", slaveList)
	log.Println("Number of slave list: ", len(slaveList))
	masterHandler := rpc.Master{
		SlaveList: slaveList,
		Working:   map[string]bool{},
	}
	// Register rpc.
	pb.RegisterMasterServer(grpcServer, &masterHandler)
	// Listen for data node changes.
	list, _, event, err := conn.ChildrenW("/data")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			e := <-event
			list, _, event, err = conn.ChildrenW("/data")
			if err != nil {
				panic(err)
			}
			if e.Type == zk.EventNodeChildrenChanged {
				// Relabel.
				masterHandler.SlaveList, _, err = getSlaveList(conn)
				if err != nil {
					panic(err)
				}
				log.Println("SlaveList changed: ", masterHandler.SlaveList)
				if len(list) == 0 {
					// Data will be lost.
					go notifyTransData(conn)
					continue
				}
				// Detect failures.
				for k, v := range masterHandler.Working {
					if v && !utils.Include(list, k) {
						// Something terrible happened.
						go notifyTransData(conn)
						masterHandler.Working[k] = false
						break
					}
				}
				for _, v := range list {
					masterHandler.Working[v] = true
				}
			}
		}
	}()
	return nil
}

/**
 * Notify the standby data node to transform.
 * @param conn: zk client.
 */
func notifyTransData(conn *zk.Conn)  {
	// Delete one standby znode.
	sb, _, err := conn.Children("/sb")
	if err != nil {
		log.Println(err)
		return
	}
	if len(sb) == 0 {
		log.Println("No standby data node available")
		return
	}
	err = conn.Delete("/sb/" + sb[0], -1)
	if err != nil {
		log.Println(err)
	}
}

/**
 * Register master.
 * @param conn
 * @param addr: master address
 * @return err
 */
func registerMaster(conn *zk.Conn, addr string) (err error) {
	// Create master znode if not exist.
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
	// Ephemeral.
	_, err = conn.Create("/master/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}

// Nothing to do here.
func DeregisterMaster(conn *zk.Conn) error {
	return nil
}

/*
 * Get and label the slaves here.
 * @return slaveMeta: all slave meta
 * @return labelList: sorted label list
 * @return err
 */
func getSlaveList(conn *zk.Conn) ([]rpc.SlaveMeta, []int, error) {
	list, _, err := conn.Children("/data")
	if err != nil {
		return nil, nil, err
	}
	slaveList := make([]rpc.SlaveMeta, len(list))
	labelList := make([]int, 0)
	for i, host := range list {
		slaveList[i].Host = host
		// Get the label data.
		val, _, err := conn.Get("/data/" + host)
		if err != nil {
			return nil, nil, err
		}
		if string(val) != "" {
			intV, err := utils.Str2Int(string(val))
			if err != nil {
				return nil, nil, err
			}
			slaveList[i].Label = intV
			labelList = append(labelList, intV)
		} else {
			// Initialize to -1.
			slaveList[i].Label = -1
		}
	}
	// Sort it to a ring.
	sort.Ints(labelList)
	for i, slave := range slaveList {
		if slave.Label == -1 {
			// calculate the label, and add it to the labelList
			labelList, slaveList[i].Label = utils.Label(labelList)
			_, err := conn.Set("/data/"+slave.Host, []byte(utils.Int2str(slaveList[i].Label)), -1)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return slaveList, labelList, err
}
