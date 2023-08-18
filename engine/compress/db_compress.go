package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/engine"
	"io/ioutil"
	"os"
)

// Compress uses gzip to compress data
func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)

	_, err := gw.Write(data)
	if err != nil {
		return nil, err
	}
	if err := gw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decompress uses gzip to decompress data
func decompress(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	gr, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	return ioutil.ReadAll(gr)
}

type DbCompress struct {
	db  *engine.DB
	dir string
}

func NewDbCompress(options config.Options) (*DbCompress, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &DbCompress{
		db:  db,
		dir: options.DirPath,
	}, nil
}

func (d *DbCompress) Put(key, value []byte) error {
	compressValue, err := compress(value)
	if err != nil {
		return err
	}
	return d.db.Put(key, compressValue)
}

func (d *DbCompress) Get(key []byte) ([]byte, error) {
	value, err := d.db.Get(key)
	if err != nil {
		return nil, err
	}
	return decompress(value)
}

func (d *DbCompress) Clean() {
	if d.db != nil {
		_ = d.db.Close()
		err := os.RemoveAll(d.dir)
		if err != nil {
			panic(err)
		}
	}
}
