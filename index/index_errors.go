package index

import "errors"

var (
	ErrOpenBPTreeFailed   = errors.New("OpenBPTreeFailedError : failed to open bptree")
	ErrCreateBucketFailed = errors.New("CreateBucketFailedError : failed to create bucket in bptree")
	ErrPutValueFailed     = errors.New("PutValueFailedError : failed to put value in bptree")
	ErrGetValueFailed     = errors.New("GetValueFailedError : failed to get value in bptree")
	ErrDeleteValueFailed  = errors.New("DeleteValueFailedError : failed to delete value in bptree")
	ErrGetIndexSizeFailed = errors.New("GetIndexSizeFailedError : failed to get index size in bptree")
)
