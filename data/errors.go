package data

import "errors"

var (
	ErrInvalidCRC = errors.New("InvalidCrcError : invalid crc value, log record maybe corrupted")
)
