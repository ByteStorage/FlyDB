package client

import (
	"context"
	"errors"
	"github.com/ByteStorage/FlyDB/lib/proto/gset"
)

func (c *Client) SAdd(key, member string) error {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SAddRequest{Key: key, Member: member}
	sadd, err := client.SAdd(context.Background(), req)
	if err != nil {
		return err
	}
	if !sadd.OK {
		return errors.New("SAdd failed")
	}
	return nil
}

func (c *Client) SAdds(key string, members []string) error {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SAddsRequest{Key: key, Members: members}
	srems, err := client.SAdds(context.Background(), req)
	if err != nil {
		return err
	}
	if !srems.OK {
		return errors.New("SAdds failed")
	}
	return nil
}

func (c *Client) SRem(key, member string) error {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SRemRequest{Key: key, Member: member}
	srem, err := client.SRem(context.Background(), req)
	if err != nil {
		return err
	}
	if !srem.OK {
		return errors.New("SRem failed")
	}
	return nil
}

func (c *Client) SRems(key string, members []string) error {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SRemsRequest{Key: key, Members: members}
	srems, err := client.SRems(context.Background(), req)
	if err != nil {
		return err
	}
	if !srems.OK {
		return errors.New("SRems failed")
	}
	return nil
}

func (c *Client) SCard(key string) (int32, error) {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return 0, errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SCardRequest{Key: key}
	scard, err := client.SCard(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return scard.Count, nil
}

func (c *Client) SMembers(key string) ([]string, error) {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SMembersRequest{Key: key}
	smembers, err := client.SMembers(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return smembers.Members, nil
}

func (c *Client) SIsMember(key, member string) (bool, error) {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return false, errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SIsMemberRequest{Key: key, Member: member}
	sismember, err := client.SIsMember(context.Background(), req)
	if err != nil {
		return false, err
	}
	return sismember.IsMember, nil
}

func (c *Client) SUnion(keys []string) ([]string, error) {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SUnionRequest{Keys: keys}
	sunion, err := client.SUnion(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return sunion.Members, nil
}

func (c *Client) SInter(keys []string) ([]string, error) {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SInterRequest{Keys: keys}
	sinter, err := client.SInter(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return sinter.Members, nil
}

func (c *Client) SDiff(keys []string) ([]string, error) {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return nil, errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SDiffRequest{Keys: keys}
	sdiff, err := client.SDiff(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return sdiff.Members, nil
}

func (c *Client) SUnionStore(destination string, keys []string) error {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SUnionStoreRequest{Destination: destination, Keys: keys}
	sunionstore, err := client.SUnionStore(context.Background(), req)
	if err != nil {
		return err
	}
	if !sunionstore.OK {
		return errors.New("SUnionStore failed")
	}
	return nil
}

func (c *Client) SInterStore(destinationKey string, keys []string) error {
	client, err := newSetGrpcClient(c.Addr)
	if err != nil {
		return errors.New("new grpc client error: " + err.Error())
	}
	req := &gset.SInterStoreRequest{Destination: destinationKey, Keys: keys}
	sinterstore, err := client.SInterStore(context.Background(), req)
	if err != nil {
		return err
	}
	if !sinterstore.OK {
		return errors.New("SInterStore failed")
	}
	return nil
}
