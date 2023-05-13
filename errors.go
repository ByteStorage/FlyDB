package flydb

import "errors"

var (
	ErrKeyIsEmpty             = errors.New("KeyEmptyError : the key is empty")
	ErrIndexUpdateFailed      = errors.New("IndexUpdateFailError : failed to update index")
	ErrKeyNotFound            = errors.New("KeyNotFoundError : key is not found in database")
	ErrDataFailNotFound       = errors.New("DataFailNotFoundError : data file is not found")
	ErrDataDirectoryCorrupted = errors.New("DataDirectoryCorruptedError : the databases directory maybe corrupted")
	ErrExceedMaxBatchNum      = errors.New("ExceedMaxBatchNumError : exceed the max batch num")
	ErrMergeIsProgress        = errors.New("MergeIsProgressError : merge is in progress, try again later")

	ErrOptionDirPathIsEmpty          = errors.New("OptionDirPathError : database dir path is empty")
	ErrOptionDataFileSizeNotPositive = errors.New("OptionDataFileSizeError : database data file size must be greater than 0")

	ErrOpenBPTreeFailed   = errors.New("OpenBPTreeFailedError : failed to open bptree")
	ErrCreateBucketFailed = errors.New("CreateBucketFailedError : failed to create bucket in bptree")
)
