package main

import (
	"github.com/isalb729/ds-kv/src/rpc"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/samuel/go-zookeeper/zk"
	"google.golang.org/grpc"
	"log"
	"sort"
	"time"
)

const (
	Interval = 2 * time.Second
)

func InitMaster(grpcServer *grpc.Server, conn *zk.Conn, addr string) error {
	// TODO: lock register
	// defer lock
	err := registerMaster(conn, addr)
	if err != nil {
		return err
	}
	log.Printf("Registered master: %s\n", addr)
	// get slave list and add label
	slaveList, err := getSlaveList(conn)
	if err != nil {
		return err
	}
	log.Println("Get slave list: ", slaveList)
	log.Println("Number of slave list: ", len(slaveList))
	masterHandler := rpc.Master{SlaveList: slaveList}
	pb.RegisterMetaServer(grpcServer, &masterHandler)
	//  Listen slave list.
	_, _, event, err := conn.ChildrenW("/data")
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
				/* after this event channel is closed */
				case e := <-event:
					if e.Type == zk.EventNodeChildrenChanged {
						masterHandler.SlaveList, err = getSlaveList(conn)
						if err != nil {
							panic(err)
						}
						log.Println("SlaveList changed: ", masterHandler.SlaveList)
					}
					_, _, event, err = conn.ChildrenW("/data")
					if err != nil {
						// TODO: do deregister after fatal or use panic and recover
						// When doing master replica
						log.Fatalln(err)
					}

			}
			time.Sleep(Interval)
		}
	}()
	// TODO: unlock register
	return nil
}

func registerMaster(conn *zk.Conn, addr string) (err error) {
	// create master directory if not exist
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
	// ephemeral
	_, err = conn.Create("/master/"+addr, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	return err
}


// TODO: replica master
func deregisterMaster(conn *zk.Conn, addr string) (err error) {
	return err
}

/*
 * label the slaves here
 * @return []rpc.SlaveMeta
 */
func getSlaveList(conn *zk.Conn) ([]rpc.SlaveMeta, error) {
	list, _, err := conn.Children("/data")
	if err != nil {
		return nil, err
	}
	slaveList := make([]rpc.SlaveMeta, len(list))
	labelList := make([]int, 0)
	for i, host := range list {
		slaveList[i].Host = host
		val, _, err := conn.Get("/data/" + host)
		if err != nil {
			return nil, err
		}
		if string(val) != "" {
			intV, err := utils.Str2Int(string(val))
			if err != nil {
				return nil, err
			}
			slaveList[i].Label = intV
			labelList = append(labelList, intV)
		} else {
			// initialize to -1
			slaveList[i].Label = -1
		}
	}
	// sort it to a ring
	sort.Ints(labelList)
	for i, slave := range slaveList {
		if slave.Label == -1 {
			// calculate the label, and add it to the labelList
			labelList, slaveList[i].Label = utils.Label(labelList)
			_, err := conn.Set("/data/" + slave.Host, []byte(utils.Int2str(slaveList[i].Label)), -1)
			if err != nil {
				return nil, err
			}
		}
	}
	return slaveList, err
}

