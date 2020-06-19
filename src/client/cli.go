package client

import (
	"context"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"time"
)

//import "github.com/isalb729/ds-kv/src/rpc"

func Concurrent(cli *KvCli) {

}

func Sequential(cli *KvCli) {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	_, _ = cli.Mc.GetSlave(ctx, &pb.GetSlaveRequest{
		Key:                  "s",
	})


	//get()
	//put()
	//del()
}
