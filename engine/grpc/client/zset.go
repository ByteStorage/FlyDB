package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
	"github.com/ByteStorage/FlyDB/structure"
)

func (c *Client) ZAdd(key string, score int, member string, value interface{}) error {
	client, err := newZSetGrpcClient(c.Addr)
	req := &gzset.ZAddRequest{}
	switch v := value.(type) {
	case string:
		req = &gzset.ZAddRequest{Key: key, Member: &gzset.ZSetValue{Score: int32(score), Member: member, Value: &gzset.ZSetValue_StringValue{StringValue: v}}}
	default:
		return fmt.Errorf("unsupported value type")
	}
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	sadd, err := client.ZAdd(context.Background(), req)
	if err != nil {
		return err
	}
	if !sadd.Success {
		return errors.New("SAdd failed")
	}
	return nil
}

func (c *Client) ZAdds(key string, vals ...structure.ZSetValue) error {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}

	zmembers := make([]*gzset.ZSetValue, len(vals))
	for i, member := range vals {
		switch v := member.Value.(type) {
		case string:
			zmembers[i] = &gzset.ZSetValue{Score: int32(member.Score), Member: member.Member, Value: &gzset.ZSetValue_StringValue{StringValue: v}}
		case int32:
			zmembers[i] = &gzset.ZSetValue{Score: int32(member.Score), Member: member.Member, Value: &gzset.ZSetValue_Int32Value{Int32Value: v}}
		case int64:
			zmembers[i] = &gzset.ZSetValue{Score: int32(member.Score), Member: member.Member, Value: &gzset.ZSetValue_Int64Value{Int64Value: v}}
		case float32:
			zmembers[i] = &gzset.ZSetValue{Score: int32(member.Score), Member: member.Member, Value: &gzset.ZSetValue_Float32Value{Float32Value: v}}
		case float64:
			zmembers[i] = &gzset.ZSetValue{Score: int32(member.Score), Member: member.Member, Value: &gzset.ZSetValue_Float64Value{Float64Value: v}}
		case bool:
			zmembers[i] = &gzset.ZSetValue{Score: int32(member.Score), Member: member.Member, Value: &gzset.ZSetValue_BoolValue{BoolValue: v}}
		case []byte:
			zmembers[i] = &gzset.ZSetValue{Score: int32(member.Score), Member: member.Member, Value: &gzset.ZSetValue_BytesValue{BytesValue: v}}
		default:
			return errors.New("unknown value type")
		}
	}
	req := &gzset.ZAddsRequest{Key: key, Members: zmembers}
	sadd, err := client.ZAdds(context.Background(), req)
	if err != nil {
		return err
	}
	if !sadd.Success {
		return errors.New("SAdd failed")
	}
	return nil
}

func (c *Client) ZRem(key, member string) error {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZRemRequest{Key: key, Member: member}
	resp, err := client.ZRem(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New("ZRem failed")
	}
	return nil
}

func (c *Client) ZRems(key string, members []string) error {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZRemsRequest{Key: key, Members: members}
	resp, err := client.ZRems(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New("ZRem failed")
	}
	return nil
}

func (c *Client) ZScore(key, member string) (int32, error) {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return 0, errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZScoreRequest{Key: key, Member: member}
	score, err := client.ZScore(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return score.Score, nil
}

func (c *Client) ZRank(key, member string) (int32, error) {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return 0, errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZRankRequest{Key: key, Member: member}
	rank, err := client.ZRank(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return rank.Rank, nil
}

func (c *Client) ZRevRank(key, member string) (int32, error) {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return 0, errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZRevRankRequest{Key: key, Member: member}
	rank, err := client.ZRevRank(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return rank.Rank, nil
}

func (c *Client) ZRange(key string, start, stop int32) ([]*structure.ZSetValue, error) {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZRangeRequest{Key: key, Start: start, End: stop}
	result, err := client.ZRange(context.Background(), req)
	if err != nil {
		return nil, err
	}
	var members []*structure.ZSetValue
	for _, rv := range result.GetMembers() {
		members = append(members, &structure.ZSetValue{
			Score:  int(rv.Score),
			Member: rv.Member,
			Value:  rv.Value,
		})
	}
	return members, nil
}

func (c *Client) ZCount(key string, min, max int32) (int32, error) {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return 0, errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZCountRequest{Key: key, Min: min, Max: max}
	count, err := client.ZCount(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

func (c *Client) ZRevRange(key string, start, stop int32) ([]*structure.ZSetValue, error) {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZRevRangeRequest{Key: key, StartRank: start, EndRank: stop}
	result, err := client.ZRevRange(context.Background(), req)
	if err != nil {
		return nil, err
	}
	var members []*structure.ZSetValue
	for _, rv := range result.GetMembers() {
		members = append(members, &structure.ZSetValue{
			Score:  int(rv.Score),
			Member: rv.Member,
			Value:  rv.Value,
		})
	}
	return members, nil
}

func (c *Client) ZCard(key string) (int32, error) {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return 0, errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZCardRequest{Key: key}
	card, err := client.ZCard(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return card.Count, nil
}

func (c *Client) ZIncrBy(key, member string, increment int32) error {
	client, err := newZSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gzset.ZIncrByRequest{Key: key, Member: member, IncBy: increment}
	resp, err := client.ZIncrBy(context.Background(), req)
	if err != nil {
		return err
	}
	if !resp.Exists {
		return errors.New("ZIncrBy failed")
	}
	return nil
}
