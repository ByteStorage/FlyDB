package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initTestSetDb() (*SetStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestSetStructure")
	opts.DirPath = dir
	str, _ := NewSetStructure(opts)
	return str, &opts
}
func TestSAdd(t *testing.T) {
	s, _ := initTestSetDb()

	err := s.SAdd("destination", "key1")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}
func TestSAdds(t *testing.T) {
	s, _ := initTestSetDb()

	err := s.SAdds("destination", "key1", "key2")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}
func TestSRems(t *testing.T) {
	s, _ := initTestSetDb()

	err := s.SRems("destination", "key1", "key2")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}
func TestSRem(t *testing.T) {
	s, _ := initTestSetDb()

	err := s.SRem("destination", "key1")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}
func TestSCard(t *testing.T) {
	s, _ := initTestSetDb()

	got, err := s.SCard("destination")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	assert.NotNil(t, got)

}
func TestSMembers(t *testing.T) {
	s, _ := initTestSetDb()

	mem, err := s.SMembers("destination")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	assert.IsType(t, []string{}, mem)

}
func TestSIsMember(t *testing.T) {
	s, _ := initTestSetDb()

	got, err := s.SIsMember("destination", "key1")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	assert.False(t, got)

}

func TestNewSetStructure(t *testing.T) {
	// expect error
	opts := config.DefaultOptions
	opts.DirPath = "" // the cause of error
	setDB, err := NewSetStructure(opts)
	assert.NotNil(t, err)
	assert.Nil(t, setDB)
	// expect no error
	opts = config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestSetStructure")
	opts.DirPath = dir
	setDB, _ = NewSetStructure(opts)
	assert.NotNil(t, setDB)
	assert.IsType(t, &engine.DB{}, setDB.db)
}

func TestSUnion(t *testing.T) {
	s, _ := initTestSetDb()
	keys := []string{"key1", "key2"}

	_, err := s.SUnion(keys...)
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}

func TestSInter(t *testing.T) {
	s, _ := initTestSetDb()
	keys := []string{"key1", "key2"}

	_, err := s.SInter(keys...)
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}

func TestSDiff(t *testing.T) {
	s, _ := initTestSetDb()

	_, err := s.SDiff("key1", "key2")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}

func TestSUnionStore(t *testing.T) {
	s, _ := initTestSetDb()

	err := s.SUnionStore("destination", "key1", "key2")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}

func TestSInterStore(t *testing.T) {
	s, _ := initTestSetDb()

	err := s.SInterStore("destination", "key1", "key2")
	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}

}
