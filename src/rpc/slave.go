package rpc

import (
	"context"
	"github.com/isalb729/ds-kv/src/rpc/pb"
)

type KvOp struct {
	MasterList []string
}

func (s KvOp) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	panic("implement me")
}

func (s KvOp) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	panic("implement me")
}

func (s KvOp) Del(ctx context.Context, request *pb.DelRequest) (*pb.DelResponse, error) {
	panic("implement me")
}

