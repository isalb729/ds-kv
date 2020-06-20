package rpc

import (
	"context"
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
)
/*
 *  Master has to take good care of label.
 */
type SlaveMeta struct {
	Label int
	Host string
}

type Master struct {
	SlaveList []SlaveMeta
}


/**
 * client to master
 * send key to get the host of corresponding server
 */
func (m Master) GetSlave(ctx context.Context, request *pb.GetSlaveRequest) (*pb.GetSlaveResponse, error) {
	// consistent hash and return key
	//request.Key
	hash := int(utils.BasicHash(request.Key) % utils.HoleNum)
	dist := utils.HoleNum
	var addr string
	for _, v := range m.SlaveList {
		if utils.Dist(v.Label, hash, utils.HoleNum) < dist {
			dist = utils.Dist(v.Label, hash, utils.HoleNum)
			addr = v.Host
		}
	}
	if addr == "" {
		return nil, fmt.Errorf("fail to get data node address")
	}
	return  &pb.GetSlaveResponse{
		Addr:                 addr,
	}, nil
}
