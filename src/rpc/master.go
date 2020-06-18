package rpc

import (
	"context"
	"github.com/isalb729/ds-kv/src/rpc/pb"
)
/*
 *  Master has to take good care of label.
 */
type SlaveMeta struct {
	label string
	host string
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

}
