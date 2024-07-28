package client

import (
	"context"
	"errors"

	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
)

type Range struct {
	Start int32
	End   int32
}

type Incr struct {
	Member string
	Inc    int32
}

func (c *Client) ZAdd(key string, member *gzset.ZSetValue) (*gzset.ZAddResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZAddRequest{Key: key, Member: member}
		response *gzset.ZAddResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZAdd(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZAdds(key string, member []*gzset.ZSetValue) (*gzset.ZAddsResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZAddsRequest{Key: key, Members: member}
		response *gzset.ZAddsResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZAdds(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZRem(key string, member string) (*gzset.ZRemResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZRemRequest{Key: key, Member: member}
		response *gzset.ZRemResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZRem(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZRems(key string, members []string) (*gzset.ZRemsResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZRemsRequest{Key: key, Members: members}
		response *gzset.ZRemsResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZRems(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZScore(key string, member string) (*gzset.ZScoreResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZScoreRequest{Key: key, Member: member}
		response *gzset.ZScoreResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZScore(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZRank(key string, member string) (*gzset.ZRankResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZRankRequest{Key: key, Member: member}
		response *gzset.ZRankResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZRank(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZRevRank(key string, member string) (*gzset.ZRevRankResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZRevRankRequest{Key: key, Member: member}
		response *gzset.ZRevRankResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZRevRank(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZRange(key string, requestRange *Range) (*gzset.ZRangeResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZRangeRequest{Key: key, Start: requestRange.Start, End: requestRange.End}
		response *gzset.ZRangeResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZRange(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZCount(key string, requestRange *Range) (*gzset.ZCountResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZCountRequest{Key: key, Min: requestRange.Start, Max: requestRange.End}
		response *gzset.ZCountResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZCount(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZRevRange(key string, requestRange *Range) (*gzset.ZRevRangeResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZRevRangeRequest{Key: key, StartRank: requestRange.Start, EndRank: requestRange.End}
		response *gzset.ZRevRangeResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZRevRange(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZCard(key string) (*gzset.ZCardResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZCardRequest{Key: key}
		response *gzset.ZCardResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZCard(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) ZIncrBy(key string, incr *Incr) (*gzset.ZIncrByResponse, error) {
	var (
		client   gzset.GZSetServiceClient
		request  = &gzset.ZIncrByRequest{Key: key, Member: incr.Member, IncBy: incr.Inc}
		response *gzset.ZIncrByResponse
		err      error
	)

	if client, err = c.newZSetGrpcClient(); err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}

	if response, err = client.ZIncrBy(context.Background(), request); err != nil {
		return nil, err
	}

	return response, nil
}
