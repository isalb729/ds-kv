package server

import (
	"context"
	"github.com/isalb729/ds-kv/src/rpc"
)

type Master struct {
	slaveList []string
}
/**
 * Implement metaHandler interface.
 */
func (m Master) GetSlave(ctx context.Context, request *rpc.GetSlaveRequest) (*rpc.GetSlaveResponse, error) {
	// jump consistent hash and return key
	panic("implement me")
}



