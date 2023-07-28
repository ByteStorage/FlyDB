package service

import (
	"context"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/proto/gset"
	"github.com/ByteStorage/FlyDB/structure"
)

type SetService interface {
	Base
	gset.GSetServiceServer
}

type set struct {
	dbs *structure.SetStructure
	gset.GSetServiceServer
}

func (l *set) CloseDb() error {
	return l.dbs.Stop()
}

func NewSetService(options config.Options) (SetService, error) {
	setStructure, err := structure.NewSetStructure(options)
	if err != nil {
		return nil, err
	}
	return &set{dbs: setStructure}, nil
}

func (s *set) SAdd(ctx context.Context, req *gset.SAddRequest) (*gset.EmptyResponse, error) {
	err := s.dbs.SAdd(req.Key, req.Member)
	if err != nil {
		return &gset.EmptyResponse{}, err
	}
	return &gset.EmptyResponse{OK: true}, nil
}

func (s *set) SAdds(ctx context.Context, req *gset.SAddsRequest) (*gset.EmptyResponse, error) {
	for _, member := range req.Members {
		err := s.dbs.SAdd(req.Key, member)
		if err != nil {
			return &gset.EmptyResponse{}, err
		}
	}
	return &gset.EmptyResponse{OK: true}, nil
}

func (s *set) SRem(ctx context.Context, req *gset.SRemRequest) (*gset.EmptyResponse, error) {
	err := s.dbs.SRem(req.Key, req.Member)
	if err != nil {
		return &gset.EmptyResponse{}, err
	}
	return &gset.EmptyResponse{OK: true}, nil
}

func (s *set) SRems(ctx context.Context, req *gset.SRemsRequest) (*gset.EmptyResponse, error) {
	err := s.dbs.SRems(req.Key, req.Members...)
	if err != nil {
		return &gset.EmptyResponse{}, err
	}
	return &gset.EmptyResponse{OK: true}, nil
}

func (s *set) SCard(ctx context.Context, req *gset.SCardRequest) (*gset.SCardResponse, error) {
	count, err := s.dbs.SCard(req.Key)
	if err != nil {
		return nil, err
	}
	return &gset.SCardResponse{Count: int32(count)}, nil
}

func (s *set) SMembers(ctx context.Context, req *gset.SMembersRequest) (*gset.SMembersResponse, error) {
	members, err := s.dbs.SMembers(req.Key)
	if err != nil {
		return nil, err
	}
	return &gset.SMembersResponse{Members: members}, nil
}

func (s *set) SIsMember(ctx context.Context, req *gset.SIsMemberRequest) (*gset.SIsMemberResponse, error) {
	isMember, err := s.dbs.SIsMember(req.Key, req.Member)
	if err != nil {
		return nil, err
	}
	return &gset.SIsMemberResponse{IsMember: isMember}, nil
}

func (s *set) SUnion(ctx context.Context, req *gset.SUnionRequest) (*gset.SUnionResponse, error) {
	union, err := s.dbs.SUnion(req.Keys...)
	if err != nil {
		return nil, err
	}
	return &gset.SUnionResponse{Members: union}, nil
}

func (s *set) SInter(ctx context.Context, req *gset.SInterRequest) (*gset.SInterResponse, error) {
	inter, err := s.dbs.SInter(req.Keys...)
	if err != nil {
		return nil, err
	}
	return &gset.SInterResponse{Members: inter}, nil
}

func (s *set) SDiff(ctx context.Context, req *gset.SDiffRequest) (*gset.SDiffResponse, error) {
	diff, err := s.dbs.SDiff(req.Keys...)
	if err != nil {
		return nil, err
	}
	return &gset.SDiffResponse{Members: diff}, nil
}

func (s *set) SUnionStore(ctx context.Context, req *gset.SUnionStoreRequest) (*gset.EmptyResponse, error) {
	err := s.dbs.SUnionStore(req.Destination, req.Keys...)
	if err != nil {
		return &gset.EmptyResponse{}, err
	}
	return &gset.EmptyResponse{OK: true}, nil
}

func (s *set) SInterStore(ctx context.Context, req *gset.SInterStoreRequest) (*gset.EmptyResponse, error) {
	err := s.dbs.SInterStore(req.Destination, req.Keys...)
	if err != nil {
		return &gset.EmptyResponse{}, err
	}
	return &gset.EmptyResponse{OK: true}, nil
}
