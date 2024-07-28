package service

import (
	"context"
	"errors"
	"math"

	pbany "github.com/golang/protobuf/ptypes/any"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"github.com/ByteStorage/FlyDB/lib/proto/gzset"
	"github.com/ByteStorage/FlyDB/structure"
)

type ZSetService interface {
	Base
	gzset.GZSetServiceServer
}

var _ ZSetService = &zSet{}

type zSet struct {
	dbs *structure.ZSetStructure
	gzset.GZSetServiceServer
}

func NewZSetService(options config.Options) (ZSetService, error) {
	zSetStructure, err := structure.NewZSetStructure(options)
	if err != nil {
		return nil, err
	}
	return &zSet{dbs: zSetStructure}, nil
}

func (z *zSet) CloseDb() error {
	return z.dbs.Stop()
}

func (z *zSet) checkScoreIntIsInt32(score int) bool {
	return score >= math.MinInt32 && score <= math.MaxInt32
}

func (z *zSet) ZAdd(_ context.Context, request *gzset.ZAddRequest) (*gzset.ZAddResponse, error) {
	var (
		err error
	)

	if err = z.dbs.ZAdd(request.Key, int(request.Member.Score), request.Member.Member,
		string(request.Member.Value.Value)); err != nil {
		return &gzset.ZAddResponse{Success: false}, err
	}

	return &gzset.ZAddResponse{Success: true}, nil
}

func (z *zSet) ZAdds(_ context.Context, request *gzset.ZAddsRequest) (*gzset.ZAddsResponse, error) {
	for _, member := range request.Members {
		var (
			err error
		)

		if !z.checkScoreIntIsInt32(int(member.Score)) {
			return &gzset.ZAddsResponse{Success: false}, nil
		}

		if err = z.dbs.ZAdd(request.Key, int(member.Score), member.Member, string(member.Value.GetValue())); err != nil {
			return &gzset.ZAddsResponse{Success: false}, err
		}

	}

	return &gzset.ZAddsResponse{Success: true}, nil
}

func (z *zSet) ZRem(_ context.Context, request *gzset.ZRemRequest) (*gzset.ZRemResponse, error) {
	if err := z.dbs.ZRem(request.Key, request.Member); err != nil {
		return &gzset.ZRemResponse{Success: false}, err
	}

	return &gzset.ZRemResponse{Success: true}, nil
}

func (z *zSet) ZRems(_ context.Context, request *gzset.ZRemsRequest) (*gzset.ZRemsResponse, error) {
	if err := z.dbs.ZRems(request.Key, request.Members...); err != nil {
		return &gzset.ZRemsResponse{Success: false}, err
	}

	return &gzset.ZRemsResponse{Success: true}, nil
}

func (z *zSet) ZScore(_ context.Context, request *gzset.ZScoreRequest) (*gzset.ZScoreResponse, error) {
	var (
		err   error
		score int
	)

	if score, err = z.dbs.ZScore(request.Key, request.Member); err != nil {
		return &gzset.ZScoreResponse{}, err
	}

	return &gzset.ZScoreResponse{Score: int32(score), Exists: true}, nil
}

func (z *zSet) ZRank(_ context.Context, request *gzset.ZRankRequest) (*gzset.ZRankResponse, error) {
	var (
		err  error
		rank int
	)

	if rank, err = z.dbs.ZRank(request.Key, request.Member); err != nil {
		return &gzset.ZRankResponse{}, err
	}

	return &gzset.ZRankResponse{Rank: int32(rank), Exists: true}, nil
}

func (z *zSet) ZRevRank(_ context.Context, request *gzset.ZRevRankRequest) (*gzset.ZRevRankResponse, error) {
	var (
		err  error
		rank int
	)

	if rank, err = z.dbs.ZRevRank(request.Key, request.Member); err != nil {
		return &gzset.ZRevRankResponse{}, err
	}

	return &gzset.ZRevRankResponse{Rank: int32(rank), Exists: true}, nil
}

func (z *zSet) structureZSetValue2GZSetValue(structureZSetValue structure.ZSetValue) (*gzset.ZSetValue, error) {
	var (
		structureValueBytes []byte
		err                 error
		gzSetValue          = &gzset.ZSetValue{}
		value               []byte
		dec                 *encoding.MessagePackCodecDecoder
	)

	structureValueBytes, _ = structureZSetValue.MarshalBinary()

	dec = encoding.NewMessagePackDecoder(structureValueBytes)

	if err = dec.Decode(&gzSetValue.Member); err != nil {
		return gzSetValue, err
	}

	if err = dec.Decode(&gzSetValue.Score); err != nil {
		return gzSetValue, err
	}

	err = dec.Decode(&value)
	gzSetValue.Value = &pbany.Any{Value: value}
	return gzSetValue, err
}

func (z *zSet) structureZSetValueList2GZSetValueList(structureZSetValueList []structure.ZSetValue) ([]*gzset.
	ZSetValue, error) {
	var (
		err             error
		responseMembers = make([]*gzset.ZSetValue, 0, len(structureZSetValueList))
	)

	for _, member := range structureZSetValueList {
		var (
			gzSetValue *gzset.ZSetValue
		)
		if gzSetValue, err = z.structureZSetValue2GZSetValue(member); err != nil {
			return make([]*gzset.ZSetValue, 0), err
		}
		responseMembers = append(responseMembers, gzSetValue)
	}

	return responseMembers, nil
}

func (z *zSet) ZRange(_ context.Context, request *gzset.ZRangeRequest) (*gzset.ZRangeResponse, error) {
	var (
		err             error
		members         []structure.ZSetValue
		responseMembers = make([]*gzset.ZSetValue, 0, len(members))
	)

	if members, err = z.dbs.ZRange(request.Key, int(request.Start), int(request.End)); err != nil {
		return &gzset.ZRangeResponse{}, err
	}

	if responseMembers, err = z.structureZSetValueList2GZSetValueList(members); err != nil {
		return &gzset.ZRangeResponse{}, err
	}

	return &gzset.ZRangeResponse{Members: responseMembers}, nil
}

func (z *zSet) ZCount(_ context.Context, request *gzset.ZCountRequest) (*gzset.ZCountResponse, error) {
	var (
		err   error
		count int
	)

	if count, err = z.dbs.ZCount(request.Key, int(request.Min), int(request.Max)); err != nil {
		return &gzset.ZCountResponse{}, err
	}

	return &gzset.ZCountResponse{Count: int32(count)}, nil
}

func (z *zSet) ZRevRange(_ context.Context, request *gzset.ZRevRangeRequest) (*gzset.ZRevRangeResponse, error) {
	var (
		err             error
		members         []structure.ZSetValue
		responseMembers = make([]*gzset.ZSetValue, 0, len(members))
	)

	if members, err = z.dbs.ZRevRange(request.Key, int(request.StartRank), int(request.EndRank)); err != nil {
		return &gzset.ZRevRangeResponse{}, err
	}

	if responseMembers, err = z.structureZSetValueList2GZSetValueList(members); err != nil {
		return &gzset.ZRevRangeResponse{}, err
	}

	return &gzset.ZRevRangeResponse{Members: responseMembers}, nil
}

func (z *zSet) ZCard(_ context.Context, request *gzset.ZCardRequest) (*gzset.ZCardResponse, error) {
	var (
		err   error
		count int
	)

	if count, err = z.dbs.ZCard(request.Key); err != nil {
		return &gzset.ZCardResponse{}, err
	}

	return &gzset.ZCardResponse{Count: int32(count)}, nil
}

func (z *zSet) ZIncrBy(_ context.Context, request *gzset.ZIncrByRequest) (*gzset.ZIncrByResponse, error) {
	var (
		err      error
		newScore int
	)

	if err = z.dbs.ZIncrBy(request.Key, request.Member, int(request.IncBy)); err != nil {
		return &gzset.ZIncrByResponse{}, err
	}

	if newScore, err = z.dbs.ZScore(request.Key, request.Member); err != nil {
		return &gzset.ZIncrByResponse{}, err
	}

	if !z.checkScoreIntIsInt32(newScore) {
		return &gzset.ZIncrByResponse{NewScore: 0, Exists: true}, errors.New("new score is not int32")
	}

	return &gzset.ZIncrByResponse{NewScore: int32(newScore), Exists: true}, nil
}
