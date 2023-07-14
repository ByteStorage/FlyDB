package structure

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
)

// Due to the complexity of each operation is at least O(n)
// So we can directly use slice to implement the list at the bottom level
// If the implementation of the db is improved later, we need to switch to a linked list
type ListStructure struct {
	db *engine.DB
}

// NewListStructure returns a new ListStructure
// It will return a nil ListStructure if the database cannot be opened
// or the database cannot be created
// The database will be created if it does not exist
func NewListStructure(options config.Options) (*ListStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &ListStructure{db: db}, nil
}

// LPush adds a value to the left of the list corresponding to the key
// If the key does not exist, it will create the key
func (l *ListStructure) LPush(key []byte, value []byte) error {
	// Check if value is empty
	if value == nil {
		return ErrInvalidValue
	}

	// Get the list
	lstArr, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	// Use new slice space, add data at the head, and append all original data afterwards
	tmpArr := make([][]byte, 1+len(lstArr))
	tmpArr[0] = value
	copy(tmpArr[1:], lstArr)
	lstArr = tmpArr

	// Store to db
	return l.setListToDB(key, lstArr)
}

// LPushs adds one or more values to the left of the list corresponding to the key
// If the key does not exist, it will create the key
func (l *ListStructure) LPushs(key []byte, values ...[]byte) error {
	// Check if values are valid
	if len(values) == 0 {
		return ErrInvalidValue
	}
	for _, v := range values {
		if len(v) == 0 {
			return ErrInvalidValue
		}
	}

	// Get the list
	lstArr, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	// Use new slice space, add data at the head, and append all original data afterwards
	tmpArr := make([][]byte, len(values)+len(lstArr))
	copy(tmpArr[:len(values)], values)
	copy(tmpArr[len(values):], lstArr)
	lstArr = tmpArr

	// Store to db
	return l.setListToDB(key, lstArr)
}

// RPush adds a value to the right of the list corresponding to the key
// If the key does not exist, it will create the key
func (l *ListStructure) RPush(key []byte, value []byte) error {
	// Check if value is empty
	if value == nil {
		return ErrInvalidValue
	}

	// Get the list
	lstArr, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	// Append the new data to the end
	lstArr = append(lstArr, value)

	// Store to db
	return l.setListToDB(key, lstArr)
}

// RPushs appends one or more values to the right side of a list associated with a key.
// If the key does not exist, it will be created.
func (l *ListStructure) RPushs(key []byte, values ...[]byte) error {
	// Check if values are valid
	if len(values) == 0 {
		return ErrInvalidValue
	}
	for _, v := range values {
		if len(v) == 0 {
			return ErrInvalidValue
		}
	}

	// Get the list
	lstArr, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	// Append new values to the end
	lstArr = append(lstArr, values...)

	// Store in the database
	return l.setListToDB(key, lstArr)
}

// LPop returns and removes the leftmost value of a list associated with a key.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
func (l *ListStructure) LPop(key []byte) ([]byte, error) {
	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if len(lstArr) == 0 {
		return nil, ErrListEmpty
	}

	leftData := lstArr[0]
	lstArr = lstArr[1:]

	// Store in the database
	return leftData, l.setListToDB(key, lstArr)
}

// RPop returns and removes the rightmost value of a list associated with a key.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
func (l *ListStructure) RPop(key []byte) ([]byte, error) {
	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if len(lstArr) == 0 {
		return nil, ErrListEmpty
	}

	rightData := lstArr[len(lstArr)-1]
	lstArr = lstArr[:len(lstArr)-1]

	// Store in the database
	return rightData, l.setListToDB(key, lstArr)
}

// LRange returns a range of elements from a list associated with a key.
// The range is inclusive, including both the start and stop indices.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
// Negative indices can be used, where -1 represents the last element of the list,
// -2 represents the second last element, and so on.
func (l *ListStructure) LRange(key []byte, start int, stop int) ([][]byte, error) {
	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if len(lstArr) == 0 {
		return nil, ErrListEmpty
	}

	lstLen := len(lstArr)

	// Calculate the correct indices
	start = (start%lstLen + lstLen) % lstLen
	stop = (stop%lstLen + lstLen) % lstLen

	// Return empty if the range length is less than 1
	if stop-start+1 < 1 {
		return nil, nil
	}

	return lstArr[start : stop+1], nil
}

// LLen returns the size of a list associated with a key.
// If the key does not exist, an error is returned.
func (l *ListStructure) LLen(key []byte) (int, error) {
	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return 0, err
	}

	return len(lstArr), nil
}

// LRem removes elements from a list associated with a key based on the count and value parameters.
// The count can have the following values:
// count > 0: Remove count occurrences of the value from the beginning of the list.
// count < 0: Remove count occurrences of the value from the end of the list.
// count = 0: Remove all occurrences of the value from the list.
// If the key does not exist, an error is returned.
func (l *ListStructure) LRem(key []byte, count int, value []byte) error {
	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return err
	}

	// Store whether an element is removed
	isRemoved := make([]bool, len(lstArr))

	num := 0
	// Process different counts to calculate the isRemoved array
	if count > 0 {
		for i := 0; i < len(lstArr) && num < count; i++ {
			if bytes.Equal(lstArr[i], value) {
				isRemoved[i] = true
				num++
			}
		}
	} else if count < 0 {
		for i := len(lstArr) - 1; i >= 0 && num < -count; i-- {
			if bytes.Equal(lstArr[i], value) {
				isRemoved[i] = true
				num++
			}
		}
	} else {
		for i := 0; i < len(lstArr); i++ {
			if bytes.Equal(lstArr[i], value) {
				isRemoved[i] = true
				num++
			}
		}
	}

	// Create a new slice to store the list after removal
	tmpList := make([][]byte, 0, len(lstArr)-num)
	for i := 0; i < len(lstArr); i++ {
		if !isRemoved[i] {
			tmpList = append(tmpList, lstArr[i])
		}
	}
	lstArr = tmpList

	// Store in the database
	return l.setListToDB(key, lstArr)
}

// LSet sets the value of an element in a list associated with a key based on the index.
// If the index is out of range, an error is returned.
// If the list is empty, an error is returned.
func (l *ListStructure) LSet(key []byte, index int, value []byte) error {
	// Check if the value is valid
	if len(value) == 0 {
		return ErrInvalidValue
	}

	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return err
	}

	// Return error if the list is empty
	if len(lstArr) == 0 {
		return ErrListEmpty
	}

	// Check if the index is out of range
	if index < 0 || index >= len(lstArr) {
		return ErrIndexOutOfRange
	}

	lstArr[index] = value

	// Store in the database
	return l.setListToDB(key, lstArr)
}

// LTrim retains a range of elements in a list associated with a key.
// The range is inclusive, including both the start and stop indices.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
// Negative indices can be used, where -1 represents the last element of the list,
// -2 represents the second last element, and so on.
func (l *ListStructure) LTrim(key []byte, start int, stop int) error {
	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return err
	}

	if len(lstArr) == 0 {
		return ErrListEmpty
	}

	lstLen := len(lstArr)

	// Calculate the correct indices
	start = (start%lstLen + lstLen) % lstLen
	stop = (stop%lstLen + lstLen) % lstLen

	if stop-start+1 < 1 { // Empty the list if the range length is less than 1
		lstArr = make([][]byte, 0)
	} else {
		lstArr = lstArr[start : stop+1]
	}

	// Store in the database
	return l.setListToDB(key, lstArr)
}

// LIndex returns the value of an element in a list associated with a key based on the index.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
// Negative indices can be used, where -1 represents the last element of the list,
// -2 represents the second last element, and so on.
func (l *ListStructure) LIndex(key []byte, index int) ([]byte, error) {
	// Get the list
	lstArr, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if len(lstArr) == 0 {
		return nil, ErrListEmpty
	}

	lstLen := len(lstArr)

	// Calculate the correct index
	index = (index%lstLen + lstLen) % lstLen

	return lstArr[index], nil
}

// RPOPLPUSH removes the last element from one list and pushes it to another list.
// If the source list is empty, an error is returned.
// If the destination list is empty, it is created.
// Atomicity is not guaranteed.
func (l *ListStructure) RPOPLPUSH(source []byte, destination []byte) error {
	// Get the source list
	lstArr1, err := l.getListFromDB(source, false)
	if err != nil {
		return err
	}

	// Get the destination list
	lstArr2, err := l.getListFromDB(destination, true)
	if err != nil {
		return err
	}

	// Return error if the source list is empty
	if len(lstArr1) == 0 {
		return ErrListEmpty
	}

	// Insert lstArr1's last element at the beginning of lstArr2
	tmpLst := make([][]byte, len(lstArr2)+1)
	tmpLst[0] = lstArr1[len(lstArr1)-1]
	copy(tmpLst[1:], lstArr2)
	lstArr2 = tmpLst

	// Truncate lstArr1's last element
	lstArr1 = lstArr1[:len(lstArr1)-1]

	// Store in the database
	err = l.setListToDB(source, lstArr1)
	if err != nil {
		return err
	}
	return l.setListToDB(destination, lstArr2)
}

var (
	// ErrListEmpty is returned if the list is empty.
	ErrListEmpty = errors.New("Wrong operation: list is empty")
	// ErrIndexOutOfRange is returned if the index out of range.
	ErrIndexOutOfRange = errors.New("Wrong operation: index out of range")
)

// getListFromDB retrieves data from the database. When isKeyCanNotExist is true, it returns an empty slice if the key doesn't exist instead of an error.
func (l *ListStructure) getListFromDB(key []byte, isKeyCanNotExist bool) ([][]byte, error) {
	if isKeyCanNotExist {
		// Get data corresponding to the key from the database
		dbData, err := l.db.Get(key)

		// Since the key might not exist, we need to handle ErrKeyNotFound separately as it is a valid case
		if err != nil && err != _const.ErrKeyNotFound {
			return nil, err
		}

		// Deserialize the data into a list
		lstArr, err := l.decodeList(dbData)
		if err != nil {
			if len(dbData) != 0 {
				return nil, err
			} else {
				lstArr = make([][]byte, 0)
			}
		}
		return lstArr, nil
	} else {
		// Get data corresponding to the key from the database
		dbData, err := l.db.Get(key)
		if err != nil {
			return nil, err
		}

		// Deserialize the data into a list
		lstArr, err := l.decodeList(dbData)
		if err != nil {
			return nil, err
		}
		return lstArr, nil
	}
}

// setListToDB stores the data into the database.
func (l *ListStructure) setListToDB(key []byte, lstArr [][]byte) error {
	// Serialize into binary array
	encValue, err := l.encodeList(lstArr)
	if err != nil {
		return err
	}
	// Store in the database
	return l.db.Put(key, encValue)
}

// encodeList encodes the value
// format: [type][len1][value1][len2][value2]...
// len: variable number of bytes
// value: len bytes
func (l *ListStructure) encodeList(data [][]byte) ([]byte, error) {
	dataLen := len(data)

	// Calculate the required buffer space in advance
	bufMaxLen := 1
	for i := 0; i < dataLen; i++ {
		bufMaxLen += len(data[i]) + binary.MaxVarintLen64
	}
	buf := make([]byte, bufMaxLen)

	buf[0] = List

	bufIndex := 1

	for i := 0; i < dataLen; i++ {
		bufIndex += binary.PutVarint(buf[bufIndex:], int64(len(data[i])))
		bufIndex += copy(buf[bufIndex:], data[i])
	}

	return buf[:bufIndex], nil
}

// decodeList decodes the value
// format: [type][len1][value1][len2][value2]...
// len: variable number of bytes
// value: len bytes
func (l *ListStructure) decodeList(value []byte) ([][]byte, error) {
	// Check the length of the value
	if len(value) < 1 {
		return nil, ErrInvalidValue
	}

	// Check the type of the value
	if value[0] != List {
		return nil, ErrInvalidType
	}

	valueLen := len(value)

	nowIndex := 1

	lstLen := 0

	// Calculate the length of the list
	for nowIndex < valueLen {
		length, lenOfLen := binary.Varint(value[nowIndex:])
		nowIndex += lenOfLen + int(length)
		lstLen++
	}

	result := make([][]byte, 0, lstLen)
	nowIndex = 1
	for nowIndex < valueLen {
		// Read the data length
		length, lenOfLen := binary.Varint(value[nowIndex:])

		// Check the number of bytes read
		if lenOfLen <= 0 {
			return nil, ErrInvalidValue
		}

		// Jump to the start of the data
		nowIndex += lenOfLen

		// Check if the next operation will go out of bounds
		if nowIndex+int(length) > valueLen {
			return nil, ErrInvalidValue
		}

		// Add the data to the result
		result = append(result, make([]byte, length))
		copy(result[len(result)-1], value[nowIndex:nowIndex+int(length)])

		// Jump to the next data length
		nowIndex += int(length)
	}

	return result, nil
}
