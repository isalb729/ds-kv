package client

import (
	"context"
	"github.com/isalb729/ds-kv/src/rpc"
	"google.golang.org/grpc"
	"time"
)

type KvCli struct {
	Mc rpc.MasterClient
}

/**
 * @return err: error
 * @return created: created or updated
 */
func (cli *KvCli) put(key, value string) (error, bool) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	rsp, err := cli.Mc.GetSlave(ctx, &rpc.GetSlaveRequest{
		Key: key,
	})
	if err != nil {
		return err, false
	}
	conn, err := grpc.DialContext(ctx, rsp.Addr)
	if err != nil {
		return err, false
	}
	kvClient := rpc.NewKvClient(conn)
	putRsp, err := kvClient.Put(ctx, &rpc.PutRequest{
		Key:                  key,
		Value:                value,
	})
	if err != nil {
		return err, false
	}
	return nil, putRsp.Created
}

/**
 * @return err: error
 * @return value: corresponding value of the key
 */
func (cli *KvCli) get(key string) (err error, value string) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	rsp, err := cli.Mc.GetSlave(ctx, &rpc.GetSlaveRequest{
		Key: key,
	})
	if err != nil {
		return err, ""
	}
	conn, err := grpc.DialContext(ctx, rsp.Addr)
	if err != nil {
		return err, ""
	}
	kvClient := rpc.NewKvClient(conn)
	getRsp, err := kvClient.Get(ctx, &rpc.GetRequest{
		Key:                  key,
	})
	if err != nil {
		return err, ""
	}
	return nil, getRsp.Value
}

/**
 * @return err: error
 * @return deleted: deleted or not found
 */
func (cli *KvCli) del(key string) (error, bool) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	rsp, err := cli.Mc.GetSlave(ctx, &rpc.GetSlaveRequest{
		Key: key,
	})
	if err != nil {
		return err, false
	}
	conn, err := grpc.DialContext(ctx, rsp.Addr)
	if err != nil {
		return err, false
	}
	kvClient := rpc.NewKvClient(conn)
	delRsp, err := kvClient.Del(ctx, &rpc.DelRequest{
		Key:                  key,
	})
	if err != nil {
		return err, false
	}
	return nil, delRsp.Deleted
}
