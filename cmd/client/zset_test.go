//go:build zset_integration

package client

import "testing"

func TestZAdd(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zadd member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZAdd",
	}

	cmdMetaData.cmdTest(t)
}

func TestZCount(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zcount member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZCount",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZCount for-zcount 0 1",
		expectResult: []string{
			"ZCount data success, count: 1",
		},
		testcaseName: "TestZCount",
	}

	cmdMetaData.cmdTest(t)
}

func TestZCountKeyNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZCount for-zcout-key-not-exist 0 1",
		expectResult: []string{
			"ZCount data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZCountKeyNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetCard(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zcard member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetCard",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZCard for-zcard",
		expectResult: []string{
			"ZCard data success, count: 1",
		},
		testcaseName: "TestZSetCard",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetCardKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZCard for-zcard-key-not-exist",
		expectResult: []string{
			"ZCard data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetCardKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRem(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrem member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRem",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRem for-zrem member",
		expectResult: []string{
			"ZRem data success",
		},
		testcaseName: "TestZSetRem",
	}

	cmdMetaData.cmdTest(t)

	cardCmdMetaData := &cmdMeta{
		cmd: "ZCard for-zrem",
		expectResult: []string{
			"ZCard data success, count: 0",
		},
		testcaseName: "TestZSetRem",
	}

	cardCmdMetaData.cmdTest(t)
}

func TestZSetRemKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZRem for-zrem-key-not-exist member",
		expectResult: []string{
			"ZRem data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRemKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRemMemberIsNotExist(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrem-member-not-exist member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRemMemberIsNotExist",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRem for-zrem-member-not-exist member-not-exist",
		expectResult: []string{
			"ZRem data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRemMemberIsNotExist",
	}

	cmdMetaData.cmdTest(t)

	cardCmdMetaData := &cmdMeta{
		cmd: "ZCard for-zrem-member-not-exist",
		expectResult: []string{
			"ZCard data success, count: 1",
		},
		testcaseName: "TestZSetRemMemberIsNotExist",
	}

	cardCmdMetaData.cmdTest(t)
}

func TestZSetRems(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrems member1 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRems",
	}

	addCmdMetaData.cmdTest(t)

	addCmdMetaData = &cmdMeta{
		cmd: "ZAdd for-zrems member2 2 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRems",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRems for-zrems member1 member2",
		expectResult: []string{
			"ZRems data success",
		},
		testcaseName: "TestZSetRems",
	}

	cmdMetaData.cmdTest(t)

	cardCmdMetaData := &cmdMeta{
		cmd: "ZCard for-zrems",
		expectResult: []string{
			"ZCard data success, count: 0",
		},
		testcaseName: "TestZSetRems",
	}

	cardCmdMetaData.cmdTest(t)
}

func TestZSetRemsKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZRems for-zrems-key-not-exist member1 member2",
		expectResult: []string{
			"ZRems data error:  rpc error: code = Unknown desc = failed to get or create ZSet from DB with key 'for-zrems-key-not-exist': KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = failed to get or create ZSet from DB with key 'for-zrems-key-not-exist': KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRemsKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRemsMemberIsNotExist(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrems-member-not-exist member1 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRemsMemberIsNotExist",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRems for-zrems-member-not-exist member1 member2",
		expectResult: []string{
			"ZRems data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRemsMemberIsNotExist",
	}

	cmdMetaData.cmdTest(t)

	cardCmdMetaData := &cmdMeta{
		cmd: "ZCard for-zrems-member-not-exist",
		expectResult: []string{
			"ZCard data success, count: 1",
		},
		testcaseName: "TestZSetRemsMemberIsNotExist",
	}

	cardCmdMetaData.cmdTest(t)
}

func TestZSetScore(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zscore member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetScore",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZScore for-zscore member",
		expectResult: []string{
			"ZScore data success, score: 1",
		},
		testcaseName: "TestZSetScore",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetScoreKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZScore for-zscore-key-not-exist member",
		expectResult: []string{
			"ZScore data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetScoreKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetScoreMemberIsNotExist(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zscore-member-not-exist member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetScoreMemberIsNotExist",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZScore for-zscore-member-not-exist member-not-exist",
		expectResult: []string{
			"ZScore data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetScoreMemberIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRank(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrank member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRank",
	}

	addCmdMetaData.cmdTest(t)

	addCmdMetaData = &cmdMeta{
		cmd: "ZAdd for-zrank member2 2 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRank",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRank for-zrank member",
		expectResult: []string{
			"ZRank data success, rank: 1",
		},
		testcaseName: "TestZSetRank",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRankKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZRank for-zrank-key-not-exist member",
		expectResult: []string{
			"ZRank data error:  rpc error: code = Unknown desc = failed to get or create ZSet from DB with key 'for-zrank-key-not-exist': KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = failed to get or create ZSet from DB with key 'for-zrank-key-not-exist': KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRankKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRankMemberIsNotExist(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrank-member-not-exist member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRankMemberIsNotExist",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRank for-zrank-member-not-exist member-not-exist",
		expectResult: []string{
			"ZRank data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRankMemberIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRevRank(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrevrank member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRevRank",
	}

	addCmdMetaData.cmdTest(t)

	addCmdMetaData = &cmdMeta{
		cmd: "ZAdd for-zrevrank member2 2 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRevRank",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRevRank for-zrevrank member",
		expectResult: []string{
			"ZRevRank data success, rank: 2",
		},
		testcaseName: "TestZSetRevRank",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRevRankKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZRevRank for-zrevrank-key-not-exist member",
		expectResult: []string{
			"ZRevRank data error:  rpc error: code = Unknown desc = failed to get or create ZSet from DB with key 'for-zrevrank-key-not-exist': KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = failed to get or create ZSet from DB with key 'for-zrevrank-key-not-exist': KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRevRankKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRevRankMemberIsNotExist(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrevrank-member-not-exist member 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRevRankMemberIsNotExist",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRevRank for-zrevrank-member-not-exist member-not-exist",
		expectResult: []string{
			"ZRevRank data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRevRankMemberIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRange(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrange member1 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRange",
	}

	addCmdMetaData.cmdTest(t)

	addCmdMetaData = &cmdMeta{
		cmd: "ZAdd for-zrange member2 2 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRange",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRange for-zrange 0 1",
		expectResult: []string{
			"ZRange data success",
			"Member: member1, Score: 1, Value: value",
		},
		testcaseName: "TestZSetRange",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRangeKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZRange for-zrange-key-not-exist 0 1",
		expectResult: []string{
			"ZRange data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRangeKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRevRange(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zrevrange member1 1 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRevRange",
	}

	addCmdMetaData.cmdTest(t)

	addCmdMetaData = &cmdMeta{
		cmd: "ZAdd for-zrevrange member2 2 value",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetRevRange",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZRevRange for-zrevrange 0 1",
		expectResult: []string{
			"ZRevRange data success",
			"Member: member2, Score: 2, Value: value",
		},
		testcaseName: "TestZSetRevRange",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetRevRangeKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "ZRevRange for-zrevrange-key-not-exist 0 1",
		expectResult: []string{
			"ZRevRange data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestZSetRevRangeKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestZSetIncrBy(t *testing.T) {
	addCmdMetaData := &cmdMeta{
		cmd: "ZAdd for-zincrby member 1 2",
		expectResult: []string{
			"ZAdd data success",
		},
		testcaseName: "TestZSetIncrBy",
	}

	addCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "ZIncrBy for-zincrby member 1",
		expectResult: []string{
			"ZIncrBy data success, new score: 2",
		},
		testcaseName: "TestZSetIncrBy",
	}

	cmdMetaData.cmdTest(t)
}
