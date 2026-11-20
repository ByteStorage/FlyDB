//go:build hash_integration

package client

import (
	"testing"
)

func TestHSetData(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HSet for-hset-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHSetData",
	}

	cmdMetaData.cmdTest(t)
}

func TestHSetDataAlreadyExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HSet for-hset-already-exist field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "HSetDataAlreadyExist",
	}

	cmdMetaData.cmdTest(t)

	rePutCmdMetaData := &cmdMeta{
		cmd: "HSet for-hset-already-exist field re-value",
		expectResult: []string{
			"put data error:  put failed",
			"error: put failed",
		},
		testcaseName: "TestHSetDataAlreadyExist",
	}

	rePutCmdMetaData.cmdTest(t)

	getCmdMetaData := &cmdMeta{
		cmd: "HGet for-hset-already-exist field",
		expectResult: []string{
			"re-value",
		},
		testcaseName: "TestHSetDataAlreadyExist",
	}

	getCmdMetaData.cmdTest(t)
}

func TestHGetData(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hget-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "HGetData",
	}

	putCmdMetaData.cmdTest(t)

	getCmdMetaData := &cmdMeta{
		cmd: "HGet for-hget-data field",
		expectResult: []string{
			"value",
		},
		testcaseName: "TestHGetData",
	}

	getCmdMetaData.cmdTest(t)
}

func TestHGetKeyNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HGet for-hget-not-exists field",
		expectResult: []string{
			"Field or Key does not exist in the hash",
		},
		testcaseName: "TestHGetKeyNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHGetFiledNotExists(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hget-field-not-exists field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHGetFiledNotExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HGet for-hget-field-not-exists field-not-exists",
		expectResult: []string{
			"get data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestHGetFiledNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHashDel(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hdel-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHashDel",
	}

	putCmdMetaData.cmdTest(t)

	getCmdMetaData := &cmdMeta{
		cmd: "HGet for-hdel-data field",
		expectResult: []string{
			"value",
		},
		testcaseName: "TestHashDel",
	}

	getCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HDel for-hdel-data field",
		expectResult: []string{
			"delete key success",
		},
		testcaseName: "TestHashDel",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData = &cmdMeta{
		cmd: "HGet for-hdel-data field",
		expectResult: []string{
			"Field or Key does not exist in the hash",
		},
		testcaseName: "TestHashDel",
	}

	getCmdMetaData.cmdTest(t)
}

func TestHDelNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HDel for-hdel-not-exists field",
		expectResult: []string{
			"delete key error:  del failed",
			"error: del failed",
		},
		testcaseName: "TestHDelNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHExists(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hexists-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HExists for-hexists-data field",
		expectResult: []string{
			"Field exists in the hash",
		},
		testcaseName: "TestHExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHExistsButFieldNotExist(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hexists-but-field-not-exist-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHExistsButFieldNotExist",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HExists for-hexists-but-field-not-exist-data field-not-exists",
		expectResult: []string{
			"Field does not exist in the hash",
		},
		testcaseName: "TestHExistsButNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestHExistsButKeyNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HExists for-hexists-but-key-not-exist-data field",
		expectResult: []string{
			"Field does not exist in the hash",
		},
		testcaseName: "TestHExistsButKeyNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestHLen(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hlen-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHLen",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HLen for-hlen-data",
		expectResult: []string{
			"1",
		},
		testcaseName: "TestHLen",
	}

	cmdMetaData.cmdTest(t)
}

func TestHLenKeyNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HLen for-hlen-key-not-exists",
		expectResult: []string{
			"0",
		},
		testcaseName: "TestHLenKeyNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHUpdate(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hupdate-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHUpdate",
	}

	putCmdMetaData.cmdTest(t)

	getCmdMetaData := &cmdMeta{
		cmd: "HGet for-hupdate-data field",
		expectResult: []string{
			"value",
		},
		testcaseName: "TestHUpdate",
	}

	getCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HUpdate for-hupdate-data field new-value",
		expectResult: []string{
			"Hash updated successfully",
		},
		testcaseName: "TestHUpdate",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData = &cmdMeta{
		cmd: "HGet for-hupdate-data field",
		expectResult: []string{
			"new-value",
		},
		testcaseName: "TestHUpdate",
	}

	getCmdMetaData.cmdTest(t)
}

func TestHUpdateKeyNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HUpdate for-hupdate-key-not-exists field value",
		expectResult: []string{
			"HUpdate error:  HSet failed",
			"error: HSet failed",
		},
		testcaseName: "TestHUpdateKeyNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHUpdateFieldNotExists(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hupdate-field-not-exists field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHUpdateFieldNotExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HUpdate for-hupdate-field-not-exists field-not-exists value",
		expectResult: []string{
			"HUpdate error:  HSet failed",
			"error: HSet failed",
		},
		testcaseName: "TestHUpdateFieldNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByKey(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hincr-by-key field 1",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHIncrByKey",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HIncrby for-hincr-by-key field 1",
		expectResult: []string{
			"2",
		},
		testcaseName: "TestHIncrByKey",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByKeyNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HIncrby for-hincr-by-key-not-exists field 1",
		expectResult: []string{
			"0",
		},
		testcaseName: "TestHIncrByKeyNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByFieldNotExists(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hincr-by-field-not-exists field 1",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHIncrByFieldNotExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HIncrby for-hincr-by-field-not-exists field-not-exists 1",
		expectResult: []string{
			"0",
		},
		testcaseName: "TestHIncrByFieldNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByFieldNotInt(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hincr-by-field-not-int field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHIncrByFieldNotInt",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HIncrby for-hincr-by-field-not-int field 1",
		expectResult: []string{
			"HIncrBy error:  rpc error: code = Unknown desc = strconv.ParseInt: parsing \"value\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.ParseInt: parsing \"value\": invalid syntax",
		},
		testcaseName: "TestHIncrByFieldNotInt",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByFloatKey(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hincr-by-float-key field 1",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHIncrByFloatKey",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HIncrbyfloat for-hincr-by-float-key field 1.1",
		expectResult: []string{
			"2.1",
		},
		testcaseName: "TestHIncrByFloatKey",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByFloatKeyNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HIncrbyfloat for-hincr-by-float-key-not-exists field 1.1",
		expectResult: []string{
			"0",
		},
		testcaseName: "TestHIncrByFloatKeyNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByFloatFieldNotExists(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hincr-by-float-field-not-exists field 1",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHIncrByFloatFieldNotExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HIncrbyfloat for-hincr-by-float-field-not-exists field-not-exists 1.1",
		expectResult: []string{
			"0",
		},
		testcaseName: "TestHIncrByFloatFieldNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHIncrByFloatFieldNotFloat(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hincr-by-float-field-not-float field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHIncrByFloatFieldNotFloat",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HIncrbyfloat for-hincr-by-float-field-not-float field 1.1",
		expectResult: []string{
			"HIncrByFloat error:  rpc error: code = Unknown desc = strconv.ParseFloat: parsing \"value\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.ParseFloat: parsing \"value\": invalid syntax",
		},
		testcaseName: "TestHIncrByFloatFieldNotFloat",
	}

	cmdMetaData.cmdTest(t)
}

func TestHDecrByKey(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hdecr-by-key field 2",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHDecrByKey",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HDecrby for-hdecr-by-key field 1",
		expectResult: []string{
			"1",
		},
		testcaseName: "TestHDecrByKey",
	}

	cmdMetaData.cmdTest(t)
}

func TestHDecrByKeyNotExists(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HDecrby for-hdecr-by-key-not-exists field 1",
		expectResult: []string{
			"0",
		},
		testcaseName: "TestHDecrByKeyNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHDecrByFieldNotExists(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hdecr-by-field-not-exists field 2",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHDecrByFieldNotExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HDecrby for-hdecr-by-field-not-exists field-not-exists 1",
		expectResult: []string{
			"0",
		},
		testcaseName: "TestHDecrByFieldNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHDecrByFieldNotInt(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hdecr-by-field-not-int field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHDecrByFieldNotInt",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HDecrby for-hdecr-by-field-not-int field 1",
		expectResult: []string{
			"HDecrBy error:  rpc error: code = Unknown desc = strconv.ParseInt: parsing \"value\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.ParseInt: parsing \"value\": invalid syntax",
		},
		testcaseName: "TestHDecrByFieldNotInt",
	}

	cmdMetaData.cmdTest(t)
}

func TestHMove(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hmove-source source source",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHMove",
	}

	putCmdMetaData.cmdTest(t)

	putCmdMetaData = &cmdMeta{
		cmd: "HSet for-hmove-dest dest dest",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHMove",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HMove for-hmove-source source for-hmove-dest",
		expectResult: []string{
			"Field moved successfully",
		},
		testcaseName: "TestHMove",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := &cmdMeta{
		cmd: "HGet for-hmove-source source-field",
		expectResult: []string{
			"get data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestHMove",
	}

	getCmdMetaData.cmdTest(t)

	getCmdMetaData = &cmdMeta{
		cmd: "HGet for-hmove-dest source",
		expectResult: []string{
			"source",
		},
		testcaseName: "TestHMove",
	}

	getCmdMetaData.cmdTest(t)
}

func TestHSetNx(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HSetnx for-hset-nx key value",
		expectResult: []string{
			"Field set successfully",
		},
		testcaseName: "TestHSetNx",
	}

	cmdMetaData.cmdTest(t)
}

func TestHSetNxFieldExists(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-hset-nx-field-exists key value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHSetNxFieldExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HSetnx for-hset-nx-field-exists key value",
		expectResult: []string{
			"HSetNX error:  HSet failed",
			"error: HSet failed",
		},
		testcaseName: "TestHSetNxFieldExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestHType(t *testing.T) {
	putCmdMetaData := &cmdMeta{
		cmd: "HSet for-htype-data field value",
		expectResult: []string{
			"put data success",
		},
		testcaseName: "TestHType",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := &cmdMeta{
		cmd: "HType for-htype-data field",
		expectResult: []string{
			"Type of field field in hash for-htype-data is hash",
		},
		testcaseName: "TestHType",
	}

	cmdMetaData.cmdTest(t)
}

func TestHTypeKeyNotExist(t *testing.T) {
	cmdMetaData := &cmdMeta{
		cmd: "HType for-htype-key-not-exist field",
		expectResult: []string{
			"Field or Key does not exist in the hash",
		},
		testcaseName: "TestHTypeKeyNotExist",
	}

	cmdMetaData.cmdTest(t)
}
