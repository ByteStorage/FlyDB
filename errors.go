package flydb

import "errors"

var (
	ErrKeyIsEmpty                    = errors.New("The key is empty")
	ErrIndexUpdataFailed             = errors.New("failed to update index")
	ErrKeyNotFound                   = errors.New("ket not found in database")
	ErrDataFailNotFound              = errors.New("data file is not found")
	ErrDataDirctoryCorrupted         = errors.New("the databases directory maybe corrupted")
	ErrOptionDirPathIsEmpty          = errors.New("database dir path is empty")
	ErrOptionDataFileSizeNotPositive = errors.New("database data file size must be greater than 0")
)
