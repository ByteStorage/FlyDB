package main

import (
	"fmt"
	"sync"
)

type KeyValueDB struct {
	storage map[string]interface{}
	mu      sync.RWMutex
}

func NewKeyValueDB() *KeyValueDB {
	return &KeyValueDB{storage: make(map[string]interface{})}
}

func (db *KeyValueDB) Set(key string, value interface{}) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.storage[key] = value
}

func (db *KeyValueDB) Get(key string) interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.storage[key]
}

type SimpleSQL struct {
	db     *KeyValueDB
	tables map[string][]string
}

func NewSimpleSQL() *SimpleSQL {
	return &SimpleSQL{db: NewKeyValueDB(), tables: make(map[string][]string)}
}

func (s *SimpleSQL) CreateTable(tableName string, columns []string) {
	s.tables[tableName] = columns
	s.db.Set(tableName, []map[string]interface{}{})
}

func (s *SimpleSQL) Insert(tableName string, values []interface{}) {
	if _, ok := s.tables[tableName]; !ok {
		fmt.Println("Table not found!")
		return
	}

	columns := s.tables[tableName]
	if len(columns) != len(values) {
		fmt.Println("Column count mismatch!")
		return
	}

	row := make(map[string]interface{})
	for i, col := range columns {
		row[col] = values[i]
	}

	tableData := s.db.Get(tableName).([]map[string]interface{})
	tableData = append(tableData, row)
	s.db.Set(tableName, tableData)
}

func (s *SimpleSQL) Select(tableName string, where map[string]interface{}) []map[string]interface{} {
	if _, ok := s.tables[tableName]; !ok {
		fmt.Println("Table not found!")
		return nil
	}

	tableData := s.db.Get(tableName).([]map[string]interface{})
	var results []map[string]interface{}
	for _, row := range tableData {
		match := true
		for col, val := range where {
			if row[col] != val {
				match = false
				break
			}
		}
		if match {
			results = append(results, row)
		}
	}
	return results
}

func main() {
	sqlEngine := NewSimpleSQL()
	sqlEngine.CreateTable("users", []string{"id", "name", "age"})
	sqlEngine.Insert("users", []interface{}{1, "Alice", 25})
	sqlEngine.Insert("users", []interface{}{2, "Bob", 30})
	fmt.Println(sqlEngine.Select("users", nil))                                     // 查询所有用户
	fmt.Println(sqlEngine.Select("users", map[string]interface{}{"name": "Alice"})) // 查询名为Alice的用户
}
