package rpc

import (
	"context"
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"log"
	"sync"
)

type KvOp struct {
	DataDir    string
	StoreLevel int
	// read write lock for each file
	RwLock  map[string]*sync.RWMutex
}

func (kv *KvOp) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	// TODO lock
	fmt.Println("entering get")
	defer fmt.Println("exiting get")
	key := request.Key
	err, path := utils.GetPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data := map[string]interface{}{}
	if kv.RwLock[path] == nil {
		kv.RwLock[path] = new(sync.RWMutex)
	}
	kv.RwLock[path].RLock()
	defer kv.RwLock[path].RUnlock()
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if data[key] == nil {
		return &pb.GetResponse{
			Ok:    false,
			Value: "",
		}, nil
	} else {
		return &pb.GetResponse{
			Ok:    true,
			Value: data[key].(string),
		}, nil
	}
}

func (kv *KvOp) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	fmt.Println("entering getAll")
	defer fmt.Println("exiting getAll")
	// never deadlock
	for _, lock := range kv.RwLock {
		lock.RLock()
	}
	defer func() {
		for _, lock := range kv.RwLock {
			lock.RUnlock()
		}
	}()
	paths, err := utils.ReadAllFiles(kv.DataDir)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var kvs []*pb.GetAllResponse_Kvs
	for _, path := range paths {
		data := map[string]interface{}{}
		err = utils.ReadMap(path, &data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for k, v := range data {
			kvs = append(kvs, &pb.GetAllResponse_Kvs{
				Key:   k,
				Value: v.(string),
			})
		}
	}
	return &pb.GetAllResponse{
		Kvs: kvs,
	}, nil
}

func (kv *KvOp) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	fmt.Println("entering put")
	defer fmt.Println("exiting put")
	key := request.Key
	val := request.Value
	err, path := utils.GetPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	created := true

	data := map[string]interface{}{}
	if kv.RwLock[path] == nil {
		kv.RwLock[path] = new(sync.RWMutex)
	}
	kv.RwLock[path].Lock()
	defer kv.RwLock[path].Unlock()
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if data[key] != nil {
		if data[key] == val {
			return &pb.PutResponse{
				Created: false,
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
		Created: created,
	}, nil
}

func (kv *KvOp) Del(ctx context.Context, request *pb.DelRequest) (*pb.DelResponse, error) {
	fmt.Println("entering del")
	defer fmt.Println("exiting del")
	key := request.Key
	err, path := utils.GetPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	if kv.RwLock[path] == nil {
		kv.RwLock[path] = new(sync.RWMutex)
	}
	kv.RwLock[path].Lock()
	defer kv.RwLock[path].Unlock()
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if data[key] == nil {
		return &pb.DelResponse{
			Deleted: false,
		}, nil
	} else {
		delete(data, key)
		err = utils.WriteMap(path, data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &pb.DelResponse{
			Deleted: true,
		}, nil
	}
}

func (kv *KvOp) MoveData(ctx context.Context, request *pb.MoveDataRequest) (*pb.MoveDataResponse, error) {
	fmt.Println("entering move")
	defer fmt.Println("exiting move")
	for _, lock := range kv.RwLock {
		lock.Lock()
	}
	defer func() {
		for _, lock := range kv.RwLock {
			lock.Unlock()
		}
	}()
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
			if utils.ShouldBeMoved(k, request.ToLabel, request.FromLabel) {
				kvs = append(kvs, &pb.MoveDataResponse_Kv{
					Key:   k,
					Value: v.(string),
				})
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
		Kvs: kvs,
	}, nil
}
