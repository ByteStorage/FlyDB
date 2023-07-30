package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
	"github.com/ByteStorage/FlyDB/structure"
)

type ZSetService interface {
	Base
	gzset.GZSetServiceServer
}

type zset struct {
	dbs *structure.ZSetStructure
	gzset.GZSetServiceServer
}

func (z *zset) CloseDb() error {
	return z.dbs.Stop()
}

func NewZSetService(options config.Options) (ZSetService, error) {
	zsetStructure, err := structure.NewZSetStructure(options)
	if err != nil {
		return nil, err
	}
	return &zset{dbs: zsetStructure}, nil
}

func (z *zset) ZAdd(ctx context.Context, req *gzset.ZAddRequest) (*gzset.ZAddResponse, error) {
	zSetValue := req.Member
	var value string
	switch v := zSetValue.Value.(type) {
	case *gzset.ZSetValue_StringValue:
		value = v.StringValue
	default:
		return nil, fmt.Errorf("unsupported value type")
	}
	err := z.dbs.ZAdd(req.Key, int(zSetValue.Score), zSetValue.Member, value)
	if err != nil {
		return &gzset.ZAddResponse{}, err
	}
	return &gzset.ZAddResponse{Success: true}, nil
}

func (z *zset) ZAdds(ctx context.Context, req *gzset.ZAddsRequest) (*gzset.ZAddsResponse, error) {
	members := make([]structure.ZSetValue, len(req.Members))
	for i, member := range req.Members {
		var value interface{}
		switch v := member.Value.(type) {
		case *gzset.ZSetValue_StringValue:
			value = v.StringValue
		case *gzset.ZSetValue_Int32Value:
			value = v.Int32Value
		case *gzset.ZSetValue_Int64Value:
			value = v.Int64Value
		case *gzset.ZSetValue_Float32Value:
			value = v.Float32Value
		case *gzset.ZSetValue_Float64Value:
			value = v.Float64Value
		case *gzset.ZSetValue_BoolValue:
			value = v.BoolValue
		case *gzset.ZSetValue_BytesValue:
			value = v.BytesValue
		default:
			return nil, fmt.Errorf("unsupported value type")
		}
		members[i] = structure.ZSetValue{
			Score:  int(member.Score),
			Member: member.Member,
			Value:  value,
		}
	}
	err := z.dbs.ZAdds(req.Key, members...)
	if err != nil {
		return &gzset.ZAddsResponse{}, err
	}

	return &gzset.ZAddsResponse{Success: true}, nil
}

func (z *zset) ZRem(ctx context.Context, req *gzset.ZRemRequest) (*gzset.ZRemResponse, error) {
	err := z.dbs.ZRem(req.Key, req.Member)
	if err != nil {
		return &gzset.ZRemResponse{}, err
	}
	return &gzset.ZRemResponse{Success: true}, nil
}

func (z *zset) ZRems(ctx context.Context, req *gzset.ZRemsRequest) (*gzset.ZRemsResponse, error) {
	err := z.dbs.ZRems(req.Key, req.Members...)
	if err != nil {
		return &gzset.ZRemsResponse{}, err
	}
	return &gzset.ZRemsResponse{Success: true}, nil
}

func (z *zset) ZScore(ctx context.Context, req *gzset.ZScoreRequest) (*gzset.ZScoreResponse, error) {
	score, err := z.dbs.ZScore(req.Key, req.Member)
	if err != nil {
		return &gzset.ZScoreResponse{}, err
	}
	return &gzset.ZScoreResponse{Score: int32(score)}, nil
}

func (z *zset) ZRank(ctx context.Context, req *gzset.ZRankRequest) (*gzset.ZRankResponse, error) {
	rank, err := z.dbs.ZRank(req.Key, req.Member)
	if err != nil {
		return nil, err
	}
	return &gzset.ZRankResponse{Rank: int32(rank)}, nil
}

func (z *zset) ZRevRank(ctx context.Context, req *gzset.ZRevRankRequest) (*gzset.ZRevRankResponse, error) {
	rank, err := z.dbs.ZRevRank(req.Key, req.Member)
	if err != nil {
		return nil, err
	}
	return &gzset.ZRevRankResponse{Rank: int32(rank)}, nil
}

func (z *zset) ZRange(ctx context.Context, req *gzset.ZRangeRequest) (*gzset.ZRangeResponse, error) {
	rangeValues, err := z.dbs.ZRange(req.Key, int(req.Start), int(req.End))
	if err != nil {
		return nil, err
	}
	var scores []int32
	var members []string
	var values []*gzset.Value
	for _, rv := range rangeValues {
		scores = append(scores, int32(rv.Score))
		members = append(members, rv.Member)
		switch v := rv.Value.(type) {
		case string:
			// Wrap the string value in the protobuf message for strings.
			values = append(values, &gzset.Value{Value: &gzset.Value_StringValue{StringValue: v}})
		case int32:
			values = append(values, &gzset.Value{Value: &gzset.Value_Int32Value{Int32Value: v}})
		case int64:
			values = append(values, &gzset.Value{Value: &gzset.Value_Int64Value{Int64Value: v}})
		case float32:
			values = append(values, &gzset.Value{Value: &gzset.Value_Float32Value{Float32Value: v}})
		case float64:
			values = append(values, &gzset.Value{Value: &gzset.Value_Float64Value{Float64Value: v}})
		case bool:
			values = append(values, &gzset.Value{Value: &gzset.Value_BoolValue{BoolValue: v}})
		case []byte:
			values = append(values, &gzset.Value{Value: &gzset.Value_BytesValue{BytesValue: v}})
		default:
			return nil, errors.New("unknown value type")
		}
	}
	return &gzset.ZRangeResponse{Score: scores, Members: members, Values: values}, nil
}

func (z *zset) ZCount(ctx context.Context, req *gzset.ZCountRequest) (*gzset.ZCountResponse, error) {
	count, err := z.dbs.ZCount(req.Key, int(req.Min), int(req.Max))
	if err != nil {
		return &gzset.ZCountResponse{}, err
	}
	return &gzset.ZCountResponse{Count: int32(count)}, nil
}

func (z *zset) ZRevRange(ctx context.Context, req *gzset.ZRevRangeRequest) (*gzset.ZRevRangeResponse, error) {
	rangeValues, err := z.dbs.ZRevRange(req.Key, int(req.StartRank), int(req.EndRank))
	if err != nil {
		return nil, err
	}
	var scores []int32
	var members []string
	var values []*gzset.Value
	for _, rv := range rangeValues {
		scores = append(scores, int32(rv.Score))
		members = append(members, rv.Member)
		switch v := rv.Value.(type) {
		case string:
			// Wrap the string value in the protobuf message for strings.
			values = append(values, &gzset.Value{Value: &gzset.Value_StringValue{StringValue: v}})
		case int32:
			values = append(values, &gzset.Value{Value: &gzset.Value_Int32Value{Int32Value: v}})
		case int64:
			values = append(values, &gzset.Value{Value: &gzset.Value_Int64Value{Int64Value: v}})
		case float32:
			values = append(values, &gzset.Value{Value: &gzset.Value_Float32Value{Float32Value: v}})
		case float64:
			values = append(values, &gzset.Value{Value: &gzset.Value_Float64Value{Float64Value: v}})
		case bool:
			values = append(values, &gzset.Value{Value: &gzset.Value_BoolValue{BoolValue: v}})
		case []byte:
			values = append(values, &gzset.Value{Value: &gzset.Value_BytesValue{BytesValue: v}})
		default:
			return nil, errors.New("unknown value type")
		}
	}
	return &gzset.ZRevRangeResponse{Score: scores, Members: members, Values: values}, nil
}

func (z *zset) ZCard(ctx context.Context, req *gzset.ZCardRequest) (*gzset.ZCardResponse, error) {
	count, err := z.dbs.ZCard(req.Key)
	if err != nil {
		return &gzset.ZCardResponse{}, err
	}
	return &gzset.ZCardResponse{Count: int32(count)}, nil
}

func (z *zset) ZIncrBy(ctx context.Context, req *gzset.ZIncrByRequest) (*gzset.ZIncrByResponse, error) {
	err := z.dbs.ZIncrBy(req.Key, req.Member, int(req.IncBy))
	if err != nil {
		return nil, err
	}
	return &gzset.ZIncrByResponse{Exists: true}, nil
}
