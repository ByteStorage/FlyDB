package structure

import (
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"regexp"
	"time"
)

type SetStructure struct {
	db *engine.DB
}

var (
	ErrMemberNotFound    = errors.New("ErrMemberNotFound: member not found")
	ErrSetNotInitialized = errors.New("wrong operation: set not initialized")
)

// FSets serves as a set data structure where every key in the map is an element of the set.
type FSets map[string]struct{}

func NewSetStructure(options config.Options) (*SetStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &SetStructure{db: db}, nil
}

// SAdd adds a member to the set stored at key.
//
// If the set did not exist, a new set will be created
// and the member will be added to it.
func (s *SetStructure) SAdd(key, member string, ttl int64) error {
	return s.SAdds(key, ttl, member)
}

/*
SAdds attempts to add multiple members to a set identified by the 'key' parameter. The members to add are variadic and come as 'members...' in the function signature. This means the function accepts an arbitrary number of strings to be added to the set.

Parameters:
  - 's': Pointer to an instance of SetStructure. This is the receiver of the SAdds function.
  - 'key': String that is used as the identifier for the set.
  - 'members': Variadic parameter that represents the members to be added. Since it's a variadic parameter, it can be any number of strings.

Return:
  - It returns an error if it encounters one during execution. Errors might occur if the key is found to be empty or if problems arise when retrieving the set from the database or saving it back.

Internal logic:
 1. The function first checks if the provided key is an empty string. If so, it returns an error indicating that the key is empty.
 2. Converts the key into bytes format for internal use.
 3. It tries to get the set associated with the key from the database. If it encounters an error during retrieval, it returns that error.
 4. If it can successfully retrieve the set, it attempts to add the members to the set.
 5. After adding the members, it tries to save this updated set back into the database. If this operation throws an error, it returns that error.

All the methods like 'getZSetFromDB', 'setZSetToDB', and 'add' handle the lower-level logic associated with database interaction and set manipulation.
*/
func (s *SetStructure) SAdds(key string, ttl int64, members ...string) error {
	fs, err := s.checkAndGetSet(key, true)
	if err != nil {
		return err
	}
	var expirationTime time.Duration
	fs.add(members...)
	expirationTime = time.Duration(ttl) * time.Second
	return s.setSetToDB(stringToBytesWithKey(key), fs, expirationTime)
}

// SRem removes a member from a set
func (s *SetStructure) SRem(key, member string) error {
	return s.SRems(key, 0, member)
}

// SRems removes multiple members from a set
func (s *SetStructure) SRems(key string, ttl int64, members ...string) error {
	fs, err := s.checkAndGetSet(key, false)
	if err != nil {
		return err
	}
	var expirationTime time.Duration
	if err = fs.remove(members...); err != nil {
		return err
	}
	expirationTime = time.Duration(ttl) * time.Second
	return s.setSetToDB(stringToBytesWithKey(key), fs, expirationTime)
}

// SCard gets the cardinality (size) of a set
func (s *SetStructure) SCard(key string) (int, error) {
	fs, err := s.checkAndGetSet(key, false)
	if err != nil {
		return -1, err
	}
	return len(*fs), nil
}

// SMembers gets all members of a set identified by the provided key.
// If the set does not exist or the key belongs to a different data type,
// the function will return an error.
// It returns a slice of strings which are the members of the set.
func (s *SetStructure) SMembers(key string) ([]string, error) {
	fs, err := s.checkAndGetSet(key, false)
	if err != nil {
		return nil, err
	}
	// type of fs is map[string]struct{}
	// we need to take the keys and make them a slice
	var members []string
	for k := range *fs {
		members = append(members, k)
	}

	return members, nil
}

// SIsMember checks if a member exists in a set
func (s *SetStructure) SIsMember(key, member string) (bool, error) {
	fs, err := s.checkAndGetSet(key, false)
	if err != nil {
		return false, err
	}

	return fs.exists(member), nil
}

// SUnion gets the union of multiple sets
func (s *SetStructure) SUnion(keys ...string) ([]string, error) {
	// if there are no keys provided, then there's no union
	if len(keys) == 0 {
		return nil, nil
	}
	mem := make(map[string]struct{})
	var members []string
	for _, key := range keys {
		fs, err := s.checkAndGetSet(key, false)
		if err != nil {
			return nil, err
		}
		for k := range *fs {
			if _, ok := mem[k]; !ok {
				mem[k] = struct{}{}
				members = append(members, k)
			}
		}
	}
	return members, nil
}

// SInter returns the intersection of multiple sets.
// The parameter 'keys' is a variadic parameter, meaning
// it can accept any number of arguments. These arguments are the keys
// of the sets that are to be intersected.
func (s *SetStructure) SInter(keys ...string) ([]string, error) {
	// if there are no keys, then there's no intersection
	if len(keys) == 0 {
		return nil, nil
	}
	// All elements of 'first' are stored in 'mem' map because it's checked
	// against each of the subsequent sets.
	first, err := s.checkAndGetSet(keys[0], false)
	if err != nil {
		return nil, err
	}

	mem := make(map[string]struct{})
	for k := range *first {
		mem[k] = struct{}{}
	}

	var members []string
	// For each other key, we get its set and its members.
	for _, key := range keys[1:] {
		fs, err := s.checkAndGetSet(key, false)
		if err != nil {
			return nil, err
		}
		// We check if each member of the current set is in 'mem'.
		// If yes, we add it to 'members' array, and remove it from 'mem' map.
		for k := range *fs {
			if _, ok := mem[k]; ok {
				delete(mem, k)
				members = append(members, k)
			}
		}
	}

	return members, nil
}

// SDiff computes the difference between the first and subsequent sets.
// It returns keys unique to the first set and an error if applicable.
func (s *SetStructure) SDiff(keys ...string) ([]string, error) {
	// If there are no keys, then there's no difference
	if len(keys) == 0 {
		return nil, nil
	}

	first, err := s.checkAndGetSet(keys[0], false)
	if err != nil {
		return nil, err
	}

	// Initialize a set for members of first set
	mem := make(map[string]struct{})
	for k := range *first {
		mem[k] = struct{}{}
	}

	var diffMembers []string
	// If member is in the first set and also in another set, remove it.
	for _, key := range keys[1:] {
		fs, err := s.checkAndGetSet(key, false)
		if err != nil {
			return nil, err
		}
		for k := range *fs {
			// Only delete the existing members in the first set. Do not add members from other sets
			delete(mem, k)
		}
	}
	// Remaining members in mem are unique to first set.
	for k := range mem {
		diffMembers = append(diffMembers, k)
	}

	return diffMembers, nil
}

// Keys returns all the keys of the set structure
func (s *SetStructure) Keys(regx string) ([]string, error) {
	toRegexp := convertToRegexp(regx)
	compile, err := regexp.Compile(toRegexp)
	if err != nil {
		return nil, err
	}
	var keys []string
	byteKeys := s.db.GetListKeys()
	for _, key := range byteKeys {
		if compile.MatchString(string(key)) {
			// check if deleted
			if !s.exists(string(key)) {
				continue
			}
			keys = append(keys, string(key))
		}
	}
	return keys, nil
}

// SUnionStore calculates and stores the union of multiple sets
// in a destination set.
//
// 'destination' is the name of the set where the result will be stored.
//
// 'keys' is a variadic parameter that represents the names of all the sets
// whose union is to be calculated.
//
// The function returns an error if there's any issue computing the union or
// storing the result in the destination set.
//
// Usage example:
//
//	set := NewSetStructure(opts) // initialize t
//	err := set.SUnionStore("result", "set1", "set2", "set3")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (s *SetStructure) SUnionStore(destination string, keys ...string) error {
	union, err := s.SUnion(keys...)
	if err != nil {
		return err
	}
	return s.SAdds(destination, 0, union...)
}

/*
SInterStore stores the intersection of multiple sets in a destination set.

Function parameters:
  - destination (string): The key for the set where the intersection result will be stored.
  - keys (string[]): An array of keys for the sets to be intersected.

This function first uses the SInter function to find the intersection of the provided sets (keys).
In case an error occurs during this process (the SInter function returns an error), it will immediately return that error.
Otherwise, it proceeds to use the SAdds function to store the results of the intersection into the destination set.

The function will return the error from the SAdds function if it occurs or nil if the operation is successful.

Errors can occur when there's an issue with accessing the data structures involved in the operations. For instance, trying to perform set operations on non-existing keys or keys that point to non-set data types.

Returns:
  - error: An error object that describes an error that occurred during the computation (if any).
*/
func (s *SetStructure) SInterStore(destination string, keys ...string) error {
	inter, err := s.SInter(keys...)
	if err != nil {
		return err
	}
	return s.SAdds(destination, 0, inter...)
}

func (s *SetStructure) checkAndGetSet(key string, createIfNotExist bool) (*FSets, error) {
	// Check if value is empty
	if err := checkKey(key); err != nil {
		return nil, err
	}
	keyBytes := stringToBytesWithKey(key)
	// Get the list
	set, _, err := s.getSetFromDB(keyBytes, createIfNotExist)
	if err != nil {
		return nil, err
	}

	return set, nil
}

// internal/private functions

// add adds new elements to the set data structure (FSets).
// It allows one or more string parameters. If an element already exists, it won't be added again.
func (s *FSets) add(member ...string) {

	for _, m := range member {
		// Check if the member to be added already exists in the FSets or not.
		// existence of the element in the FSets needs to be checked, not the value.
		if _, exists := (*s)[m]; !exists {
			// If the elements do not already exist, add it to the set.
			// The value here is an empty struct, because we are only interested in the keys.
			(*s)[m] = struct{}{}
		}
	}
}

// remove is a method on the FSet struct that takes in a variable number
// of string parameters. It takes these strings to be "members", and generates errors
// if any of these strings are not currently present in the FSet.
//
// Parameters:
//
//	member: a variable list of strings that are supposed to be "members" of the FSet.
//
// Returns:
//
//	If all of the strings in the members parameter are present in the FSet,
//	it then proceeds to remove each of these members from the FSet,
//	and returns nil, signifying success.
//
//	If any string from the members parameter does not exist in the FSet,
//	the function immediately returns an ErrInvalidArgs error,
//	which signifies that an invalid argument has been provided to the function.
func (s *FSets) remove(member ...string) error {
	if s == nil {
		return ErrSetNotInitialized
	}
	for _, m := range member {
		// check to see if all members exist
		if _, exists := (*s)[m]; !exists {
			return ErrMemberNotFound
		}
	}
	// iterate again
	for _, m := range member {
		// remove elements from FSets
		delete(*s, m)
	}
	return nil
}

// GetSetFromDB retrieves a set from database given a key. If createIfNotExist is true,
// a new set will be created if the key is not found. It returns the file sets and any write error encountered.
func (s *SetStructure) getSetFromDB(key []byte, createIfNotExist bool) (*FSets, int64, error) {
	if s.db == nil {
		return nil, 0, ErrSetNotInitialized
	}
	var zSetValueWithTTL FSetWithTTL
	dbData, err := s.db.Get(key)

	// If key is not found, return nil for both; otherwise return the error.
	if err != nil {
		if errors.Is(err, _const.ErrKeyNotFound) && createIfNotExist {
			return &FSets{}, 0, nil
		}
		return nil, 0, err
	} else {
		err = encoding.NewMessagePackDecoder(dbData).Decode(&zSetValueWithTTL) // Decode the value along with TTL
		if err != nil {
			return nil, 0, err
		}

		expiration := zSetValueWithTTL.TTL
		if expiration != 0 && expiration < time.Now().UnixNano() {
			return nil, -1, _const.ErrKeyIsExpired
		}

		return zSetValueWithTTL.ZSet, zSetValueWithTTL.TTL, nil // Return the zSetValue and TTL
	}
}

type FSetWithTTL struct {
	ZSet *FSets
	TTL  int64
}

// setSetToDB
func (s *SetStructure) setSetToDB(key []byte, zSetValue *FSets, ttl time.Duration) error {
	// create a new zSetValueWithTTL struct
	var expire int64 = 0

	if ttl != 0 {
		expire = time.Now().Add(ttl).UnixNano()
	}

	valueWithTTL := &FSetWithTTL{
		ZSet: zSetValue,
		TTL:  expire,
	}

	val := encoding.NewMessagePackEncoder()
	err := val.Encode(valueWithTTL) // Encode the value along with TTL
	if err != nil {
		return err
	}
	return s.db.Put(key, val.Bytes())
}

func (s *SetStructure) exists(key string, member ...string) bool {
	if err := checkKey(key); err != nil {
		return false
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, _, err := s.getSetFromDB(keyBytes, false)

	if err != nil {
		return false
	}
	return zSet.exists(member...)
}
func (s *FSets) exists(member ...string) bool {
	for _, s2 := range member {
		if _, ok := (*s)[s2]; !ok {
			return false
		}
	}
	return true
}

func (s *SetStructure) Stop() error {
	err := s.db.Close()
	return err
}

func (s *SetStructure) TTL(k string) (int64, error) {
	keyBytes := stringToBytesWithKey(k)
	_, expire, err := s.getSetFromDB(keyBytes, false)
	if err != nil {
		return -1, err
	}

	now := time.Now().UnixNano() / int64(time.Second)
	expire = expire / int64(time.Second)

	remainingTTL := expire - now

	//println("re",remainingTTL)
	if remainingTTL <= 0 {
		return 0, nil // Return 0 TTL for expired keys
	}

	return remainingTTL, nil
}

func (s *SetStructure) Size(key string) (string, error) {

	members, err := s.SMembers(key)
	if err != nil {
		return "", err
	}
	var sizeInBytes int

	// Calculate the size of the value
	for _, v := range members {
		toString, err := interfaceToString(v)
		if err != nil {
			return "", err
		}
		sizeInBytes += len(toString)
	}

	// Convert bytes to corresponding units (KB, MB...)
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)

	var size string
	switch {
	case sizeInBytes < KB:
		size = fmt.Sprintf("%dB", sizeInBytes)
	case sizeInBytes < MB:
		size = fmt.Sprintf("%.2fKB", float64(sizeInBytes)/KB)
	case sizeInBytes < GB:
		size = fmt.Sprintf("%.2fMB", float64(sizeInBytes)/MB)
	}

	return size, nil
}

func (s *SetStructure) SDel(key string) error {
	byteKey := stringToBytesWithKey(key)
	err := s.db.Delete(byteKey)
	if err != nil {
		return err
	}
	return nil
}
