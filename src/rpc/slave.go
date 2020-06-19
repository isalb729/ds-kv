package rpc

import (
	"context"
	"github.com/isalb729/ds-kv/src/rpc/pb"
)

type KvOp struct {
	MasterList []string
	// TODO
	DataDir string
}

func (s KvOp) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	//err := utils.CreateDataDir("hello")
	////_ = utils.WriteMap("hello/world", map[string]interface{}{"000": "0", "11": "new"})
	//err = utils.AppendMap("hello/world", map[string]interface{}{"a": "1234", "999": "s456"})
	//data := make(map[string]interface{})
	//err = utils.ReadMap("hello/world", &data)
	//fmt.Println(err)
	//fmt.Println(data)
	panic("implement me")
}

func (s KvOp) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	panic("implement me")
}

func (s KvOp) Del(ctx context.Context, request *pb.DelRequest) (*pb.DelResponse, error) {
	panic("implement me")
}

