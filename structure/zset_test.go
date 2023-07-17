package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func initZSetDB() (*ZSetStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestZSetStructure")
	opts.DirPath = dir
	hash, _ := NewZSetStructure(opts)
	return hash, &opts
}

func TestSortedSet(t *testing.T) {
	type test struct {
		name        string
		input       map[string]int
		want        *ZSetNodes
		expectError bool
	}

	zs := newZSetNodes()
	err := zs.InsertNode(3, "banana", "hello")
	err = zs.InsertNode(1, "apple", "hello")
	err = zs.InsertNode(2, "pear", "hello")
	err = zs.InsertNode(44, "orange", "hello")
	err = zs.InsertNode(9, "strawberry", "delish")
	err = zs.InsertNode(15, "dragon-fruit", "nonDelish")
	t.Log(zs.skipList.getRank(9, "strawberry"))
	t.Log(zs.skipList.getNodeByRank(1))
	t.Log(zs.skipList.getNodeByRank(2))
	t.Log(zs.skipList.getNodeByRank(3))
	t.Log(zs.skipList.getNodeByRank(5))
	//var bufEnc bytes.Buffer
	//enc := gob.NewEncoder(&bufEnc)
	//err = enc.Encode(zs)
	//assert.NoError(t, err)
	b, err := zs.Bytes()
	t.Log(b)

	fromBytes := newZSetNodes()
	//buf := bytes.NewBuffer(bufEnc.Bytes())
	//gd := gob.NewDecoder(buf)
	//err = gd.Decode(fromBytes.FromBytes(b))
	//assert.NoError(t, err)

	t.Log(fromBytes.FromBytes(b))
	//t.Log(fromBytes)
	assert.NoError(t, err)

	tests := []test{
		{
			name:        "empty",
			input:       map[string]int{},
			want:        &ZSetNodes{},
			expectError: false,
		},
		{
			name:  "three fruits",
			input: map[string]int{"banana": 3, "apple": 2, "pear": 4, "peach": 40},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, nil)
		})
	}

}

func TestSortedSet_Bytes(t *testing.T) {

}

func TestNewSkipList(t *testing.T) {
	s := newSkipList()

	assert := assert.New(t)
	assert.Equal(1, s.level)
	assert.Nil(s.head.prev)
	assert.Equal(0, s.head.value.score)
	assert.Equal("", s.head.value.member)
}

func TestNewSkipListNode(t *testing.T) {
	score := 10
	key := "test_key"
	value := "test_value"
	level := 5

	node := newSkipListNode(level, score, key, value)

	// Validate node's value
	if node.value.score != score || node.value.member != key || node.value.value != value {
		t.Errorf("Unexpected value in node, got: %v, want: {score: %d, key: %s, val: %s}.\n", node.value, score, key, value)
	}

	// Validate node's level slice length
	if len(node.level) != level {
		t.Errorf("Unexpected length of node's level slice, got: %d, want: %d.\n", len(node.level), level)
	}

	// Validate each SkipListLevel in the level slice
	for _, l := range node.level {
		if l.next != nil || l.span != 0 {
			t.Errorf("Unexpected SkipListLevel, got: %v, want: {forward: nil, span: 0}.\n", l)
		}
	}
}
func TestZAdd(t *testing.T) {
	zs, _ := initZSetDB()
	type testCase struct {
		key    string
		score  int
		member string
		value  string
		err    error
	}

	testCases := []testCase{
		{"key", 10, "member", "value", nil},
		{"", 10, "member", "value", _const.ErrKeyIsEmpty},
	}

	for _, tc := range testCases {
		err := zs.ZAdd(tc.key, tc.score, tc.member, tc.value)
		// Adjust according to your error handling
		if err != tc.err {
			t.Errorf("Expected error to be %v, but got %v", tc.err, err)
		}
	}
}

func TestSkipList_delete(t *testing.T) {
	type deleteTest struct {
		name       string
		score      int
		member     string
		targetList []testZSetNodeValue
		inputList  []testZSetNodeValue
	}

	vals := []testZSetNodeValue{
		{score: 1, member: "mem1", value: nil},
		{score: 2, member: "mem2", value: nil},
		{score: 3, member: "mem3", value: nil},
		{score: 4, member: "mem4", value: nil},
		{score: 5, member: "mem5", value: nil},
	}

	// Omitted: Add some nodes into sl...

	tests := []deleteTest{
		{
			name:       "Delete Test 1",
			score:      15,
			member:     "member1",
			targetList: []testZSetNodeValue{{score: 3, member: "mem3"}}, // result of adding nodes into sl
			inputList:  vals,
		},
		// Add more test cases here...
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			head := newZSetNodes()
			populateSkipListFromSlice(head, test.inputList)

			for _, value := range test.targetList {
				// check if the insertion has been performed
				assert.True(t, head.exists(value.score, value.member))
				// delete the target members
				assert.NoError(t, head.RemoveNode(value.member))
				// check to see if the deletion has been correctly performed
				assert.False(t, head.exists(value.score, value.member))

			}
		})
	}
}

type testZSetNodeValue struct {
	score  int
	member string
	value  interface{}
}

func populateSkipListFromSlice(nodes *ZSetNodes, zSetNodeValues []testZSetNodeValue) {
	// Iterate over the zsetNodes array
	for _, zSetNode := range zSetNodeValues {
		_ = nodes.InsertNode(zSetNode.score, zSetNode.member, zSetNode.value)
	}
}
