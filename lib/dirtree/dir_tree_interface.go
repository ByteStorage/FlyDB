package dirtree

type DirTreeInterface interface {
	// InsertFile inserts a file into the directory tree
	InsertFile(path string) bool
	// DeleteFile deletes a file from the directory tree
	DeleteFile(filename string) bool
	// DeleteDir deletes a directory from the directory tree
	DeleteDir(path string) bool
	// MkDir creates a directory in the directory tree
	MkDir(path string) bool
}
