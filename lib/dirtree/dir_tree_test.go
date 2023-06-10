// Copyright 2023 The promisedb Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dirtree

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestDirTree_MkDir(t *testing.T) {
	/*
		Build a directory tree as follows
		└─ /
		   ├─ b
		   │  └─ bb
		   ├─ c
		   └─ a
			  └─ aa
				 └─ aaa
	*/
	dirTree := NewDirTree()

	// normally mkdir
	dirTree.MkDir("/a")
	dirTree.MkDir("/b")
	dirTree.MkDir("/c")

	// insert sub dir
	dirTree.MkDir("/a/aa")
	dirTree.MkDir("/a/aa/aaa")

	// insert duplicate folder
	dirTree.MkDir("/b/bb")
	dirTree.MkDir("/b/bb")

	dirTree.DebugDirTree()
}

func TestDirTree_MkDir_Batch(t *testing.T) {
	dirTree := NewDirTree()

	dirCount := 10000

	start := time.Now()

	for i := 0; i < dirCount; i++ {
		filePath := generateRandomPath(true)
		dirTree.MkDir(filePath)
	}

	elapsed := time.Since(start)

	fmt.Printf("Inserted %d dirs in %s\n", dirCount, elapsed)
}

func TestDirTree_DeleteDir(t *testing.T) {

	dirTree := NewDirTree()

	dirTree.MkDir("/a")
	dirTree.MkDir("/a/aa")
	dirTree.MkDir("/b")
	dirTree.MkDir("/b/bb")

	// normally delete dir
	assert.True(t, dirTree.DeleteDir("/a"))

	// delete child node
	assert.True(t, dirTree.DeleteDir("/b/bb"))
	assert.True(t, dirTree.DeleteDir("/b"))

	// remove non-existing node
	assert.False(t, dirTree.DeleteDir("/c"))

	dirTree.DebugDirTree()
}

func TestDirTree_DeleteDir_Batch(t *testing.T) {
	dirTree := NewDirTree()

	dirCount := 10000

	for i := 0; i < dirCount; i++ {
		filePath := generateRandomPath(true)
		dirTree.MkDir(filePath)
	}

	start := time.Now()

	assert.True(t, dirTree.DeleteDir("/"))

	elapsed := time.Since(start)

	fmt.Printf("deleted %d dirs in %s\n", dirCount, elapsed)

	dirTree.DebugDirTree()
}

func TestDirTree_InsertFile(t *testing.T) {
	dirTree := NewDirTree()

	dirTree.MkDir("/a")
	dirTree.MkDir("/b/bb")

	// normally insert a file
	assert.True(t, dirTree.InsertFile("/test.txt"))

	// insert a duplicate file
	assert.False(t, dirTree.InsertFile("/test.txt"))

	// insert two files into children nodes
	assert.True(t, dirTree.InsertFile("/a/test.txt"))
	assert.True(t, dirTree.InsertFile("/b/bb/test.txt"))

	// insert a file that it's parent node does not exist
	assert.True(t, dirTree.InsertFile("/c/test.txt"))

	dirTree.DebugDirTree()
}

func TestDirTree_InsertFile_Batch(t *testing.T) {
	dirTree := NewDirTree()

	fileCount := 100000

	start := time.Now()

	for i := 0; i < fileCount; i++ {
		filePath := generateRandomPath(false)
		dirTree.InsertFile(filePath)
	}

	elapsed := time.Since(start)

	fmt.Printf("Inserted %d files in %s\n", fileCount, elapsed)
}

func TestDirTree_DeleteFile(t *testing.T) {
	dirTree := NewDirTree()

	dirTree.MkDir("/a")
	dirTree.MkDir("/b/bb")

	// create a directory tree
	assert.True(t, dirTree.InsertFile("/test.txt"))
	assert.True(t, dirTree.InsertFile("/a/test.txt"))
	assert.True(t, dirTree.InsertFile("/b/bb/test.txt"))
	assert.True(t, dirTree.InsertFile("/c/cc/ccc/test.txt"))

	// normally delete a file
	assert.True(t, dirTree.DeleteFile("/test.txt"))

	// delete a file that does not exist
	assert.False(t, dirTree.DeleteFile("/test.txt"))

	// delete files in children nodes
	assert.True(t, dirTree.DeleteFile("/b/bb/test.txt"))
	assert.True(t, dirTree.DeleteFile("/c/cc/ccc/test.txt"))

	dirTree.DebugDirTree()
}

func TestDirTree_DeleteFile_Batch(t *testing.T) {
	dirTree := NewDirTree()

	fileCount := 100000

	filePaths := make([]string, fileCount)

	for i := 0; i < fileCount; i++ {
		filePath := generateRandomPath(false)
		dirTree.InsertFile(filePath)
		filePaths[i] = filePath
	}

	start := time.Now()

	for i := 0; i < fileCount; i++ {
		dirTree.DeleteFile(filePaths[i])
	}

	elapsed := time.Since(start)

	fmt.Printf("deleted %d files in %s\n", fileCount, elapsed)
}

// generateRandomPath generate random paths from levels 1-20
func generateRandomPath(isDir bool) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	pathLength := rand.Intn(20) + 1
	components := make([]string, pathLength)

	for i := 0; i < pathLength; i++ {
		componentLength := rand.Intn(10) + 1
		component := ""

		for j := 0; j < componentLength; j++ {
			component += string(chars[rand.Intn(len(chars))])
		}

		components[i] = component
	}

	path := "/" + strings.Join(components, "/")
	if !isDir {
		path = path + ".txt"
	}

	return path
}
