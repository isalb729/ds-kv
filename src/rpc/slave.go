package rpc

import (
	"context"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"google.golang.org/grpc"
	"log"
	"sync"
)

type KvOp struct {
	DataDir    string
	StoreLevel int
	// Read write lock for each file.
	RwLock  map[string]*sync.RWMutex
	// Standby node address.
	// Send backup data.
	Sb []string
	GLock sync.Mutex
}

type Sb struct {
	DataDir string
	StoreLevel int
	// Lock for each file.
	Lock  map[string]*sync.Mutex
	// Avoid concurrent map writing.
	GLock sync.Mutex
}

/**
 * Standby node put.
 */
func (sb *Sb) Put(ctx context.Context, request *pb.PutRequest) (*pb.NoResponse, error) {
	key := request.Key
	val := request.Value
	err, path := utils.GetPath(sb.DataDir, key, sb.StoreLevel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data := map[string]interface{}{}
	// Avoid concurrent map writing.
	sb.GLock.Lock()
	if sb.Lock[path] == nil {
		sb.Lock[path] = new(sync.Mutex)
	}
	sb.GLock.Unlock()
	sb.Lock[path].Lock()
	defer sb.Lock[path].Unlock()
	// Read data and append.
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = utils.AppendMap(path, map[string]interface{}{key: val})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.NoResponse{
	}, nil
}

/**
 * Delete standby data.
 */
func (sb *Sb) Del(ctx context.Context, request *pb.DelRequest) (*pb.NoResponse, error) {
	key := request.Key
	err, path := utils.GetPath(sb.DataDir, key, sb.StoreLevel)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	sb.GLock.Lock()
	if sb.Lock[path] == nil {
		sb.Lock[path] = new(sync.Mutex)
	}
	sb.GLock.Unlock()
	sb.Lock[path].Lock()
	defer sb.Lock[path].Unlock()
	// Read and rewrite.
	err = utils.ReadMap(path, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if data[key] == nil {
		return &pb.NoResponse{
		}, nil
	} else {
		delete(data, key)
		err = utils.WriteMap(path, data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &pb.NoResponse{
		}, nil
	}
}

/**
 * Get data.
 */
func (kv *KvOp) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	key := request.Key
	err, path := utils.GetPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data := map[string]interface{}{}
	kv.GLock.Lock()
	if kv.RwLock[path] == nil {
		kv.RwLock[path] = new(sync.RWMutex)
	}
	kv.GLock.Unlock()
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

/**
 * Get all data.
 * For debugging.
 */
func (kv *KvOp) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	// never deadlock
	for _, lock := range kv.RwLock {
		lock.RLock()
	}
	defer func() {
		for _, lock := range kv.RwLock {
			lock.RUnlock()
		}
	}()
	// Read all path.
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

/**
 * Put data.
 */
func (kv *KvOp) Put(ctx context.Context, request *pb.PutRequest) (*pb.PutResponse, error) {
	key := request.Key
	val := request.Value
	err, path := utils.GetPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	created := true

	data := map[string]interface{}{}
	kv.GLock.Lock()
	if kv.RwLock[path] == nil {
		kv.RwLock[path] = new(sync.RWMutex)
	}
	kv.GLock.Unlock()
	kv.RwLock[path].Lock()
	defer kv.RwLock[path].Unlock()
	done := make(chan bool)
	defer func() {
		<-done
	}()
	// Send to standby.
	go func() {
		for _, sb := range kv.Sb {
			sbConn, err := grpc.Dial(sb, grpc.WithInsecure())
			if err != nil {
				log.Println("Put to", sb, ":", err)
				continue
			}
			sbClient := pb.NewDataStandByClient(sbConn)
			_, err = sbClient.Put(context.Background(), &pb.PutRequest{
				Key:                  key,
				Value:                val,
			})
			if err != nil {
				log.Println("Put to", sb, ":", err)
			}
		}
		done <- true
	}()
	// Read and append.
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

/**
 * Delete data.
 */
func (kv *KvOp) Del(ctx context.Context, request *pb.DelRequest) (*pb.DelResponse, error) {
	key := request.Key
	err, path := utils.GetPath(kv.DataDir, key, kv.StoreLevel)
	if err != nil {
		return nil, err
	}
	data := map[string]interface{}{}
	kv.GLock.Lock()
	if kv.RwLock[path] == nil {
		kv.RwLock[path] = new(sync.RWMutex)
	}
	kv.GLock.Unlock()
	kv.RwLock[path].Lock()
	defer kv.RwLock[path].Unlock()
	done := make(chan bool)
	defer func() {
		<-done
	}()
	go func() {
		for _, sb := range kv.Sb {
			sbConn, err := grpc.Dial(sb, grpc.WithInsecure())
			if err != nil {
				log.Println("Del to", sb, ":", err)
				continue
			}
			sbClient := pb.NewDataStandByClient(sbConn)
			_, err = sbClient.Del(context.Background(), &pb.DelRequest{
				Key:                  key,
			})
			if err != nil {
				log.Println("Del to", sb, ":", err)
			}
		}
		done <- true
	}()
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

/**
 * Move data.
 * For registration or deregistration of data node.
 */
func (kv *KvOp) MoveData(ctx context.Context, request *pb.MoveDataRequest) (*pb.MoveDataResponse, error) {
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
	// Find the one that's closer to the adjcent.
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
