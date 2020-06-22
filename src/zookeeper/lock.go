package zookeeper

import (
	"github.com/isalb729/ds-kv/src/utils"
	"github.com/samuel/go-zookeeper/zk"
	"sort"
)

/*
 * create sequential node;wait until to be the first;doop;delete
 */

// https://zookeeper.apache.org/doc/r3.3.5/zookeeperProgrammers.html
// https://zookeeper.apache.org/doc/r3.3.5/recipes.html
// sequential lock
func Lock(conn *zk.Conn, name string) (string, error) {
	path, err := conn.CreateProtectedEphemeralSequential("/slock/" + name + "/0", nil, zk.WorldACL(zk.PermAll))
	if err != nil {
		return "", err
	}
	list, _, err := conn.Children("/slock/" + name)
	if err != nil {
		return "", err
	}
	last := getLastName(path[len(utils.ParseDir(path)) + 1:], list)
	if last != "" {
		_, _, event, err := conn.ExistsW("/slock/" + name + "/" +  last)
		if err != nil {
			return "", err
		}
		for e := <-event; ; {
			list, _, err := conn.Children("/slock/" + name)
			if err != nil {
				return "", err
			}
			last = getLastName(path[len(utils.ParseDir(path)) + 1:], list)
			if e.Type == zk.EventNodeDeleted && last == ""{
				break
			}
			_, _, event, err = conn.ExistsW("/slock/" + name + "/" +  last)
			if err != nil {
				return "", err
			}
		}
	}
	return path, nil

}

func getId(path string) uint64 {
	id, _ := utils.Str2UInt64(path[len(path) - 10:])
	return id
}

func getLastName(name string, list []string) string {
	sort.Slice(list, func(i, j int) bool {
		return getId(list[i]) < getId(list[j])
	})
	for k, v := range list {
		if v == name {
			if k > 0 {
				return list[k - 1]
			} else {
				return ""
			}
		}
	}
	// possibly the node itself has been deleted by zk
	return ""
}

func UnLock(conn *zk.Conn, name string) error {
	 return conn.Delete(name, -1)
}
