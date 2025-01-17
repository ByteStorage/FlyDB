package engine

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ByteStorage/FlyDB/lib/logger"
)

const (
	PutOperationName  = "put"
	DelOperationName  = "delete"
	LogConsoleEncoder = "console"
	engineDBLog       = "engine/db.go"
	logMetaCount      = 5
	logTimeLayout     = "2006-01-02T15:04:05.000Z0700"
	infoLogLevel      = "INFO"
	errorLogLevel     = "ERROR"
)

var (
	ErrLogIsEmpty                  = errors.New("log is empty")
	ErrLogMetaIsNotEnough          = errors.New("log meta count is not enough")
	ErrLogNotContainOperationField = errors.New("log not contain operation field")
	ErrLogUnSupportOperation       = errors.New("unsupported operation")
	ErrNotSupportEncoder           = errors.New("not support encoder")
)

var (
	supportEncoder = []string{LogConsoleEncoder}
)

type Operation struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func NewPutOperation(key, value []byte) *Operation {
	return &Operation{
		Name:  PutOperationName,
		Key:   string(key),
		Value: string(value),
	}
}

func NewDeleteOperation(key []byte) *Operation {
	return &Operation{
		Name: DelOperationName,
		Key:  string(key),
	}
}

type OperationLogHandlerOption func(*OperationLogHandler)

func WithLogConsoleEncoder() OperationLogHandlerOption {
	return func(o *OperationLogHandler) {
		o.logEncoder = LogConsoleEncoder
	}
}

func WithDB(db *DB) OperationLogHandlerOption {
	return func(o *OperationLogHandler) {
		o.db = db
	}
}

type OperationLogHandler struct {
	logEncoder    string
	logFilePath   string
	logLevel      string
	db            *DB
	logLinesChan  chan string
	operationChan chan *Operation
	errorChan     chan error
}

func defaultOperationLogHandler() *OperationLogHandler {
	return &OperationLogHandler{
		logLinesChan:  make(chan string, 1024),
		operationChan: make(chan *Operation, 1024),
		errorChan:     make(chan error, 1024),
	}
}

func NewOperationLogHandler(optionList ...OperationLogHandlerOption) *OperationLogHandler {
	handler := defaultOperationLogHandler()

	for _, option := range optionList {
		option(handler)
	}

	return handler
}

func (o *OperationLogHandler) readLog(logFilePath string) error {
	var (
		file    *os.File
		err     error
		scanner *bufio.Scanner
	)

	if file, err = os.Open(logFilePath); err != nil {
		return err
	}

	defer file.Close()

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		o.logLinesChan <- scanner.Text()
	}

	return scanner.Err()
}

type logMeta struct {
	logTime   time.Time
	logLevel  string
	file      string
	Operation *Operation `json:"operation"`
}

type operationErr struct {
	log string
	err error
	op  *Operation
}

func (o *operationErr) Error() string {
	operationErrStringBuilder := strings.Builder{}
	if len(o.log) > 0 {
		operationErrStringBuilder.WriteString("log: ")
		operationErrStringBuilder.WriteString(o.log)
		operationErrStringBuilder.WriteString("\t")
	}
	if o.err != nil {
		operationErrStringBuilder.WriteString("err: ")
		operationErrStringBuilder.WriteString(o.err.Error())
		operationErrStringBuilder.WriteString("\t")
	}
	if o.op != nil {
		operationErrStringBuilder.WriteString("operation: ")
		operationErrStringBuilder.WriteString(fmt.Sprintf("%+v", o.op))
	}

	return operationErrStringBuilder.String()
}

var _ error = &operationErr{}

func (o *OperationLogHandler) parseLogTime(logTime string) (time.Time, error) {
	var (
		parsedTime time.Time
		err        error
	)

	if parsedTime, err = time.Parse(logTimeLayout, logTime); err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

// decodeLogConsoleMeta decode log from console encoder
// 2024-08-11T20:23:05.599+0800	INFO	engine/db.go:68	open db	{"options": {"DirPath":"./data","DataFileSize":268435456,"SyncWrite":false,"IndexType":2,"FIOType":3}}
// time: 2024-08-11T20:23:05.599+0800
// level: INFO
// file: engine/db.go:68
// message: open db
// fields: {"options": {"DirPath":"./data","DataFileSize":268435456,"SyncWrite":false,"IndexType":2,"FIOType":3}}
func (o *OperationLogHandler) decodeLogConsoleMeta(log string) (*logMeta, error) {
	if len(log) == 0 {
		return nil, &operationErr{
			err: ErrLogIsEmpty,
		}
	}

	var (
		logSplit    = strings.Split(log, "\t")
		logMetaData = &logMeta{}
		err         error
	)

	if len(logSplit) < logMetaCount {
		return nil, &operationErr{
			log: log,
			err: ErrLogMetaIsNotEnough,
		}
	}

	if logMetaData.logTime, err = o.parseLogTime(logSplit[0]); err != nil {
		return nil, &operationErr{
			log: log,
			err: err,
		}
	}

	logMetaData.logLevel = logSplit[1]
	logMetaData.file = logSplit[2]
	if err = json.Unmarshal([]byte(logSplit[len(logSplit)-1]), &logMetaData); err != nil {
		return nil, &operationErr{
			log: log,
			err: err,
		}
	}

	if logMetaData.Operation == nil {
		return nil, &operationErr{
			log: log,
			err: ErrLogNotContainOperationField,
		}
	}

	return logMetaData, nil
}

func (o *OperationLogHandler) checkLogConsoleMeta(logMetaData *logMeta, start time.Time, end time.Time) bool {
	if logMetaData.logTime.Before(start) || logMetaData.logTime.After(end) {
		return false
	}

	if !strings.Contains(logMetaData.file, engineDBLog) {
		return false
	}

	if logMetaData.logLevel != o.logLevel {
		return false
	}

	return true
}

func (o *OperationLogHandler) decodeLogConsoleEncode(start time.Time, end time.Time) {
	for line := range o.logLinesChan {
		var (
			logMetaData *logMeta
			err         error
		)

		if logMetaData, err = o.decodeLogConsoleMeta(line); err != nil {
			o.errorChan <- err
			continue
		}

		if !o.checkLogConsoleMeta(logMetaData, start, end) {
			continue
		}

		o.operationChan <- logMetaData.Operation
	}
}

func (o *OperationLogHandler) restoreOperation() {
	for operation := range o.operationChan {
		var (
			err error
		)
		switch operation.Name {
		case PutOperationName:
			if err = o.db.Put([]byte(operation.Key), []byte(operation.Value)); err != nil {
				o.errorChan <- &operationErr{
					err: err,
					op:  operation,
				}
			}
		case DelOperationName:
			if err = o.db.Delete([]byte(operation.Key)); err != nil {
				o.errorChan <- &operationErr{
					err: err,
					op:  operation,
				}
			}
		default:
			o.errorChan <- &operationErr{
				err: ErrLogUnSupportOperation,
				op:  operation,
			}
		}
	}
}

func (o *OperationLogHandler) combinedError() error {
	var (
		combinedErr error
	)

	for err := range o.errorChan {
		if combinedErr == nil {
			combinedErr = fmt.Errorf("restore operation error")
		}
		combinedErr = fmt.Errorf("%w\n%v", combinedErr, err)
	}

	return combinedErr
}

func (o *OperationLogHandler) checkEncoderSupport() bool {
	for _, encoder := range supportEncoder {
		if o.logEncoder == encoder {
			return true
		}
	}

	return false
}

func (o *OperationLogHandler) restoreLogOperation(start time.Time, end time.Time) error {
	var (
		waitGroup = &sync.WaitGroup{}
	)

	if !o.checkEncoderSupport() {
		return ErrNotSupportEncoder
	}

	waitGroup.Add(3)

	go func() {
		defer waitGroup.Done()
		if err := o.readLog(o.logFilePath); err != nil {
			o.errorChan <- err
		}
		close(o.logLinesChan)
	}()

	go func() {
		defer waitGroup.Done()
		o.decodeLogConsoleEncode(start, end)
		close(o.operationChan)
	}()

	go func() {
		defer waitGroup.Done()
		o.restoreOperation()
	}()

	go func() {
		waitGroup.Wait()
		close(o.errorChan)
	}()

	return o.combinedError()
}

// RestoreWithTime will restore operation log with time range [start, end]
func (o *OperationLogHandler) RestoreWithTime(start time.Time, end time.Time) error {
	o.logFilePath = logger.LogLocation
	o.logLevel = infoLogLevel
	return o.restoreLogOperation(start, end)
}

// RestoreAfterStart will restore operation log after start time
func (o *OperationLogHandler) RestoreAfterStart(start time.Time) error {
	return o.RestoreWithTime(start, time.Now().Add(time.Hour*24*365))
}

// RestoreBeforeEnd will restore operation log before end time
func (o *OperationLogHandler) RestoreBeforeEnd(end time.Time) error {
	return o.RestoreWithTime(time.Time{}, end)
}

// Restore will restore operation log with all time range
func (o *OperationLogHandler) Restore() error {
	return o.RestoreWithTime(time.Time{}, time.Now().Add(time.Hour*24*365))
}

// ReExecuteWithTime re-execute operation log with time range [start, end]
func (o *OperationLogHandler) ReExecuteWithTime(start time.Time, end time.Time) error {
	o.logFilePath = logger.ErrorLogLocation
	o.logLevel = errorLogLevel
	return o.restoreLogOperation(start, end)
}

// ReExecuteAfterStart re-execute operation log after start time
func (o *OperationLogHandler) ReExecuteAfterStart(start time.Time) error {
	return o.ReExecuteWithTime(start, time.Now().Add(time.Hour*24*365))
}

// ReExecuteBeforeEnd re-execute operation log before end time
func (o *OperationLogHandler) ReExecuteBeforeEnd(end time.Time) error {
	return o.ReExecuteWithTime(time.Time{}, end)
}

// ReExecute re-execute operation log with all time range
func (o *OperationLogHandler) ReExecute() error {
	return o.ReExecuteWithTime(time.Time{}, time.Now().Add(time.Hour*24*365))
}
