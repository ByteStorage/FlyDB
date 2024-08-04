package client

import (
	"testing"
)

func TestStringPut(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "put test test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringPut",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringPutEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "put test",
		expectResult: []string{"key or value is empty"},
		testcaseName: "TestStringPutEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringGet(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-get test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringGet",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "get for-get",
		expectResult: []string{"test"},
		testcaseName: "TestStringGet",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringGetEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "get",
		expectResult: []string{"key is empty"},
		testcaseName: "TestStringGetEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringGetNotExist(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "get not-exist",
		expectResult: []string{
			"get data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringGetNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringDel(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-del test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringDel",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "delete for-del",
		expectResult: []string{"delete key success"},
		testcaseName: "TestStringDel",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd: "get for-del",
		expectResult: []string{
			"get data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringDel",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringDelEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "delete",
		expectResult: []string{"key is empty"},
		testcaseName: "TestStringDelEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringDelNotExist(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "delete not-exist",
		expectResult: []string{"delete key success"},
		testcaseName: "TestStringDelNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringType(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-type test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringType",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "type for-type",
		expectResult: []string{"Type: string"},
		testcaseName: "TestStringType",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringTypeEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "type",
		expectResult: []string{"key is empty"},
		testcaseName: "TestStringTypeEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringTypeNotExist(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "type not-exist",
		expectResult: []string{
			"get type error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringTypeNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringStrLen(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-strlen test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringStrLen",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "strlen for-strlen",
		expectResult: []string{"4"},
		testcaseName: "TestStringStrLen",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringStrLenEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "strlen",
		expectResult: []string{"key is empty"},
		testcaseName: "TestStringStrLenEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringStrLenNotExist(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "strlen not-exist",
		expectResult: []string{
			"get string length error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringStrLenNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringGetSet(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-getset test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringGetSet",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "getset for-getset test2",
		expectResult: []string{"test"},
		testcaseName: "TestStringGetSet",
	}

	cmdMetaData.cmdTest(t)

	cmdMetaData = cmdMeta{
		cmd:          "get for-getset",
		expectResult: []string{"test2"},
		testcaseName: "TestStringGetSet",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringGetSetEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "getset",
		expectResult: []string{"key or value is empty"},
		testcaseName: "TestStringGetSetEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringGetSetNotExist(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "getset not-exist test",
		expectResult: []string{
			"getset operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringGetSetNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringAppend(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-append test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringAppend",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "append for-append test2",
		expectResult: []string{"Append is successful"},
		testcaseName: "TestStringAppend",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-append",
		expectResult: []string{"testtest2"},
		testcaseName: "TestStringAppend",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringAppendEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "append test",
		expectResult: []string{"key or value is empty"},
		testcaseName: "TestStringAppendEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringAppendNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "append not-exist test",
		expectResult: []string{
			"Append operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringAppendNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringIncr(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-incr 1",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringIncr",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "incr for-incr",
		expectResult: []string{"Incr operation success"},
		testcaseName: "TestStringIncr",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-incr",
		expectResult: []string{"2"},
		testcaseName: "TestStringIncr",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringIncrEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "incr",
		expectResult: []string{"key is empty"},
		testcaseName: "TestStringIncrEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringIncrNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "incr not-exist",
		expectResult: []string{
			"incr operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringIncrNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringIncrNotNumber(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-incr not-number",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringIncrNotNumber",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd: "incr for-incr",
		expectResult: []string{
			"incr operation error:  rpc error: code = Unknown desc = strconv.Atoi: parsing \"not-number\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.Atoi: parsing \"not-number\": invalid syntax",
		},
		testcaseName: "TestStringIncrNotNumber",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringIncrBy(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-incrby 1",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringIncrBy",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "incrby for-incrby 2",
		expectResult: []string{"IncrBy operation success"},
		testcaseName: "TestStringIncrBy",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-incrby",
		expectResult: []string{"3"},
		testcaseName: "TestStringIncrBy",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringIncrByNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "incrby not-exist 2",
		expectResult: []string{
			"incrby operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringIncrByNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringIncrByNotNumber(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-incrby not-number",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringIncrByNotNumber",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd: "incrby for-incrby 1",
		expectResult: []string{
			"incrby operation error:  rpc error: code = Unknown desc = strconv." +
				"Atoi: parsing \"not-number\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.Atoi: parsing \"not-number\": invalid syntax",
		},
		testcaseName: "TestStringIncrByNotNumber",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringIncrByFloat(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-incrbyfloat 1.1",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringIncrByFloat",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "incrbyfloat for-incrbyfloat 2.2",
		expectResult: []string{"IncrByFloat operation success"},
		testcaseName: "TestStringIncrByFloat",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-incrbyfloat",
		expectResult: []string{"3.3"},
		testcaseName: "TestStringIncrByFloat",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringIncrByFloatNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "incrbyfloat not-exist 2.2",
		expectResult: []string{
			"incrbyfloat operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in" +
				" database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringIncrByFloatNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringIncrByFloatNotNumber(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-incrbyfloat not-number",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringIncrByFloatNotNumber",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd: "incrbyfloat for-incrbyfloat 1.1",
		expectResult: []string{
			"incrbyfloat operation error:  rpc error: code = Unknown desc = strconv." +
				"ParseFloat: parsing \"not-number\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.ParseFloat: parsing \"not-number\": invalid syntax",
		},
		testcaseName: "TestStringIncrByFloatNotNumber",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringDecr(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-decr 1",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringDecr",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "decr for-decr",
		expectResult: []string{"Decr operation success"},
		testcaseName: "TestStringDecr",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-decr",
		expectResult: []string{"0"},
		testcaseName: "TestStringDecr",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringDecrEmpty(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "decr",
		expectResult: []string{"key is empty"},
		testcaseName: "TestStringDecrEmpty",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringDecrNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "decr not-exist",
		expectResult: []string{
			"decr operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringDecrNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringDecrNotNumber(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-decr not-number",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringDecrNotNumber",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd: "decr for-decr",
		expectResult: []string{
			"decr operation error:  rpc error: code = Unknown desc = strconv." +
				"Atoi: parsing \"not-number\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.Atoi: parsing \"not-number\": invalid syntax",
		},
		testcaseName: "TestStringDecrNotNumber",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringDecrBy(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-decrby 3",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringDecrBy",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "decrby for-decrby 2",
		expectResult: []string{"DecrBy operation success"},
		testcaseName: "TestStringDecrBy",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-decrby",
		expectResult: []string{"1"},
		testcaseName: "TestStringDecrBy",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringDecrByNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "decrby not-exist 2",
		expectResult: []string{
			"decrby operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringDecrByNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringDecrByNotNumber(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-decrby not-number",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringDecrByNotNumber",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd: "decrby for-decrby 1",
		expectResult: []string{
			"decrby operation error:  rpc error: code = Unknown desc = strconv." +
				"Atoi: parsing \"not-number\": invalid syntax",
			"error: rpc error: code = Unknown desc = strconv.Atoi: parsing \"not-number\": invalid syntax",
		},
		testcaseName: "TestStringDecrByNotNumber",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringExists(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-exists test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "exists for-exists",
		expectResult: []string{"key is exist for-exists"},
		testcaseName: "TestStringExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringExistsNotExist(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "exists not-exist",
		expectResult: []string{"key is not exist not-exist"},
		testcaseName: "TestStringExistsNotExist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringExpire(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-string-expire test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringExpire",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "expire for-string-expire 10",
		expectResult: []string{"expire key success"},
		testcaseName: "TestStringExpire",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringExpireNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "expire not-exist 10",
		expectResult: []string{
			"expire operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringExpireNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringPersist(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-persist test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringPersist",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "persist for-persist",
		expectResult: []string{"key is Persist"},
		testcaseName: "TestStringPersist",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringPersistNotExists(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd: "persist not-exist",
		expectResult: []string{
			"persist operation error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in" +
				" database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringPersistNotExists",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringMGet(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-mget test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringMGet",
	}

	putCmdMetaData.cmdTest(t)

	putCmdMetaData = cmdMeta{
		cmd:          "put for-mget2 test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringMGet",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "mget for-mget for-mget2",
		expectResult: []string{"[test test]"},
		testcaseName: "TestStringMGet",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringMGetNotExists(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-mget test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringMGet",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd: "mget for-mget not-exist",
		expectResult: []string{
			"get data error:  rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
			"error: rpc error: code = Unknown desc = KeyNotFoundError : key is not found in database",
		},
		testcaseName: "TestStringMGet",
	}

	cmdMetaData.cmdTest(t)
}

func TestStringMSet(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-mset test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringMSet",
	}

	putCmdMetaData.cmdTest(t)

	putCmdMetaData = cmdMeta{
		cmd:          "put for-mset2 test2",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringMSet",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd:          "mset for-mset mset for-mset2 mset",
		expectResult: []string{"Data successfully set."},
		testcaseName: "TestStringMSet",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-mset",
		expectResult: []string{"mset"},
		testcaseName: "TestStringMSet",
	}

	getCmdMetaData.cmdTest(t)

	getCmdMetaData = cmdMeta{
		cmd:          "get for-mset2",
		expectResult: []string{"mset"},
		testcaseName: "TestStringMSet",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringMSetNX(t *testing.T) {
	cmdMetaData := cmdMeta{
		cmd:          "msetnx for-msetnx msetnx for-msetnx2 msetnx",
		expectResult: []string{"Data successfully set."},
		testcaseName: "TestStringMSetNX",
	}

	cmdMetaData.cmdTest(t)

	getCmdMetaData := cmdMeta{
		cmd:          "get for-msetnx",
		expectResult: []string{"msetnx"},
		testcaseName: "TestStringMSetNX",
	}

	getCmdMetaData.cmdTest(t)

	getCmdMetaData = cmdMeta{
		cmd:          "get for-msetnx2",
		expectResult: []string{"msetnx"},
		testcaseName: "TestStringMSetNX",
	}

	getCmdMetaData.cmdTest(t)
}

func TestStringMSetNXExists(t *testing.T) {
	putCmdMetaData := cmdMeta{
		cmd:          "put for-msetnxexists test",
		expectResult: []string{"put data success"},
		testcaseName: "TestStringMSetNXExists",
	}

	putCmdMetaData.cmdTest(t)

	cmdMetaData := cmdMeta{
		cmd: "msetnx for-msetnxexists msetnx",
		expectResult: []string{
			"set data error: MSetNX failed  any key has existed",
			"MSetNX failed  any key has existed",
		},
		testcaseName: "TestStringMSetNXExists",
	}

	cmdMetaData.cmdTest(t)
}
