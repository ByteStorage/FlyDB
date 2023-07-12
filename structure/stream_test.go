package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	"os"
	"testing"
)

func initStreamDB() (*StreamStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestStreamStructure_Get")
	opts.DirPath = dir
	stc, _ := NewStreamStructure(opts)
	return stc, &opts
}

func TestStreamStructure_XRead(t *testing.T) {
	stc, _ := initStreamDB()
	defer stc.db.Clean()

	ok1, err := stc.XAdd("test", "1", map[string]interface{}{"name1": "flydb1"})
	t.Log(ok1, err)

	ok2, err := stc.XAdd("test", "2", map[string]interface{}{"name2": "flydb2"})
	t.Log(ok2, err)

	ok3, err := stc.XAdd("test", "3", map[string]interface{}{"name3": "flydb3"})
	t.Log(ok3, err)
	//
	ok4, err := stc.XAdd("test1", "1", map[string]interface{}{"name11": "flydb11"})
	t.Log(ok4, err)

	item, err := stc.XRead("test", 1)
	t.Log(item, err)

	item1, err := stc.XRead("test1", 1)
	t.Log(item1, err)

}
