package client

import "testing"

func TestStringLPush(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LPush for-string-lPush value",
		expectResult: []string{
			"LPush data success",
		},
		testcaseName: "TestStringLPush",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLPushsData(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lPushs value1 value2 value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLPushsData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringRPushData(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "RPush for-string-rPush value",
		expectResult: []string{
			"RPush data success",
		},
		testcaseName: "TestStringRPushData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringRPushsData(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "RPushs for-string-rPushs value1 value2 value3",
		expectResult: []string{
			"RPushs data success",
		},
		testcaseName: "TestStringRPushsData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLPopData(t *testing.T) {
	lpushsCmd := &cmdMeta{
		cmd: "LPushs for-string-lPop lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLPopData",
	}

	lpushsCmd.cmdTest(t)

	lPushCmd := &cmdMeta{
		cmd: "LPush for-string-lPop lpush-value",
		expectResult: []string{
			"LPush data success",
		},
		testcaseName: "TestStringLPopData",
	}

	lPushCmd.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LPop for-string-lPop",
		expectResult: []string{
			"LPop data success: lpush-value",
		},
		testcaseName: "TestStringLPopData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LPop for-string-lPop",
		expectResult: []string{
			"LPop data success: lpushs-value1",
		},
		testcaseName: "TestStringLPopData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LPop for-string-lPop",
		expectResult: []string{
			"LPop data success: lpushs-value2",
		},
		testcaseName: "TestStringLPopData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LPop for-string-lPop",
		expectResult: []string{
			"LPop data success: lpushs-value3",
		},
		testcaseName: "TestStringLPopData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLPopDataRightPush(t *testing.T) {
	rPushsCmd := &cmdMeta{
		cmd: "RPushs for-string-lPop-right-push rpushs-value1 rpushs-value2 rpushs-value3",
		expectResult: []string{
			"RPushs data success",
		},
		testcaseName: "TestStringLPopDataRightPush",
	}

	rPushsCmd.cmdTest(t)

	rPushCmd := &cmdMeta{
		cmd: "RPush for-string-lPop-right-push rpush-value",
		expectResult: []string{
			"RPush data success",
		},
		testcaseName: "TestStringLPopDataRightPush",
	}

	rPushCmd.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LPop for-string-lPop-right-push",
		expectResult: []string{
			"LPop data success: rpushs-value1",
		},
		testcaseName: "TestStringLPopDataRightPush",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LPop for-string-lPop-right-push",
		expectResult: []string{
			"LPop data success: rpushs-value2",
		},
		testcaseName: "TestStringLPopDataRightPush",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LPop for-string-lPop-right-push",
		expectResult: []string{
			"LPop data success: rpushs-value3",
		},
		testcaseName: "TestStringLPopDataRightPush",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LPop for-string-lPop-right-push",
		expectResult: []string{
			"LPop data success: rpush-value",
		},
		testcaseName: "TestStringLPopDataRightPush",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLPopListKeyIsEmpty(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LPop for-string-lPop-list-key-is-empty",
		expectResult: []string{
			"LPop data error: client LPop failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client LPop failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringLPopListEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringRPopData(t *testing.T) {
	rPushsCmdMetaData := &cmdMeta{
		cmd: "RPushs for-string-rPop rpushs-value1 rpushs-value2 rpushs-value3",
		expectResult: []string{
			"RPushs data success",
		},
		testcaseName: "TestStringRPopData",
	}

	rPushsCmdMetaData.cmdTest(t)

	rPushCmdMetaData := &cmdMeta{
		cmd: "RPush for-string-rPop rpush-value",
		expectResult: []string{
			"RPush data success",
		},
		testcaseName: "TestStringRPopData",
	}

	rPushCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "RPop for-string-rPop",
		expectResult: []string{
			"RPop data success: rpush-value",
		},
		testcaseName: "TestStringRPopData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "RPop for-string-rPop",
		expectResult: []string{
			"RPop data success: rpushs-value3",
		},
		testcaseName: "TestStringRPopData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "RPop for-string-rPop",
		expectResult: []string{
			"RPop data success: rpushs-value2",
		},
		testcaseName: "TestStringRPopData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "RPop for-string-rPop",
		expectResult: []string{
			"RPop data success: rpushs-value1",
		},
		testcaseName: "TestStringRPopData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringRPopDataListLeftPush(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-rPop-left-push lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringRPopDataListLeftPush",
	}

	lPushsCmdMetaData.cmdTest(t)

	lPushCmdMetaData := &cmdMeta{
		cmd: "LPush for-string-rPop-left-push lpush-value",
		expectResult: []string{
			"LPush data success",
		},
		testcaseName: "TestStringRPopDataListLeftPush",
	}

	lPushCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "RPop for-string-rPop-left-push",
		expectResult: []string{
			"RPop data success: lpushs-value3",
		},
		testcaseName: "TestStringRPopDataListLeftPush",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "RPop for-string-rPop-left-push",
		expectResult: []string{
			"RPop data success: lpushs-value2",
		},
		testcaseName: "TestStringRPopDataListLeftPush",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "RPop for-string-rPop-left-push",
		expectResult: []string{
			"RPop data success: lpushs-value1",
		},
		testcaseName: "TestStringRPopDataListLeftPush",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "RPop for-string-rPop-left-push",
		expectResult: []string{
			"RPop data success: lpush-value",
		},
		testcaseName: "TestStringRPopDataListLeftPush",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringRPopListKeyIsEmpty(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "RPop for-string-rPop-list-key-is-empty",
		expectResult: []string{
			"RPop data error: client RPop failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client RPop failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringRPopListEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLRangeData(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lRange lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLRangeData",
	}

	lPushsCmdMetaData.cmdTest(t)

	lPushCmdMetaData := &cmdMeta{
		cmd: "LPush for-string-lRange lpush-value",
		expectResult: []string{
			"LPush data success",
		},
		testcaseName: "TestStringLRangeData",
	}

	lPushCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LRange for-string-lRange 0 1",
		expectResult: []string{
			"LRange data success: [lpush-value lpushs-value1]",
		},
		testcaseName: "TestStringLRangeData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LRange for-string-lRange 0 2",
		expectResult: []string{
			"LRange data success: [lpush-value lpushs-value1 lpushs-value2]",
		},
		testcaseName: "TestStringLRangeData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LRange for-string-lRange 0 3",
		expectResult: []string{
			"LRange data success: [lpush-value lpushs-value1 lpushs-value2 lpushs-value3]",
		},
		testcaseName: "TestStringLRangeData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLRangeDataListKeyIsEmpty(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LRange for-string-lRange-list-key-is-empty 0 1",
		expectResult: []string{
			"LRange data error: client LRange failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client LRange failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringLRangeListEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLLenData(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lLen lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLLenData",
	}

	lPushsCmdMetaData.cmdTest(t)

	lPushCmdMetaData := &cmdMeta{
		cmd: "LPush for-string-lLen lpush-value",
		expectResult: []string{
			"LPush data success",
		},
		testcaseName: "TestStringLLenData",
	}

	lPushCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LLen for-string-lLen",
		expectResult: []string{
			"LLen data success: 4",
		},
		testcaseName: "TestStringLLenData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLLenDataListKeyIsEmpty(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LLen for-string-lLen-list-key-is-empty",
		expectResult: []string{
			"LLen data error: client LLen failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client LLen failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringLLenListEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLRemDataCount0(t *testing.T) {
	lPushCmdMetaData := &cmdMeta{
		cmd: "LPush for-string-lRem lpush-value1",
		expectResult: []string{
			"LPush data success",
		},
		testcaseName: "TestStringLRemDataCount0",
	}

	lPushCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LRem for-string-lRem 0 lpush-value1",
		expectResult: []string{
			"LRem data success",
		},
		testcaseName: "TestStringLRemDataCount0",
	}

	cmdMetaData.cmdTest(t)

	popCmdMetaData := &cmdMeta{
		cmd: "LPop for-string-lRem",
		expectResult: []string{
			"LPop data error: client LPop failed: rpc error: code = Unknown desc = Wrong operation: list is empty",
			"error: client LPop failed: rpc error: code = Unknown desc = Wrong operation: list is empty",
		},
		testcaseName: "TestStringLRemDataCount0",
	}

	popCmdMetaData.cmdTest(t)
}

func TestStringLRemDataCount1(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lRem lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLRemDataCount1",
	}

	lPushsCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LRem for-string-lRem 1 lpushs-value1",
		expectResult: []string{
			"LRem data success",
		},
		testcaseName: "TestStringLRemDataCount1",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LRange for-string-lRem 0 1",
		expectResult: []string{
			"LRange data success: [lpushs-value2 lpushs-value3]",
		},
		testcaseName: "TestStringLRemDataCount1",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLRemDataCountLetterThan0(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lRem lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLRemDataCountLetterThan0",
	}

	lPushsCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LRem for-string-lRem -1 lpushs-value1",
		expectResult: []string{
			"LRem data success",
		},
		testcaseName: "TestStringLRemDataCountLetterThan0",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LRange for-string-lRem 0 1",
		expectResult: []string{
			"LRange data success: [lpushs-value2 lpushs-value3]",
		},
		testcaseName: "TestStringLRemDataCountLetterThan0",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLRemDataKeyIsNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LRem for-string-lRem-key-is-not-exists 0 lpushs-value1",
		expectResult: []string{
			"LRem data error: client LRem failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client LRem failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringLRemDataKeyIsNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLRemDataNeedToRemoveNotExist(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lRem-need-to-remove-not-exist lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLRemDataNeedToRemoveNotExist",
	}

	lPushsCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LRem for-string-lRem-need-to-remove-not-exist 1 lpushs-value4",
		expectResult: []string{
			"LRem data success",
		},
		testcaseName: "TestStringLRemDataNeedToRemoveNotExist",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LRange for-string-lRem-need-to-remove-not-exist 0 2",
		expectResult: []string{
			"LRange data success: [lpushs-value1 lpushs-value2 lpushs-value3]",
		},
		testcaseName: "TestStringLRemDataNeedToRemoveNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLIndexData(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lIndex lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLIndexData",
	}

	lPushsCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LIndex for-string-lIndex 0",
		expectResult: []string{
			"LIndex data success: lpushs-value1",
		},
		testcaseName: "TestStringLIndexData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LIndex for-string-lIndex 1",
		expectResult: []string{
			"LIndex data success: lpushs-value2",
		},
		testcaseName: "TestStringLIndexData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LIndex for-string-lIndex 2",
		expectResult: []string{
			"LIndex data success: lpushs-value3",
		},
		testcaseName: "TestStringLIndexData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLIndexDataKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LIndex for-string-lIndex-key-is-not-exists 0",
		expectResult: []string{
			"LIndex data error: client LIndex failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client LIndex failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringLIndexDataKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLSetData(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lSet lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLSetData",
	}

	lPushsCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LSet for-string-lSet 0 lpushs-value4",
		expectResult: []string{
			"LSet data success",
		},
		testcaseName: "TestStringLSetData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LRange for-string-lSet 0 2",
		expectResult: []string{
			"LRange data success: [lpushs-value4 lpushs-value2 lpushs-value3]",
		},
		testcaseName: "TestStringLSetData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLSetDataKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LSet for-string-lSet-key-is-not-exists 0 lpushs-value4",
		expectResult: []string{
			"LSet data error: client LSet failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client LSet failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringLSetDataKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLSetDataIndexIsOutOfRange(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lSet-index-is-out-of-range lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLSetDataIndexIsOutOfRange",
	}

	lPushsCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LSet for-string-lSet-index-is-out-of-range 3 lpushs-value4",
		expectResult: []string{
			"LSet data error: client LSet failed: rpc error: code = Unknown desc = Wrong operation: index out of range",
			"error: client LSet failed: rpc error: code = Unknown desc = Wrong operation: index out of range",
		},
		testcaseName: "TestStringLSetDataIndexIsOutOfRange",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLTrimData(t *testing.T) {
	lPushsCmdMetaData := &cmdMeta{
		cmd: "LPushs for-string-lTrim lpushs-value1 lpushs-value2 lpushs-value3",
		expectResult: []string{
			"LPushs data success",
		},
		testcaseName: "TestStringLTrimData",
	}

	lPushsCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "LTrim for-string-lTrim 0 1",
		expectResult: []string{
			"LTrim data success",
		},
		testcaseName: "TestStringLTrimData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LLen for-string-lTrim",
		expectResult: []string{
			"LLen data success: 2",
		},
		testcaseName: "TestStringLTrimData",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = &cmdMeta{
		cmd: "LRange for-string-lTrim 0 1",
		expectResult: []string{
			"LRange data success: [lpushs-value1 lpushs-value2]",
		},
		testcaseName: "TestStringLTrimData",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringLTrimDataKeyIsNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "LTrim for-string-lTrim-key-is-not-exists 0 1",
		expectResult: []string{
			"LTrim data error: client LTrim failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: client LTrim failed: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringLTrimDataKeyIsNotExist",
	}

	cmdMetaData.cmdTest(t)
}
