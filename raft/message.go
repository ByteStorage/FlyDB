package cluster

import (
	"context"
	"fmt"
	"github.com/ByteStorage/FlyDB/lib/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"reflect"
)

func SendMessage[T any](addr string, msg T, msgType interface{}) (interface{}, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	msgValue := reflect.ValueOf(msg)
	msgTypeValue := reflect.TypeOf(msgType)

	if msgValue.Type().AssignableTo(msgTypeValue.Elem()) {
		switch msgTypeValue {
		case reflect.TypeOf((*proto.SlaveGetRequest)(nil)).Elem():
			client := proto.NewSlaveGrpcServiceClient(conn)
			response, err := client.Get(context.Background(), msgValue.Interface().(*proto.SlaveGetRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		case reflect.TypeOf((*proto.SlaveSetRequest)(nil)).Elem():
			client := proto.NewSlaveGrpcServiceClient(conn)
			response, err := client.Set(context.Background(), msgValue.Interface().(*proto.SlaveSetRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		case reflect.TypeOf((*proto.SlaveDelRequest)(nil)).Elem():
			client := proto.NewSlaveGrpcServiceClient(conn)
			response, err := client.Del(context.Background(), msgValue.Interface().(*proto.SlaveDelRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		case reflect.TypeOf((*proto.SlaveKeysRequest)(nil)).Elem():
			client := proto.NewSlaveGrpcServiceClient(conn)
			response, err := client.Keys(context.Background(), msgValue.Interface().(*proto.SlaveKeysRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		case reflect.TypeOf((*proto.MasterGetRequest)(nil)).Elem():
			client := proto.NewMasterGrpcServiceClient(conn)
			response, err := client.Get(context.Background(), msgValue.Interface().(*proto.MasterGetRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		case reflect.TypeOf((*proto.MasterSetRequest)(nil)).Elem():
			client := proto.NewMasterGrpcServiceClient(conn)
			response, err := client.Set(context.Background(), msgValue.Interface().(*proto.MasterSetRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		case reflect.TypeOf((*proto.MasterDelRequest)(nil)).Elem():
			client := proto.NewMasterGrpcServiceClient(conn)
			response, err := client.Del(context.Background(), msgValue.Interface().(*proto.MasterDelRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		case reflect.TypeOf((*proto.MasterKeysRequest)(nil)).Elem():
			client := proto.NewMasterGrpcServiceClient(conn)
			response, err := client.Keys(context.Background(), msgValue.Interface().(*proto.MasterKeysRequest))
			if err != nil {
				return nil, err
			}
			return response, nil
		default:
			return nil, fmt.Errorf("unsupported message type")
		}
	} else {
		return nil, fmt.Errorf("msg is not assignable to the specified type")
	}
}
