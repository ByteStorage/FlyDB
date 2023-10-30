package sql

import v1 "github.com/ByteStorage/FlyDB/engine/sql/engine"

type Engine interface {
	Execute(sql string, args ...interface{}) (int64, error)
}

type engine struct {
}

func NewEngine() Engine {
	return &engine{}
}

func (e *engine) Execute(sql string, args ...interface{}) (int64, error) {
	// sql parse
	parser := v1.NewSQLParser(sql)
	_, err := parser.Parse(sql)
	if err != nil {
		return 0, err
	}
	// sql optimize

	// sql execute plan

	// sql execute

	return 0, nil
}
