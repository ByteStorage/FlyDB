package dh

import (
	"errors"
	"github.com/ByteStorage/FlyDB/lib/raft"
	"github.com/klauspost/reedsolomon"
)

type EcEncoder struct {
	reedsolomon.Encoder
}

func NewEcEncoder(dataNum int, parityNum int) (EcEncoder, error) {
	encoder, err := reedsolomon.New(dataNum, parityNum)
	if err != nil {
		return EcEncoder{}, err
	}
	return EcEncoder{encoder}, nil
}

func (ec *EcEncoder) AssignData(data []byte, slaveAddrList []cluster.Slave) (map[string][]byte, error) {
	dataShards, err := ec.Split(data)
	if err != nil {
		return nil, err
	}
	// Check if the number of shards is more than the number of slaves
	if len(dataShards) > len(slaveAddrList) {
		return nil, errors.New("number of data shards is more than the number of slaves")
	}
	assignment := make(map[string][]byte)
	for i, shard := range dataShards {
		assignment[slaveAddrList[i].Addr] = shard
	}
	return assignment, nil
}
