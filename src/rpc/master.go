package rpc

import (
	"context"
	"fmt"
	"github.com/isalb729/ds-kv/src/rpc/pb"
	"github.com/isalb729/ds-kv/src/utils"
	"google.golang.org/grpc"
	"log"
)

// Master has to take good care of label.
type SlaveMeta struct {
	Label int
	Host  string
}

type Master struct {
	SlaveList []SlaveMeta
	// Available slaves.
	// Used to detect failures.
	Working   map[string]bool
}

/**
 * Label as deregistered.
 * If znode disappear and no notification, then it's a failure.
 */
func (m *Master) DeregisterNotify(ctx context.Context, request *pb.DeregisterNotifyRequest) (*pb.DeregisterNotifyResponse, error) {
	log.Println("Data server", request.Addr, "deregistered")
	m.Working[request.Addr] = false
	return &pb.DeregisterNotifyResponse{
	}, nil
}

/**
 * Client to master.
 * Send key to get the host of corresponding server.
 */
func (m Master) GetSlave(ctx context.Context, request *pb.GetSlaveRequest) (*pb.GetSlaveResponse, error) {
	// Consistent hash and return key.
	hash := int(utils.BasicHash(request.Key) % utils.HoleNum)
	dist := utils.HoleNum
	var addr string
	var lastLabel int
	// Find the closest label.
	for _, v := range m.SlaveList {
		// If the same dist, choose the one with the smallest label.
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

/**
 * For debugging.
 * Dump all data.
 */
func (m Master) DumpAll(ctx context.Context, request *pb.DumpAllRequest) (*pb.DumpAllResponse, error) {
	var data []*pb.DumpAllResponse_Data
	for _, slave := range m.SlaveList {
		conn, err := grpc.Dial(slave.Host, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		kvClient := pb.NewDataClient(conn)
		// Get all.
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
