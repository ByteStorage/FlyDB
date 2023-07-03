package dirtree

import (
	"fmt"
	"path/filepath"
	"strings"
)

var _ DirTreeInterface = &DirTree{}

/**TODO:
The tentative plan here is to maintain only one directory tree
				/
	  / Node 1 / Node 2 / Node 3
	  /file1  /file2  /file3
Roughly structured as above, each node is a directory,
There are files under each directory, so you only need to maintain one directory tree,
And do not need to maintain the node and file mapping relationship,
here DirTree just write a simple directory tree,
you can find there is no optimized directory tree
*/

// DirTree is a struct that represents a directory tree
type DirTree struct {
	// Root is the root node of the directory tree
	Root    *DirTreeNode
	NodeMap map[string]*DirTreeNode // quickly find node when deleting
}

// DirTreeNode is a struct that represents a directory tree node
type DirTreeNode struct {
	path     string
	name     string
	children map[string]*DirTreeNode
	isDir    bool
}

// NewDirTree Creates a new instance of DirTree
func NewDirTree() *DirTree {
	nodeMap := make(map[string]*DirTreeNode)
	root := &DirTreeNode{
		path:     "/",
		name:     "/",
		children: make(map[string]*DirTreeNode),
		isDir:    true,
	}
	nodeMap[root.name] = root
	return &DirTree{
		Root:    root,
		NodeMap: nodeMap,
	}
}

func (d *DirTree) InsertFile(path string) bool {
	// the file maybe already exists
	if _, ok := d.NodeMap[path]; ok {
		return false
	}

	components := splitPath(path)

	parentPath := filepath.Dir(path)

	// get the parent node of the file
	parentNode, ok := d.NodeMap[parentPath]
	if !ok {
		// the directory of this file maybe not exists, so create directories
		d.MkDir(filepath.Join(components[:len(components)-1]...))
		parentNode = d.NodeMap[parentPath]
	}

	fileName := "/" + components[len(components)-1]

	fileNode := &DirTreeNode{
		path:     path,
		name:     fileName,
		children: nil,
		isDir:    false,
	}

	parentNode.children[fileName] = fileNode
	d.NodeMap[path] = fileNode

	return true
}

func (d *DirTree) DeleteFile(path string) bool {
	if _, ok := d.NodeMap[path]; !ok {
		return false
	}

	components := splitPath(path)
	fileName := "/" + components[len(components)-1]

	parentNode := d.findParentNode(path)
	delete(parentNode.children, fileName)
	delete(d.NodeMap, path)

	return true
}

func (d *DirTree) DeleteDir(path string) bool {
	node, ok := d.NodeMap[path] // find nodes directly from nodeMap, avoiding recursion overhead
	if !ok {
		return false // directory not found
	}

	if !node.isDir {
		return false // this is not a path to a directory
	}

	// recursively delete child nodes
	for name, child := range node.children {
		d.deleteNode(child)
		delete(node.children, name)
	}

	parentNode := d.findParentNode(path)

	delete(parentNode.children, node.name)
	delete(d.NodeMap, node.path)

	return true
}

// MkDir Inserts a new DirTreeNode into the DirTree
func (d *DirTree) MkDir(path string) bool {

	components := splitPath(path)

	curNode := d.Root

	for i, component := range components {
		component = "/" + component
		child, ok := curNode.children[component]
		// if the current component does not exist, create a new one
		if !ok {
			nodePath := "/" + filepath.Join(components[:i+1]...)

			newNode := &DirTreeNode{
				path:     nodePath,
				name:     component,
				children: make(map[string]*DirTreeNode),
				isDir:    true,
			}
			curNode.children[component] = newNode
			curNode = newNode

			d.NodeMap[nodePath] = newNode
		} else {
			// otherwise go to the next level
			curNode = child
		}
	}

	return true
}

func (d *DirTree) deleteNode(node *DirTreeNode) {
	if !node.isDir {
		// Delete file node
		delete(d.NodeMap, node.path)
	} else {
		// Delete directory node and its descendants
		for name, child := range node.children {
			d.deleteNode(child)
			delete(node.children, name)
		}
		delete(d.NodeMap, node.path)
	}
}

func (d *DirTree) findParentNode(path string) *DirTreeNode {
	parentPath := filepath.Dir(path)
	return d.NodeMap[parentPath]
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

	fmt.Println(node.name)

	childCount := len(node.children)
	i := 0
	for _, child := range node.children {
		isLast := i == childCount-1
		debugDirTree(child, indent, isLast)
		i++
	}
}
