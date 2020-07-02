package client

import (
	"context"
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"google.golang.org/grpc"
	"time"
)

type KvCli struct {
	Mc pb.MasterClient
}

/**
 * Connect to the master node.
 * @param addrList: multiple master address
 * @return kvClient: client
 */
func Connect(addrList []string) *KvCli {
	if len(addrList) == 0 {
		return nil
	}
	index := 0
	var kvCli KvCli
	var masterCli pb.MasterClient
	// Polling to get the available master.
	for ; ; index = (index + 1) % len(addrList) {
		conn, err := grpc.Dial(addrList[index], grpc.WithInsecure())
		if err == nil {
			masterCli = pb.NewMasterClient(conn)
			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			_, err := masterCli.GetSlave(ctx, &pb.GetSlaveRequest{
				Key: "",
			})
			if err == nil || err.Error() == "rpc error: code = Unknown desc = fail to get data node address"{
				kvCli.Mc = masterCli
				break
			}
		}
	}
	// Listening in case of master failures.
	go func() {
		connect := true
		for {
			conn, err := grpc.Dial(addrList[index], grpc.WithInsecure())
			if err == nil {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				masterCli = pb.NewMasterClient(conn)
				_, err = masterCli.GetSlave(ctx, &pb.GetSlaveRequest{
					Key: "",
				})
				if err == nil || err.Error() == "rpc error: code = Unknown desc = fail to get data node address" {
					// Everything is normal.
					// Avoid frequent pinging.
					if connect {
						time.Sleep(300 * time.Millisecond)
					} else {
						// Get the new client.
						masterCli = pb.NewMasterClient(conn)
						kvCli.Mc = masterCli
						connect = true
					}
					continue
				}
			}
			// Master failure.
			connect = false
			index = (index + 1) % len(addrList)
		}
	}()
	return &kvCli
}

/**
 * Put operation.
 * @param key
 * @param value
 * @return err: error
 * @return created: created or updated
 */
func (cli *KvCli) Put(key, value string) (error, bool) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	// Get slave address.
	rsp, err := cli.Mc.GetSlave(ctx, &pb.GetSlaveRequest{
		Key: key,
	})
	if err != nil {
		return err, false
	}
	conn, err := grpc.DialContext(ctx, rsp.Addr, grpc.WithInsecure())
	if err != nil {
		return err, false
	}
	kvClient := pb.NewDataClient(conn)
	// Put data.
	putRsp, err := kvClient.Put(ctx, &pb.PutRequest{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return err, false
	}
	return nil, putRsp.Created
}

/**
 * Get operation.
 * @param key
 * @return err: error
 * @return value: corresponding value of the key
 */
func (cli *KvCli) Get(key string) (error, string) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	// Get slave address.
	rsp, err := cli.Mc.GetSlave(ctx, &pb.GetSlaveRequest{
		Key: key,
	})
	if err != nil {
		return err, ""
	}
	conn, err := grpc.DialContext(ctx, rsp.Addr, grpc.WithInsecure())
	if err != nil {
		return err, ""
	}
	kvClient := pb.NewDataClient(conn)
	// Get data.
	getRsp, err := kvClient.Get(ctx, &pb.GetRequest{
		Key: key,
	})
	if err != nil {
		return err, ""
	}
	if !getRsp.Ok {
		return fmt.Errorf("not found"), ""
	}
	return nil, getRsp.Value
}

/**
 * Delete operation.
 * @param key
 * @return err: error
 * @return deleted: deleted or not found
 */
func (cli *KvCli) Del(key string) (error, bool) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	// Get slave address.
	rsp, err := cli.Mc.GetSlave(ctx, &pb.GetSlaveRequest{
		Key: key,
	})
	if err != nil {
		return err, false
	}
	conn, err := grpc.DialContext(ctx, rsp.Addr, grpc.WithInsecure())
	if err != nil {
		return err, false
	}
	kvClient := pb.NewDataClient(conn)
	// Delete data.
	delRsp, err := kvClient.Del(ctx, &pb.DelRequest{
		Key: key,
	})
	if err != nil {
		return err, false
	}
	return nil, delRsp.Deleted
}

/**
 * Get all the data storage information.
 * @return err: error
 * @return response: response data
 */
func (cli *KvCli) DumpAll() (error, []*pb.DumpAllResponse_Data) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// Dump all.
	rsp, err := cli.Mc.DumpAll(ctx, &pb.DumpAllRequest{
	})
	if err != nil {
		return err, nil
	}
	return nil, rsp.Data
}
