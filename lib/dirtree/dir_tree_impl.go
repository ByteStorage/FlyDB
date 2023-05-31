package dirtree

var _ DirTreeInterface = &DirTree{}

/**TODO:
	这里初步计划只维护一个目录树
			/
   /节点1  /节点2  /节点3
   /file1  /file2  /file3
	大概结构如上，每个节点都是一个目录，
	每个目录下面都有文件，这样的话只需要维护一个目录树就可以了，
	而不需要维护节点与文件的映射关系，这里DirTree只是写了个简单的目录树，可以找找有没有优化的目录树
*/

// DirTree is a struct that represents a directory tree
type DirTree struct {
	// Root is the root node of the directory tree
	Root *DirTreeNode
}

// DirTreeNode is a struct that represents a directory tree node
type DirTreeNode struct {
	Name     string
	Children []*DirTreeNode
}

func (d *DirTree) InsertFile(path string) bool {
	//TODO implement me
	panic("implement me")
}

func (d *DirTree) DeleteFile(filename string) bool {
	//TODO implement me
	panic("implement me")
}

func (d *DirTree) DeleteDir(path string) bool {
	//TODO implement me
	panic("implement me")
}

func (d *DirTree) MkDir(path string) bool {
	//TODO implement me
	panic("implement me")
}
