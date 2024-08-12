package engine

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/lib/logger"
)

func Test_parseLogTime(t *testing.T) {
	operationLogHandler := NewOperationLogHandler(WithLogConsoleEncoder())

	testcase := []struct {
		name             string
		timeString       string
		expectTimeYear   int
		expectTimeMonth  int
		expectTimeHour   int
		expectTimeMin    int
		expectTimeSecond int
		expectErr        bool
	}{
		{
			name:             "parse time",
			timeString:       "2021-08-01T15:04:05.999+0800",
			expectTimeYear:   2021,
			expectTimeMonth:  8,
			expectTimeHour:   15,
			expectTimeMin:    4,
			expectTimeSecond: 5,
			expectErr:        false,
		},
		{
			name:             "parse time",
			timeString:       "2024-08-11T20:23:07.877+0800",
			expectTimeYear:   2024,
			expectTimeMonth:  8,
			expectTimeHour:   20,
			expectTimeMin:    23,
			expectTimeSecond: 7,
			expectErr:        false,
		},
		{
			name:       "bad time",
			timeString: "11-08-2024 20:23:05",
			expectErr:  true,
		},
		{
			name:       "empty time",
			timeString: "",
			expectErr:  true,
		},
	}

	for _, testcaseData := range testcase {
		t.Run(testcaseData.name, func(t *testing.T) {
			parsedTime, err := operationLogHandler.parseLogTime(testcaseData.timeString)
			assert.Equal(t, testcaseData.expectErr, err != nil)
			if !testcaseData.expectErr {
				assert.Equal(t, testcaseData.expectTimeYear, parsedTime.Year())
				assert.Equal(t, testcaseData.expectTimeMonth, int(parsedTime.Month()))
				assert.Equal(t, testcaseData.expectTimeHour, parsedTime.Hour())
				assert.Equal(t, testcaseData.expectTimeMin, parsedTime.Minute())
				assert.Equal(t, testcaseData.expectTimeSecond, parsedTime.Second())
			}
		})
	}
}

func Test_decodeLogConsoleMeta(t *testing.T) {
	operationLogHandler := NewOperationLogHandler(WithLogConsoleEncoder())

	testcase := []struct {
		name          string
		log           string
		expectErr     bool
		expectLogMeta *logMeta
	}{
		{
			name: "decode put log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput error\t{\"error\": \"truncate data\\\\000000000.data: Access is denied.\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
			expectLogMeta: &logMeta{
				logLevel: "INFO",
				file:     "engine/db.go:171",
				Operation: &Operation{
					Name:  "put",
					Key:   "test",
					Value: "test",
				},
			},
		},
		{
			name: "decode delete log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tdelete error\t{\"error\": \"truncate data\\\\000000000.data: Access is denied.\", " +
				"\"operation\": {\"name\":\"delete\",\"key\":\"test\"}}\n",
			expectLogMeta: &logMeta{
				logLevel: "INFO",
				file:     "engine/db.go:171",
				Operation: &Operation{
					Name: "delete",
					Key:  "test",
				},
			},
			expectErr: false,
		},
		{
			name:          "empty log",
			log:           "",
			expectErr:     true,
			expectLogMeta: nil,
		},
		{
			name: "bad log log meta is not enough",
			log: "2021-08-01T15:04:05.999+0800\tINFO\tengine/db.go:68	open db",
			expectErr:     true,
			expectLogMeta: nil,
		},
		{
			name: "bad log without operation",
			log: "2021-08-01T15:04:05.999+0800\tINFO\tengine/db.go:68	open db	{\"options\": {\"DirPath\":\"." +
				"/data\",\"DataFileSize\":268435456,\"SyncWrite\":false,\"IndexType\":2,\"FIOType\":3}}",
			expectErr:     true,
			expectLogMeta: nil,
		},
		{
			name: "bad log parse time error",
			log: "15:04:05.999+0800\tINFO\tengine/db.go:68	open db	{\"options\": {\"DirPath\":\".",
			expectErr:     true,
			expectLogMeta: nil,
		},
		{
			name: "bad log log field unmarshal error",
			log: "2021-08-01T15:04:05.999+0800\tINFO\tengine/db.go:68	open db	errorField",
			expectErr:     true,
			expectLogMeta: nil,
		},
	}

	for _, testcaseData := range testcase {
		t.Run(testcaseData.name, func(t *testing.T) {
			logMetaData, err := operationLogHandler.decodeLogConsoleMeta(testcaseData.log)
			assert.Equal(t, testcaseData.expectErr, err != nil)
			if !testcaseData.expectErr {
				assert.Equal(t, testcaseData.expectLogMeta.logLevel, logMetaData.logLevel)
				assert.Equal(t, testcaseData.expectLogMeta.file, logMetaData.file)
				assert.Equal(t, testcaseData.expectLogMeta.Operation.Name, logMetaData.Operation.Name)
				assert.Equal(t, testcaseData.expectLogMeta.Operation.Key, logMetaData.Operation.Key)
				assert.Equal(t, testcaseData.expectLogMeta.Operation.Value, logMetaData.Operation.Value)
			} else {
				assert.Nil(t, logMetaData)
			}
		})
	}
}

func Test_readLog(t *testing.T) {
	mockLogs := []string{
		"log1\n",
		"log2\n",
		"log3\n",
		"log4\n",
		"log5\n",
		"log6\n",
	}

	tmpFile, err := ioutil.TempFile("", "test")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err = os.Remove(tmpFile.Name()); err != nil {
			t.Fatal(err)
		}
	}()

	for _, log := range mockLogs {
		if _, err = tmpFile.WriteString(log); err != nil {
			t.Fatal(err)
		}
	}

	if _, err = tmpFile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	operationLogHandler := NewOperationLogHandler(WithLogConsoleEncoder())

	if err = operationLogHandler.readLog(tmpFile.Name()); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(mockLogs), len(operationLogHandler.logLinesChan))

	for _, log := range mockLogs {
		assert.Equal(t, log, <-operationLogHandler.logLinesChan+"\n")
	}
}

func Test_readLogError(t *testing.T) {
	operationLogHandler := NewOperationLogHandler(WithLogConsoleEncoder())

	err := operationLogHandler.readLog("not_exist_file")

	assert.Error(t, err)
	assert.Equal(t, 0, len(operationLogHandler.logLinesChan))
}

func TestDecodeLogConsoleEncode(t *testing.T) {
	operationLogHandler := NewOperationLogHandler(WithLogConsoleEncoder())
	operationLogHandler.logLevel = infoLogLevel
	expectOperationList := make([]*Operation, 0)

	testcase := []struct {
		name            string
		log             string
		expectErr       bool
		expectOperation *Operation
	}{
		{
			name: "decode put log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput error\t{\"error\": \"truncate data\\\\000000000.data: Access is denied.\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
			expectOperation: &Operation{
				Name:  "put",
				Key:   "test",
				Value: "test",
			},
		},
		{
			name: "decode delete log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tdelete error\t{\"error\": \"truncate data\\\\000000000.data: Access is denied.\", " +
				"\"operation\": {\"name\":\"delete\",\"key\":\"test\"}}\n",
			expectOperation: &Operation{
				Name: "delete",
				Key:  "test",
			},
		},
		{
			name: "decode put log file name not engin/db.go",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tnotEngine/db." +
				"go:171\tput error\t{\"error\": \"truncate data\\\\000000000.data: Access is denied.\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
			expectOperation: nil,
		},
		{
			name: "decode time is before start",
			log: "1999-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput error\t{\"error\": \"truncate data\\\\000000000.data: Access is denied.\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
			expectOperation: nil,
		},
		{
			name: "decode time is after end",
			log: "2800-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput error\t{\"error\": \"truncate data\\\\000000000.data: Access is denied.\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
			expectOperation: nil,
		},
		{
			name:            "decode log console error",
			log:             "",
			expectOperation: nil,
		},
	}

	for _, testcaseData := range testcase {
		if testcaseData.expectOperation != nil {
			expectOperationList = append(expectOperationList, testcaseData.expectOperation)
		}

		operationLogHandler.logLinesChan <- testcaseData.log
	}
	close(operationLogHandler.logLinesChan)

	location, _ := time.LoadLocation("Asia/Shanghai")
	operationLogHandler.decodeLogConsoleEncode(time.Date(2024, 8, 11, 20, 20, 0, 0, location), time.Now().Add(time.Hour))
	close(operationLogHandler.operationChan)
	assert.Equal(t, len(expectOperationList), len(operationLogHandler.operationChan))

	for getOperation := range operationLogHandler.operationChan {
		assert.Equal(t, expectOperationList[0].Name, getOperation.Name)
		assert.Equal(t, expectOperationList[0].Key, getOperation.Key)
		assert.Equal(t, expectOperationList[0].Value, getOperation.Value)
		expectOperationList = expectOperationList[1:]
	}
}

func TestRestoreOperation(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("./", "flydb")
	opts.DirPath = dir
	db, err := NewDB(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDB(db)

	operationLogHandler := NewOperationLogHandler(WithLogConsoleEncoder(), WithDB(db))

	testcases := []struct {
		name      string
		operation *Operation
	}{
		{
			name: "put operation",
			operation: &Operation{
				Name:  "put",
				Key:   "put_key",
				Value: "put_value",
			},
		},
		{
			name: "delete operation",
			operation: &Operation{
				Name: "delete",
				Key:  "put_key",
			},
		},
		{
			name: "unknown operation",
			operation: &Operation{
				Name: "unknown",
				Key:  "put_key",
			},
		},
	}

	for _, testcaseData := range testcases {
		operationLogHandler.operationChan <- testcaseData.operation
	}

	close(operationLogHandler.operationChan)

	operationLogHandler.restoreOperation()

	value, err := db.Get([]byte(testcases[0].operation.Key))
	assert.Error(t, err)
	assert.Nil(t, value)
}

func TestOperationLogHandler_Restore(t *testing.T) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("./", "flydb")
	opts.DirPath = dir
	db, err := NewDB(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer destroyDB(db)
	defer func() {
		if err = os.RemoveAll(filepath.Dir(logger.LogLocation)); err != nil {
			t.Fatal(err)
		}
	}()

	if err = db.Put([]byte("need_to_remove"), []byte("need to remove")); err != nil {
		t.Fatal(err)
	}

	operationLogHandler := NewOperationLogHandler(WithLogConsoleEncoder(), WithDB(db))

	testcases := []struct {
		name string
		log  string
	}{
		{
			name: "success put log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test-success\",\"value\":\"test\"}}\n",
		},
		{
			name: "success delete log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tdelete mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"delete\",\"key\":\"test-success\"}}\n",
		},
		{
			name: "error log level log",
			log: "2024-08-11T20:23:07.877+0800\tERROR\tengine/db." +
				"go:171\tdelete mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"delete\",\"key\":\"test\"}}\n",
		},
		{
			name: "time is before start",
			log: "1999-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
		},
		{
			name: "time is after end",
			log: "2800-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
		},
		{
			name: "file name is not engine/db.go",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tnotEngine/db." +
				"go:171\tput mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
		},
		{
			name: "parse time error",
			log: "15:04:05.999+0800\tINFO\tengine/db.go:68	error parse time" +
				"\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"test\",\"value\":\"test\"}}\n",
		},
		{
			name: "log field length is not enough",
			log: "2021-08-01T15:04:05.999+0800\tINFO\tengine/db.go:68	open db\n",
		},
		{
			name: "success put new log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tput mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"put\",\"key\":\"new_key\",\"value\":\"test\"}}\n",
		},
		{
			name: "delete need to remove log",
			log: "2024-08-11T20:23:07.877+0800\tINFO\tengine/db." +
				"go:171\tdelete mock\t{\"mock\": \"mock log\", " +
				"\"operation\": {\"name\":\"delete\",\"key\":\"need_to_remove\"}}\n",
		},
	}

	logFile, err := os.OpenFile(logger.LogLocation, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal(err)
	}

	for _, testcaseData := range testcases {
		if _, err = logFile.WriteString(testcaseData.log); err != nil {
			t.Fatal(err)
		}
	}

	location, _ := time.LoadLocation("Asia/Shanghai")
	err = operationLogHandler.RestoreWithTime(time.Date(2024, 8, 11, 20, 20, 0, 0, location), time.Date(2024, 8, 11, 23, 20,
		0, 0, location))
	assert.Error(t, err)
	t.Log(err)

	value, err := db.Get([]byte("test"))
	assert.Error(t, err)
	assert.Nil(t, value)

	value, err = db.Get([]byte("new_key"))
	assert.NoError(t, err)
	assert.Equal(t, "test", string(value))

	value, err = db.Get([]byte("need_to_remove"))
	assert.Error(t, err)
	assert.Nil(t, value)
}
