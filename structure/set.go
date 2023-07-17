package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
)

type SetStructure struct {
	db *engine.DB
}

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
func (s *SetStructure) SAdd(key, member string) error {

	return nil
}

// SAdds adds multiple members to a set
func (s *SetStructure) SAdds(key string, members ...string) error {
	return nil
}

// SRem removes a member from a set
func (s *SetStructure) SRem(key, member string) error {
	return nil
}

// SRems removes multiple members from a set
func (s *SetStructure) SRems(key string, members ...string) error {
	return nil
}

// SCard gets the cardinality (size) of a set
func (s *SetStructure) SCard(key string) (int, error) {
	return 0, nil
}

// SMembers gets all members of a set
func (s *SetStructure) SMembers(key string) ([]string, error) {
	return nil, nil
}

// SIsMember checks if a member exists in a set
func (s *SetStructure) SIsMember(key, member string) (bool, error) {
	return false, nil
}

// SUnion gets the union of multiple sets
func (s *SetStructure) SUnion(keys ...string) ([]string, error) {
	return nil, nil
}

// SInter gets the intersection of multiple sets
func (s *SetStructure) SInter(keys ...string) ([]string, error) {
	return nil, nil
}

// SDiff gets the difference between two sets
func (s *SetStructure) SDiff(key1, key2 string) ([]string, error) {
	return nil, nil
}

// SUnionStore stores the union of multiple sets in a destination set
func (s *SetStructure) SUnionStore(destination, key1, key2 string) error {
	return nil
}

// SInterStore stores the intersection of multiple sets in a destination set
func (s *SetStructure) SInterStore(destination, key1, key2 string) error {
	return nil
}
