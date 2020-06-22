package rpc

import (
	"context"
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"google.golang.org/grpc"
)

/*
 *  Master has to take good care of label.
 */
type SlaveMeta struct {
	Label int
	Host  string
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
	hash := int(utils.BasicHash(request.Key) % utils.HoleNum)
	dist := utils.HoleNum
	var addr string
	var lastLabel int
	for _, v := range m.SlaveList {
		// if the same dist, choose the one with the smallest label
		if utils.Dist(v.Label, hash, utils.HoleNum) < dist || utils.Dist(v.Label, hash, utils.HoleNum) == dist && lastLabel > v.Label {
			dist = utils.Dist(v.Label, hash, utils.HoleNum)
			lastLabel = v.Label
			addr = v.Host
		}
	}
	if addr == "" {
		return nil, fmt.Errorf("fail to get data node address")
	}
	return &pb.GetSlaveResponse{
		Addr: addr,
	}, nil
}

func (m Master) DumpAll(ctx context.Context, request *pb.DumpAllRequest) (*pb.DumpAllResponse, error) {
	var data []*pb.DumpAllResponse_Data
	for _, slave := range m.SlaveList {
		conn, err := grpc.Dial(slave.Host, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}

		kvClient := pb.NewKvClient(conn)
		rsp, err := kvClient.GetAll(ctx, &pb.GetAllRequest{
		})
		if err != nil {
			return nil, err
		}
		var kvls []*pb.DumpAllResponse_Data_Kvls
		for _, v := range rsp.Kvs {
			kvls = append(kvls, &pb.DumpAllResponse_Data_Kvls{
				Key:   v.Key,
				Value: v.Value,
				Label: int32(utils.BasicHash(v.Key) % utils.HoleNum),
			})
		}
		data = append(data, &pb.DumpAllResponse_Data{
			Host:  slave.Host,
			Label: int32(slave.Label),
			Kvls:  kvls,
		})
	}

	return &pb.DumpAllResponse{
		Data: data,
	}, nil
}