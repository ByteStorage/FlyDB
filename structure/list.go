package structure

import (
	"bytes"
	"errors"
	"reflect"

	"encoding/gob"

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

type listNode struct {
	Value interface{}
	Next  *listNode
}

type list struct {
	Head   *listNode
	Length int
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
func (l *ListStructure) LPush(key string, value interface{}) error {
	// Get the list
	lst, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	newNode := &listNode{
		Value: value,
		Next:  lst.Head,
	}

	lst.Head = newNode
	lst.Length++

	// Store to db
	return l.setListToDB(key, lst)
}

// LPushs adds one or more values to the left of the list corresponding to the key
// If the key does not exist, it will create the key
func (l *ListStructure) LPushs(key string, values ...interface{}) error {
	// Check if values are valid
	if len(values) == 0 {
		return ErrInvalidArgs
	}

	// Get the list
	lst, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	for i := len(values) - 1; i >= 0; i-- {
		newNode := &listNode{
			Value: values[i],
			Next:  lst.Head,
		}
		lst.Head = newNode
		lst.Length++
	}

	// Store to db
	return l.setListToDB(key, lst)
}

// RPush adds a value to the right of the list corresponding to the key
// If the key does not exist, it will create the key
func (l *ListStructure) RPush(key string, value interface{}) error {
	// Check if value is empty
	if value == nil {
		return ErrInvalidValue
	}

	// Get the list
	lst, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	// Append the new data to the end
	newNode := &listNode{
		Value: value,
		Next:  nil,
	}
	if lst.Length == 0 {
		lst.Head = newNode
	} else {
		// Find the last node
		lastNode := lst.Head
		for lastNode.Next != nil {
			lastNode = lastNode.Next
		}
		lastNode.Next = newNode
	}
	lst.Length++

	// Store to db
	return l.setListToDB(key, lst)
}

// RPushs appends one or more values to the right side of a list associated with a key.
// If the key does not exist, it will be created.
func (l *ListStructure) RPushs(key string, values ...interface{}) error {
	// Check if values are valid
	if len(values) == 0 {
		return ErrInvalidArgs
	}

	// Get the list
	lst, err := l.getListFromDB(key, true)
	if err != nil {
		return err
	}

	// Find the last node
	var lastNode *listNode
	if lst.Length == 0 {
		lastNode = nil
	} else {
		lastNode = lst.Head
		for lastNode.Next != nil {
			lastNode = lastNode.Next
		}
	}

	for _, value := range values {
		newNode := &listNode{
			Value: value,
			Next:  nil,
		}
		if lastNode == nil {
			lst.Head = newNode
		} else {
			lastNode.Next = newNode
		}
		lastNode = newNode
		lst.Length++
	}

	// Store to db
	return l.setListToDB(key, lst)
}

// LPop returns and removes the leftmost value of a list associated with a key.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
func (l *ListStructure) LPop(key string) (interface{}, error) {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if lst.Length == 0 {
		return nil, ErrListEmpty
	}

	popValue := lst.Head.Value
	lst.Head = lst.Head.Next
	lst.Length--

	// Store in the database
	return popValue, l.setListToDB(key, lst)
}

// RPop returns and removes the rightmost value of a list associated with a key.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
func (l *ListStructure) RPop(key string) (interface{}, error) {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if lst.Length == 0 {
		return nil, ErrListEmpty
	} else if lst.Length == 1 {
		popValue := lst.Head.Value
		lst.Head = nil
		lst.Length = 0
		return popValue, l.setListToDB(key, lst)
	}

	// Find the new tail
	newTail := lst.Head
	for i := 0; i < lst.Length-1; i++ {
		newTail = newTail.Next
	}
	popValue := newTail.Next.Value
	newTail.Next = nil
	lst.Length--

	// Store in the database
	return popValue, l.setListToDB(key, lst)
}

// LRange returns a range of elements from a list associated with a key.
// The range is inclusive, including both the start and stop indices.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
// Negative indices can be used, where -1 represents the last element of the list,
// -2 represents the second last element, and so on.
func (l *ListStructure) LRange(key string, start int, stop int) ([]interface{}, error) {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if lst.Length == 0 {
		return nil, ErrListEmpty
	}

	// Calculate the correct indices
	start = (start%lst.Length + lst.Length) % lst.Length
	stop = (stop%lst.Length + lst.Length) % lst.Length

	// Return empty if the range length is less than 1
	if stop < start {
		return nil, nil
	}

	nowNode := lst.Head

	for i := 0; i < start; i++ {
		nowNode = nowNode.Next
	}
	result := make([]interface{}, 0, stop-start+1)
	for i := start; i <= stop; i++ {
		result = append(result, nowNode.Value)
		nowNode = nowNode.Next
	}

	return result, nil
}

// LLen returns the size of a list associated with a key.
// If the key does not exist, an error is returned.
func (l *ListStructure) LLen(key string) (int, error) {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return 0, err
	}

	return lst.Length, nil
}

// LRem removes elements from a list associated with a key based on the count and value parameters.
// The count can have the following values:
// count > 0: Remove count occurrences of the value from the beginning of the list.
// count < 0: Remove count occurrences of the value from the end of the list.
// count = 0: Remove all occurrences of the value from the list.
// If the key does not exist, an error is returned.
func (l *ListStructure) LRem(key string, count int, value interface{}) error {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return err
	}

	// Process different counts
	if count != 0 {
		prev, curr := lst.Head, lst.Head
		for i := 0; (count > 0 && i < count || count < 0 && lst.Length+i > -count) && curr != nil; {
			if l.valueEqual(curr.Value, value) {
				if curr == lst.Head {
					lst.Head = curr.Next
				} else {
					prev.Next = curr.Next
				}
				i++
				lst.Length--
			} else {
				prev = curr
			}
			curr = curr.Next
		}
	} else {
		prev, curr := lst.Head, lst.Head
		for curr != nil {
			if l.valueEqual(curr.Value, value) {
				if curr == lst.Head {
					lst.Head = curr.Next
				} else {
					prev.Next = curr.Next
				}
				lst.Length--
			} else {
				prev = curr
			}
			curr = curr.Next
		}
	}

	// Store to db
	return l.setListToDB(key, lst)
}

// LSet sets the value of an element in a list associated with a key based on the index.
// If the index is out of range, an error is returned.
// If the list is empty, an error is returned.
func (l *ListStructure) LSet(key string, index int, value interface{}) error {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return err
	}

	// Check if the index is out of range
	if index < 0 || index >= lst.Length {
		return ErrIndexOutOfRange
	}

	nowNode := lst.Head

	for i := 0; i < index; i++ {
		nowNode = nowNode.Next
	}

	nowNode.Value = value

	// Store in the database
	return l.setListToDB(key, lst)
}

// LTrim retains a range of elements in a list associated with a key.
// The range is inclusive, including both the start and stop indices.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
// Negative indices can be used, where -1 represents the last element of the list,
// -2 represents the second last element, and so on.
func (l *ListStructure) LTrim(key string, start int, stop int) error {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return err
	}

	if lst.Length == 0 {
		return ErrListEmpty
	}

	// Calculate the correct indices
	start = (start%lst.Length + lst.Length) % lst.Length
	stop = (stop%lst.Length + lst.Length) % lst.Length

	if start > stop {
		lst = &list{
			Head:   nil,
			Length: 0,
		}
	} else {
		// Find the new head
		newHead := lst.Head
		for i := 0; i < start; i++ {
			newHead = newHead.Next
		}

		// Find the new tail
		newTail := newHead
		for i := start; i < stop; i++ {
			newTail = newTail.Next
		}

		// Disconnect the new tail from the rest of the list
		newTail.Next = nil

		// Update the list
		lst.Head = newHead
		lst.Length = stop - start + 1
	}

	// Store in the database
	return l.setListToDB(key, lst)
}

// LIndex returns the value of an element in a list associated with a key based on the index.
// If the key does not exist, an error is returned.
// If the list is empty, an error is returned.
// Negative indices can be used, where -1 represents the last element of the list,
// -2 represents the second last element, and so on.
func (l *ListStructure) LIndex(key string, index int) (interface{}, error) {
	// Get the list
	lst, err := l.getListFromDB(key, false)
	if err != nil {
		return nil, err
	}

	// Return error if the list is empty
	if lst.Length == 0 {
		return nil, ErrListEmpty
	}

	// Calculate the correct index
	index = (index%lst.Length + lst.Length) % lst.Length

	nowNode := lst.Head

	for i := 0; i < index; i++ {
		nowNode = nowNode.Next
	}

	return nowNode.Value, nil
}

// RPOPLPUSH removes the last element from one list and pushes it to another list.
// If the source list is empty, an error is returned.
// If the destination list is empty, it is created.
// Atomicity is not guaranteed.
func (l *ListStructure) RPOPLPUSH(source string, destination string) error {
	// Get the source list
	lst1, err := l.getListFromDB(source, false)
	if err != nil {
		return err
	}

	// Get the destination list
	lst2, err := l.getListFromDB(destination, true)
	if err != nil {
		return err
	}

	// Return error if the source list is empty
	if lst1.Length == 0 {
		return ErrListEmpty
	}

	// Find the last node of the source list
	lastNode := lst1.Head
	prevNode := lst1.Head
	for lastNode.Next != nil {
		prevNode = lastNode
		lastNode = lastNode.Next
	}

	// Remove the last node from the source list
	if lst1.Length == 1 {
		lst1.Head = nil
	} else {
		prevNode.Next = nil
	}
	lst1.Length--

	// Add the last node to the head of the destination list
	lastNode.Next = lst2.Head
	lst2.Head = lastNode
	lst2.Length++

	// Store in the database
	err = l.setListToDB(source, lst1)
	if err != nil {
		return err
	}
	return l.setListToDB(destination, lst2)
}

var (
	// ErrListEmpty is returned if the list is empty.
	ErrListEmpty = errors.New("Wrong operation: list is empty")
	// ErrIndexOutOfRange is returned if the index out of range.
	ErrIndexOutOfRange = errors.New("Wrong operation: index out of range")
)

// valueEqual checks if two values are equal.
// It supports comparing string, int, float64, bool, and []byte values.
func (l *ListStructure) valueEqual(value1 interface{}, value2 interface{}) bool {
	type1 := reflect.TypeOf(value1)
	type2 := reflect.TypeOf(value2)

	// If the types are not the same, the values are not equal
	if type1 != type2 {
		return false
	}

	// Compare based on type
	switch type1.Kind() {
	case reflect.String:
		return value1.(string) == value2.(string)
	case reflect.Int:
		return value1.(int) == value2.(int)
	case reflect.Float64:
		return value1.(float64) == value2.(float64)
	case reflect.Bool:
		return value1.(bool) == value2.(bool)
	case reflect.Slice:
		// Special case for []byte
		if type1.Elem().Kind() == reflect.Uint8 {
			return bytes.Equal(value1.([]byte), value2.([]byte))
		}
	}

	// For other types, return false
	return false
}

// getListFromDB retrieves data from the database. When isKeyCanNotExist is true, it returns an empty slice if the key doesn't exist instead of an error.
func (l *ListStructure) getListFromDB(key string, isKeyCanNotExist bool) (*list, error) {
	if isKeyCanNotExist {
		// Get data corresponding to the key from the database
		dbData, err := l.db.Get([]byte(key))

		// Since the key might not exist, we need to handle ErrKeyNotFound separately as it is a valid case
		if err != nil && err != _const.ErrKeyNotFound {
			return nil, err
		}

		// Deserialize the data into a list
		lst, err := l.decodeList(dbData)
		if err != nil {
			if len(dbData) != 0 {
				return nil, err
			} else {
				lst = &list{nil, 0}
			}
		}
		return lst, nil
	} else {
		// Get data corresponding to the key from the database
		dbData, err := l.db.Get([]byte(key))
		if err != nil {
			return nil, err
		}

		// Deserialize the data into a list
		lst, err := l.decodeList(dbData)
		if err != nil {
			return nil, err
		}
		return lst, nil
	}
}

// setListToDB stores the data into the database.
func (l *ListStructure) setListToDB(key string, lst *list) error {
	// Serialize into binary array
	encValue, err := l.encodeList(lst)
	if err != nil {
		return err
	}
	// Store in the database
	return l.db.Put([]byte(key), encValue)
}

// encodeList encodes the value
// format: [type][data]
func (l *ListStructure) encodeList(lst *list) ([]byte, error) {
	gob.Register(&list{})
	var buffer bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&buffer) // Will write to network.

	// Encode (send) the value.
	err := enc.Encode(lst)
	if err != nil {
		return nil, err
	}

	encodedBytes := buffer.Bytes()

	// Add a byte to the head of the slice.
	encodedBytes = append([]byte{List}, encodedBytes...)

	return encodedBytes, nil
}

// decodeList decodes the value
func (l *ListStructure) decodeList(value []byte) (*list, error) {
	// Check the length of the value
	if len(value) < 1 {
		return nil, ErrInvalidValue
	}

	// Check the type of the value
	if value[0] != List {
		return nil, ErrInvalidType
	}

	// Create a bytes.Buffer from the value (excluding the first byte)
	buffer := bytes.NewBuffer(value[1:])

	// Create a new gob.Decoder
	dec := gob.NewDecoder(buffer)

	// Create a new list to hold the decoded value
	var lst list

	// Decode the value
	err := dec.Decode(&lst)
	if err != nil {
		return nil, err
	}

	return &lst, nil
}

func (s *ListStructure) Stop() error {
	err := s.db.Close()
	return err
}
