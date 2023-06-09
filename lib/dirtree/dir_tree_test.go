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

import "testing"

func TestDirTree_MkDir(t *testing.T) {
	/*
		Build a directory tree as follows
		└─ /
		   ├─ a
		   │  └─ aa
		   │     └─ aaa
		   ├─ b
		   │  └─ bb
		   └─ c
	*/
	dirTree := NewDirTree()
	dirTree.MkDir("/a")
	dirTree.MkDir("/b")
	dirTree.MkDir("/c")
	dirTree.MkDir("/a/aa")
	dirTree.MkDir("/a/aa/aaa")
	dirTree.MkDir("/b/bb")
	dirTree.MkDir("/b/bb")
	dirTree.DebugDirTree()
}
