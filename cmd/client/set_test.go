//go:build set_integration

package client

import "testing"

func TestSetAdd(t *testing.T) {
	setCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-add member",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetAdd",
	}

	setCmdMeta.cmdTest(t)
}

func TestSetAddMemberAlreadyExist(t *testing.T) {
	setCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-add-member-already-exist member",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetAddMemberAlreadyExist",
	}

	setCmdMeta.cmdTest(t)

	reSetCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-add-member-already-exist member",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetAddMemberAlreadyExist",
	}

	reSetCmdMeta.cmdTest(t)
}

func TestSetAdds(t *testing.T) {
	setCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-adds member1 member2",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetAdds",
	}

	setCmdMeta.cmdTest(t)
}

func TestSetMembers(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-members member1 member2",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetMembers",
	}

	setAddsCmdMeta.cmdTest(t)

	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-members member3",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetMembers",
	}

	setAddCmdMeta.cmdTest(t)

	setMembersCmdMeta := &cmdMeta{
		cmd: "Smembers for-set-members",
		expectResult: []string{
			"member1",
		},
		testcaseName: "TestSetMembers",
	}

	setMembersCmdMeta.cmdTest(t)

	setMembersCmdMeta = &cmdMeta{
		cmd: "Smembers for-set-members",
		expectResult: []string{
			"member2",
		},
		testcaseName: "TestSetMembers",
	}

	setMembersCmdMeta.cmdTest(t)

	setMembersCmdMeta = &cmdMeta{
		cmd: "Smembers for-set-members",
		expectResult: []string{
			"member3",
		},
		testcaseName: "TestSetMembers",
	}

	setMembersCmdMeta.cmdTest(t)
}

func TestSetMembersKeyIsNotExist(t *testing.T) {
	setMembersCmdMeta := &cmdMeta{
		cmd: "Smembers for-set-members-key-is-not-exist",
		expectResult: []string{
			"SMembers data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetMembersKeyIsNotExist",
	}

	setMembersCmdMeta.cmdTest(t)
}

func TestSetCard(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-card member1 member2",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetCard",
	}

	setAddsCmdMeta.cmdTest(t)

	setCardCmdMeta := &cmdMeta{
		cmd: "Scard for-set-card",
		expectResult: []string{
			"SCard count: 2",
		},
		testcaseName: "TestSetCard",
	}

	setCardCmdMeta.cmdTest(t)
}

func TestSetCardKeyIsNotExist(t *testing.T) {
	setCardCmdMeta := &cmdMeta{
		cmd: "Scard for-set-card-key-is-not-exist",
		expectResult: []string{
			"SCard data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetCardKeyIsNotExist",
	}

	setCardCmdMeta.cmdTest(t)
}

func TestSetRem(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-rem member1 member2",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetRem",
	}

	setAddsCmdMeta.cmdTest(t)

	cardsCmdMeta := &cmdMeta{
		cmd: "Scard for-set-rem",
		expectResult: []string{
			"SCard count: 2",
		},
		testcaseName: "TestSetRem",
	}

	cardsCmdMeta.cmdTest(t)

	setRemCmdMeta := &cmdMeta{
		cmd: "Srem for-set-rem member1",
		expectResult: []string{
			"SRem data success",
		},
		testcaseName: "TestSetRem",
	}

	setRemCmdMeta.cmdTest(t)

	cardsCmdMeta = &cmdMeta{
		cmd: "Scard for-set-rem",
		expectResult: []string{
			"SCard count: 1",
		},
		testcaseName: "TestSetRem",
	}

	cardsCmdMeta.cmdTest(t)
}

func TestSetRemKeyIsNotExist(t *testing.T) {
	setRemCmdMeta := &cmdMeta{
		cmd: "Srem for-set-rem-key-is-not-exist member1",
		expectResult: []string{
			"SRem data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetRemKeyIsNotExist",
	}

	setRemCmdMeta.cmdTest(t)
}

func TestSetRemKeyIsEmpty(t *testing.T) {
	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-rem-key-is-empty member1",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetRemKeyIsEmpty",
	}

	setAddCmdMeta.cmdTest(t)

	setRemCmdMeta := &cmdMeta{
		cmd: "Srem for-set-rem-key-is-empty member1",
		expectResult: []string{
			"SRem data success",
		},
		testcaseName: "TestSetRemKeyIsEmpty",
	}

	setRemCmdMeta.cmdTest(t)

	setRemCmdMeta = &cmdMeta{
		cmd: "Srem for-set-rem-key-is-empty member1",
		expectResult: []string{
			"SRem data error:  rpc error: code = Unknown desc = ErrMemberNotFound: member not found",
			"error: rpc error: code = Unknown desc = ErrMemberNotFound: member not found",
		},
		testcaseName: "TestSetRemKeyIsEmpty",
	}

	setRemCmdMeta.cmdTest(t)
}

func TestSetMembersAllMemberRemove(t *testing.T) {
	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-members-all-member-remove member1",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetMembersAllMemberRemove",
	}

	setAddCmdMeta.cmdTest(t)

	setRemCmdMeta := &cmdMeta{
		cmd: "Srem for-set-members-all-member-remove member1",
		expectResult: []string{
			"SRem data success",
		},
		testcaseName: "TestSetMembersAllMemberRemove",
	}

	setRemCmdMeta.cmdTest(t)

	setMembersCmdMeta := &cmdMeta{
		cmd: "Smembers for-set-members-all-member-remove",
		expectResult: []string{
			"SMembers data: []",
		},
		testcaseName: "TestSetMembersAllMemberRemove",
	}

	setMembersCmdMeta.cmdTest(t)
}

func TestSetRems(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-rems member1 member2",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetRems",
	}

	setAddsCmdMeta.cmdTest(t)

	cardsCmdMeta := &cmdMeta{
		cmd: "Scard for-set-rems",
		expectResult: []string{
			"SCard count: 2",
		},
		testcaseName: "TestSetRems",
	}

	cardsCmdMeta.cmdTest(t)

	setRemsCmdMeta := &cmdMeta{
		cmd: "Srems for-set-rems member1 member2",
		expectResult: []string{
			"SRems data success",
		},
		testcaseName: "TestSetRems",
	}

	setRemsCmdMeta.cmdTest(t)

	cardsCmdMeta = &cmdMeta{
		cmd: "Scard for-set-rems",
		expectResult: []string{
			"SCard count: 0",
		},
		testcaseName: "TestSetRems",
	}

	cardsCmdMeta.cmdTest(t)
}

func TestSetRemsKeyIsNoExist(t *testing.T) {
	setRemsCmdMeta := &cmdMeta{
		cmd: "Srems for-set-rems-key-is-not-exist member1 member2",
		expectResult: []string{
			"SRems data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetRemsKeyIsNoExist",
	}

	setRemsCmdMeta.cmdTest(t)
}

func TestSetRemsKeyIsEmpty(t *testing.T) {
	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-rems-key-is-empty member1",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetRemsKeyIsEmpty",
	}

	setAddCmdMeta.cmdTest(t)

	setRemsCmdMeta := &cmdMeta{
		cmd: "Srems for-set-rems-key-is-empty member1",
		expectResult: []string{
			"SRems data success",
		},
		testcaseName: "TestSetRemsKeyIsEmpty",
	}

	setRemsCmdMeta.cmdTest(t)

	setRemsCmdMeta = &cmdMeta{
		cmd: "Srems for-set-rems-key-is-empty member1",
		expectResult: []string{
			"SRems data error:  rpc error: code = Unknown desc = ErrMemberNotFound: member not found",
			"error: rpc error: code = Unknown desc = ErrMemberNotFound: member not found",
		},
		testcaseName: "TestSetRemsKeyIsEmpty",
	}

	setRemsCmdMeta.cmdTest(t)
}

func TestSetRemsContainNotFoundMember(t *testing.T) {
	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-rems-contain-not-found-member member1",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetRemsContainNotFoundMember",
	}

	setAddCmdMeta.cmdTest(t)

	setRemsCmdMeta := &cmdMeta{
		cmd: "Srems for-set-rems-contain-not-found-member member1 member2",
		expectResult: []string{
			"SRems data error:  rpc error: code = Unknown desc = ErrMemberNotFound: member not found",
			"error: rpc error: code = Unknown desc = ErrMemberNotFound: member not found",
		},
		testcaseName: "TestSetRemsContainNotFoundMember",
	}

	setRemsCmdMeta.cmdTest(t)

	cardsCmdMeta := &cmdMeta{
		cmd: "Scard for-set-rems-contain-not-found-member",
		expectResult: []string{
			"SCard count: 1",
		},
		testcaseName: "TestSetRemsContainNotFoundMember",
	}

	cardsCmdMeta.cmdTest(t)
}

func TestSetIsMember(t *testing.T) {
	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-ismember member1",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetIsMember",
	}

	setAddCmdMeta.cmdTest(t)

	setIsMemberCmdMeta := &cmdMeta{
		cmd: "Sismember for-set-ismember member1",
		expectResult: []string{
			"SIsMember result: true",
		},
		testcaseName: "TestSetIsMember",
	}

	setIsMemberCmdMeta.cmdTest(t)

	setIsMemberCmdMeta = &cmdMeta{
		cmd: "Sismember for-set-ismember member2",
		expectResult: []string{
			"SIsMember result: false",
		},
		testcaseName: "TestSetIsMember",
	}

	setIsMemberCmdMeta.cmdTest(t)
}

func TestSetIsMemberKeyIsNotExist(t *testing.T) {
	setIsMemberCmdMeta := &cmdMeta{
		cmd: "Sismember for-set-ismember-key-is-not-exist member1",
		expectResult: []string{
			"SIsMember data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetIsMemberKeyIsNotExist",
	}

	setIsMemberCmdMeta.cmdTest(t)
}

func TestSetUnion(t *testing.T) {
	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-union-key-1 member1",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetUnion",
	}

	setAddCmdMeta.cmdTest(t)

	setAddCmdMeta = &cmdMeta{
		cmd: "Sadd for-set-union-key-2 member2",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetUnion",
	}

	setAddCmdMeta.cmdTest(t)

	setUnionCmdMeta := &cmdMeta{
		cmd: "Sunion for-set-union-key-1 for-set-union-key-2",
		expectResult: []string{
			"SUnion result: [member1 member2]",
		},
		testcaseName: "TestSetUnion",
	}

	setUnionCmdMeta.cmdTest(t)
}

func TestSetUnionContainNotExistKey(t *testing.T) {
	setAddCmdMeta := &cmdMeta{
		cmd: "Sadd for-set-union-contain-not-exist-key member1",
		expectResult: []string{
			"SAdd data success",
		},
		testcaseName: "TestSetUnionContainNotExistKey",
	}

	setAddCmdMeta.cmdTest(t)

	setUnionCmdMeta := &cmdMeta{
		cmd: "Sunion for-set-union-contain-not-exist-key for-set-union-contain-not-exist-key-2",
		expectResult: []string{
			"SUnion data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetUnionContainNotExistKey",
	}

	setUnionCmdMeta.cmdTest(t)
}

func TestSetInter(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-inter-key-1 member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetInter",
	}

	setAddsCmdMeta.cmdTest(t)

	setAddsCmdMeta = &cmdMeta{
		cmd: "Sadds for-set-inter-key-2 member1 key2-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetInter",
	}

	setAddsCmdMeta.cmdTest(t)

	setInterCmdMeta := &cmdMeta{
		cmd: "Sinter for-set-inter-key-1 for-set-inter-key-2",
		expectResult: []string{
			"SInter result: [member1]",
		},
		testcaseName: "TestSetInter",
	}

	setInterCmdMeta.cmdTest(t)
}

func TestSetInterContainNotExistKey(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-inter-contain-not-exist-key member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetInterContainNotExistKey",
	}

	setAddsCmdMeta.cmdTest(t)

	setInterCmdMeta := &cmdMeta{
		cmd: "Sinter for-set-inter-contain-not-exist-key for-set-inter-contain-not-exist-key-2",
		expectResult: []string{
			"SInter data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetInterContainNotExistKey",
	}

	setInterCmdMeta.cmdTest(t)
}

func TestSetDiff(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-diff-key-1 member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetDiff",
	}

	setAddsCmdMeta.cmdTest(t)

	setAddsCmdMeta = &cmdMeta{
		cmd: "Sadds for-set-diff-key-2 member1 key2-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetDiff",
	}

	setAddsCmdMeta.cmdTest(t)

	setDiffCmdMeta := &cmdMeta{
		cmd: "Sdiff for-set-diff-key-1 for-set-diff-key-2",
		expectResult: []string{
			"SDiff result: [key1-member]",
		},
		testcaseName: "TestSetDiff",
	}

	setDiffCmdMeta.cmdTest(t)
}

func TestSetDiffContainNotExistKey(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-diff-contain-not-exist-key member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetDiffContainNotExistKey",
	}

	setAddsCmdMeta.cmdTest(t)

	setDiffCmdMeta := &cmdMeta{
		cmd: "Sdiff for-set-diff-contain-not-exist-key for-set-diff-contain-not-exist-key-2",
		expectResult: []string{
			"SDiff data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetDiffContainNotExistKey",
	}

	setDiffCmdMeta.cmdTest(t)
}

func TestSetUnionStore(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-union-store-key-1 member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetUnionStore",
	}

	setAddsCmdMeta.cmdTest(t)

	setAddsCmdMeta = &cmdMeta{
		cmd: "Sadds for-set-union-store-key-2 member1 key2-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetUnionStore",
	}

	setAddsCmdMeta.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "Sunionstore for-set-unions-store-des for-set-union-store-key-1 for-set-union-store-key-2",
		expectResult: []string{
			"SUnionStore success",
		},
		testcaseName: "TestSetUnionStore",
	}

	cmdMetaData.cmdTest(t)

	membersCmdMeta := &cmdMeta{
		cmd: "Smembers for-set-unions-store-des",
		expectResult: []string{
			"member1",
		},
		testcaseName: "TestSetUnionStore",
	}

	membersCmdMeta.cmdTest(t)

	membersCmdMeta = &cmdMeta{
		cmd: "Smembers for-set-unions-store-des",
		expectResult: []string{
			"key1-member",
		},
		testcaseName: "TestSetUnionStore",
	}

	membersCmdMeta.cmdTest(t)

	membersCmdMeta = &cmdMeta{
		cmd: "Smembers for-set-unions-store-des",
		expectResult: []string{
			"key2-member",
		},
		testcaseName: "TestSetUnionStore",
	}

	membersCmdMeta.cmdTest(t)

	cardCmdMeta := &cmdMeta{
		cmd: "Scard for-set-unions-store-des",
		expectResult: []string{
			"SCard count: 3",
		},
		testcaseName: "TestSetUnionStore",
	}

	cardCmdMeta.cmdTest(t)
}

func TestSetUnionStoreContainNotExistKey(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-union-store-contain-not-exist-key member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetUnionStoreContainNotExistKey",
	}

	setAddsCmdMeta.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "Sunionstore for-set-union-store-contain-not-exist-key for-set-union-store-contain-not-exist-key-2",
		expectResult: []string{
			"SUnionStore data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetUnionStoreContainNotExistKey",
	}

	cmdMetaData.cmdTest(t)
}

func TestSetInterStore(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-inter-store-key-1 member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetInterStore",
	}

	setAddsCmdMeta.cmdTest(t)

	setAddsCmdMeta = &cmdMeta{
		cmd: "Sadds for-set-inter-store-key-2 member1 key2-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetInterStore",
	}

	setAddsCmdMeta.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "Sinterstore for-set-inter-store-des for-set-inter-store-key-1 for-set-inter-store-key-2",
		expectResult: []string{
			"SInterStore success",
		},
		testcaseName: "TestSetInterStore",
	}

	cmdMetaData.cmdTest(t)

	membersCmdMeta := &cmdMeta{
		cmd: "Smembers for-set-inter-store-des",
		expectResult: []string{
			"member1",
		},
		testcaseName: "TestSetInterStore",
	}

	membersCmdMeta.cmdTest(t)

	cardCmdMeta := &cmdMeta{
		cmd: "Scard for-set-inter-store-des",
		expectResult: []string{
			"SCard count: 1",
		},
		testcaseName: "TestSetInterStore",
	}

	cardCmdMeta.cmdTest(t)
}

func TestSetInterStoreContainNotExistKey(t *testing.T) {
	setAddsCmdMeta := &cmdMeta{
		cmd: "Sadds for-set-inter-store-contain-not-exist-key member1 key1-member",
		expectResult: []string{
			"SAdds data success",
		},
		testcaseName: "TestSetInterStoreContainNotExistKey",
	}

	setAddsCmdMeta.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "Sinterstore for-set-inter-store-contain-not-exist-key for-set-inter-store-contain-not-exist-key-2",
		expectResult: []string{
			"SInterStore data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestSetInterStoreContainNotExistKey",
	}

	cmdMetaData.cmdTest(t)
}
