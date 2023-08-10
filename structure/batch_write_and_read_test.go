package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initStringCur() (*StringStructure, *config.Options) {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	opts := config.DefaultOptions
	opts.DirPath = dir + "/flydb"
	str, _ := NewStringStructure(opts)
	return str, &opts
}

func initHashCur() (*HashStructure, *config.Options) {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	opts := config.DefaultOptions
	opts.DirPath = dir + "/flydb"
	hash, _ := NewHashStructure(opts)
	return hash, &opts
}

func initListCur() (*ListStructure, *config.Options) {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	opts := config.DefaultOptions
	opts.DirPath = dir + "/flydb"
	list, _ := NewListStructure(opts)
	return list, &opts
}

func initSetCur() (*SetStructure, *config.Options) {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	opts := config.DefaultOptions
	opts.DirPath = dir + "/flydb"
	str, _ := NewSetStructure(opts)
	return str, &opts
}

func TestBatchWrite(t *testing.T) {
	dbs, _ := initStringCur()
	defer dbs.Clean()
	dbh, _ := initHashCur()
	defer dbh.Clean()
	dbl, _ := initListCur()
	defer dbl.Stop()
	dbset, _ := initSetCur()
	defer dbset.Stop()
	err = dbs.Set("1", "1", 0)
	assert.Nil(t, err)
	_, err = dbh.HSet("2", "1", "1", 0)
	assert.Nil(t, err)
	err = dbl.LPush("3", "1")
	assert.Nil(t, err)
	err = dbset.SAdd("4", "1")
	assert.Nil(t, err)
}

func TestBatchWriteAndRead(t *testing.T) {
	dbs, _ := initStringCur()
	defer dbs.Clean()
	dbh, _ := initHashCur()
	defer dbh.Clean()
	dbl, _ := initListCur()
	defer dbl.Stop()
	dbset, _ := initSetCur()
	defer dbset.Stop()

	err := dbs.Set("1", "1", 0)
	assert.Nil(t, err)
	val, err := dbs.Get("1")
	assert.Nil(t, err)
	assert.Equal(t, "1", val)

	err = dbl.LPush("3", "1")
	assert.Nil(t, err)
	val, err = dbl.LPop("3")
	assert.Nil(t, err)
	assert.Equal(t, "1", val)

	err = dbset.SAdd("4", "1")
	assert.Nil(t, err)
	assert.Equal(t, "1", val)
	members, err := dbset.SMembers("4")
	assert.Nil(t, err)
	assert.Equal(t, "1", members[0])

	_, err = dbh.HSet("2", "1", "1", 0)
	assert.Nil(t, err)
	val, err = dbh.HGet("2", "1")
	assert.Nil(t, err)
	assert.Equal(t, "1", val)
}
