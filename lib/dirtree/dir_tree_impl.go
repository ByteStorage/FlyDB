package dirtree

import (
	"fmt"
	"strings"
)

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

// NewDirTree Creates a new instance of DirTree
func NewDirTree() *DirTree {
	return &DirTree{
		Root: &DirTreeNode{
			Name:     "/",
			Children: nil,
		},
	}
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

// MkDir Inserts a new DirTreeNode into the DirTree
func (d *DirTree) MkDir(path string) bool {

	components := splitPath(path)

	curNode := d.Root

	for _, component := range components {
		child := findChildByName(curNode.Children, component)
		// if the current component does not exist, create a new one
		if child == nil {
			newNode := &DirTreeNode{
				Name:     component,
				Children: nil,
			}
			curNode.Children = append(curNode.Children, newNode)
			curNode = newNode
		} else {
			// otherwise go to the next level
			curNode = child
		}
	}

	return true
}

func findChildByName(nodes []*DirTreeNode, name string) *DirTreeNode {
	for _, node := range nodes {
		if node.Name == name {
			return node
		}
	}
	return nil
}

func splitPath(path string) []string {
	components := strings.Split(path, "/")
	result := make([]string, 0, len(components)-1)
	for _, component := range components {
		if component != "" {
			result = append(result, component)
		}
	}
	return result
}

// DebugDirTree Print the directory tree for debug
func (d *DirTree) DebugDirTree() {
	debugDirTree(d.Root, "", true)
}

func debugDirTree(node *DirTreeNode, indent string, isLastChild bool) {
	fmt.Print(indent)
	if isLastChild {
		fmt.Print("└─ ")
		indent += "   "
	} else {
		fmt.Print("├─ ")
		indent += "│  "
	}

	fmt.Println(node.Name)

	childCount := len(node.Children)
	for i, child := range node.Children {
		isLast := i == childCount-1
		debugDirTree(child, indent, isLast)
	}
}
