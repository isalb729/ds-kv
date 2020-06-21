package rpc

import (
	"context"
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"log"
)

type KvOp struct {
	DataDir    string
	StoreLevel int
}



func getPath(base, key string, storeLevel int) (error, string) {
	if storeLevel < 1 {
		return fmt.Errorf("a store level is supposed to be at least one"), ""
	}
	primes := utils.GetPrimes(storeLevel, 3)
	hash := int(utils.BasicHash(key))
	path := base
	for _, v := range primes {
		path = path + "/" + utils.Int2str(hash % v)
	}
	return nil, path
}

func (kv KvOp) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	// TODO lock
	key := request.Key
	err, path := getPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data := map[string]interface{}{}
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if data[key] == nil {
		return &pb.GetResponse{
			Ok:                   false,
			Value:                "",
		}, nil
	} else {
		return &pb.GetResponse{
			Ok:                   true,
			Value:                data[key].(string),
		}, nil
	}
}

func (kv KvOp) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	// TODO LOCK
	key := request.Key
	val := request.Value
	err, path := getPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	created := true

	data := map[string]interface{}{}
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if data[key] != nil {
		if data[key] == val {
			return &pb.PutResponse{
				Created:              false,
			}, nil
		}
		created = false
	}

	err = utils.AppendMap(path, map[string]interface{}{key: val})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.PutResponse{
		Created:              created,
	}, nil
}

func (kv KvOp) Del(ctx context.Context, request *pb.DelRequest) (*pb.DelResponse, error) {
	// TODO lock
	key := request.Key
	err, path := getPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if data[key] == nil {
		return &pb.DelResponse{
			Deleted:              false,
		}, nil
	} else {
		delete(data, key)
		err = utils.WriteMap(path, data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &pb.DelResponse{
			Deleted:              true,
		}, nil
	}
}

func shouldBeMoved(key string, toLabel, fromLabel int32) bool {
	label := utils.BasicHash(key) % utils.HoleNum
	dto := utils.Dist(int(label), int(toLabel), utils.HoleNum)
	dfrom := utils.Dist(int(label), int(fromLabel), utils.HoleNum)
	return dto < dfrom || dto == dfrom && toLabel < fromLabel
}

func (kv KvOp) MoveData(ctx context.Context, request *pb.MoveDataRequest) (*pb.MoveDataResponse, error) {
	paths, err := utils.ReadAllFiles(kv.DataDir)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var kvs []*pb.MoveDataResponse_Kv
	for _, path := range paths {
		data := map[string]interface{}{}
		err = utils.ReadMap(path, &data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for k, v := range data {
			kvs = append(kvs, &pb.MoveDataResponse_Kv{
				Key:                  k,
				Value:                v.(string),
			})
			if shouldBeMoved(v.(string),  request.ToLabel, request.FromLabel) {
				delete(data, k)
			}
		}
		err = utils.WriteMap(path, data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return &pb.MoveDataResponse{
		Kvs:                  kvs,
	}, nil
}
