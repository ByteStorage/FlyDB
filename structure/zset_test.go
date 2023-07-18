package structure

import (
	"github.com/ByteStorage/FlyDB/config"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
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
		want        *FZSet
		expectError bool
	}

	zs := newZSetNodes()
	err := zs.InsertNode(3, "banana", "hello")
	err = zs.InsertNode(1, "apple", "hello")
	err = zs.InsertNode(2, "pear", "hello")
	err = zs.InsertNode(44, "orange", "hello")
	err = zs.InsertNode(9, "strawberry", "delish")
	err = zs.InsertNode(15, "dragon-fruit", "nonDelish")
	b, err := zs.Bytes()

	fromBytes := newZSetNodes()
	err = fromBytes.FromBytes(b)
	assert.NoError(t, err)
	assert.NotNil(t, fromBytes.skipList)
	assert.Equal(t, fromBytes.size, zs.size)
	tests := []test{
		{
			name:        "empty",
			input:       map[string]int{},
			want:        &FZSet{},
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

func TestZRem(t *testing.T) {
	mockZSetStructure, _ := initZSetDB()

	// 1. Test for Key is Empty
	err := mockZSetStructure.ZRem("", "member")
	require.Error(t, err)
	require.Equal(t, _const.ErrKeyIsEmpty, err)
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
		t.Run(tc.key, func(t *testing.T) {
			err := mockZSetStructure.ZAdd(tc.key, tc.score, tc.member, tc.value)
			// check to see if element added
			assert.Equal(t, tc.err, err)
			if tc.err == nil {
				// check if member added
				assert.True(t, mockZSetStructure.exists(tc.key, tc.score, tc.member))
			}

		})
	}
}
func TestZRems(t *testing.T) {
	mockZSetStructure, _ := initZSetDB()

	// 1. Test for Key is Empty
	err := mockZSetStructure.ZRems("", "member")
	require.Error(t, err)
	require.Equal(t, _const.ErrKeyIsEmpty, err)
	type testCase struct {
		key   string
		input []ZSetValue
		rems  []string
		want  []ZSetValue
		err   error
	}

	testCases := []testCase{
		{"key",
			[]ZSetValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
				{score: 6, member: "mem6", value: ""},
			},
			[]string{
				"mem0",
				"mem1",
				"mem6",
			}, []ZSetValue{
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
			}, nil},
		{"",
			[]ZSetValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
			},
			[]string{
				"mem0",
				"mem1",
				"mem2",
				"mem3",
				"mem4",
				"mem5",
				"mem6",
			},
			[]ZSetValue{},
			_const.ErrKeyIsEmpty},
		{
			"Key1",
			[]ZSetValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
			},
			[]string{
				"mem0",
				"mem1",
				"mem2",
				"mem3",
				"mem4",
				"mem5",
				"mem6",
			}, []ZSetValue{}, _const.ErrKeyNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			_ = mockZSetStructure.ZAdds(tc.key, tc.input...)

			//remove all the elements
			err = mockZSetStructure.ZRems(tc.key, tc.rems...)
			assert.Equal(t, tc.err, err)
			//validate
			for _, value := range tc.want {
				te := mockZSetStructure.exists(tc.key, value.score, value.member)
				assert.True(t, te)
			}
		})
	}
}
func TestZAdd(t *testing.T) {
	zs, _ := initZSetDB()
	type testCase struct {
		key    string
		score  int
		member string
		value  string
		want   ZSetValue
		err    error
	}

	testCases := []testCase{
		{
			"key",
			10,
			"member",
			"value",
			ZSetValue{member: "member"},
			nil,
		},
		{
			"",
			10,
			"member",
			"value",
			ZSetValue{member: ""},
			_const.ErrKeyIsEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			err := zs.ZAdd(tc.key, tc.score, tc.member, tc.value)
			assert.Equal(t, tc.err, err)
			if tc.err == nil {
				// check if member added
				assert.True(t, zs.exists(tc.key, tc.score, tc.member))
				err = zs.ZRem(tc.key, tc.member)
				assert.NoError(t, err)
				// should be removed successfully
				assert.False(t, zs.exists(tc.key, tc.score, tc.member))
			}
			// Adjust according to your error handling

		})
	}
}
func TestZAdds(t *testing.T) {
	zs, _ := initZSetDB()

	// 1. Test for Key is Empty
	err := zs.ZAdds("", []ZSetValue{}...)
	require.Error(t, err)
	require.Equal(t, _const.ErrKeyIsEmpty, err)
	type testCase struct {
		key   string
		input []ZSetValue
		want  []ZSetValue
		err   error
	}

	testCases := []testCase{
		{"key",
			[]ZSetValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
				{score: 6, member: "mem6", value: ""},
			},
			[]ZSetValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
				{score: 6, member: "mem6", value: ""},
			},
			nil},
		{"",
			[]ZSetValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
			},
			[]ZSetValue{},
			_const.ErrKeyIsEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			err = zs.ZAdds(tc.key, tc.input...)
			assert.Equal(t, tc.err, err)
			//validate
			for _, value := range tc.want {
				te := zs.exists(tc.key, value.score, value.member)
				assert.True(t, te)
			}
		})
	}
}
func TestZIncrBy(t *testing.T) {
	zs, _ := initZSetDB()
	err := zs.ZIncrBy("", "non-existingMember", 5)
	if err == nil {
		t.Error("Expected error for empty key not returned")
	}

	err = zs.ZIncrBy("key", "non-existingMember", 5)
	if !assert.ErrorIs(t, err, _const.ErrKeyNotFound) {
		t.Error("Expected ErrKeyNotFound for non-existing member not returned")
	}
	err = zs.ZAdd("key", 1, "existingMember", "")
	assert.NoError(t, err)
	err = zs.ZIncrBy("key", "existingMember", 5)
	assert.NoError(t, err)

	Zset, err := zs.getZSetFromDB(stringToBytesWithKey("key"))
	assert.Equal(t, 6, Zset.dict["existingMember"].score)
}
func TestZRank(t *testing.T) {
	zs, _ := initZSetDB()

	// Assume that ZAdd adds a member to a set and assigns the member a score.
	// Here the score does not matter
	err := zs.ZAdd("myKey", 1, "member1", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 2, "member2", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 3, "member3", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 4, "member4", "")
	assert.NoError(t, err)

	// Test when member is present in the set
	rank, err := zs.ZRank("myKey", "member1")
	assert.NoError(t, err)   // no error should occur
	assert.Equal(t, 1, rank) // as we inserted 'member1' first, its rank should be 1

	// Test when member is not present in the set
	rank, err = zs.ZRank("myKey", "unavailableMember")
	assert.Error(t, err)     // an error should occur
	assert.Equal(t, 0, rank) // as 'unavailableMember' is not part of set, rank should be 0

	// Test with an empty key
	rank, err = zs.ZRank("", "member")
	assert.Error(t, err)     // an error should occur
	assert.Equal(t, 0, rank) // rank should be 0 for invalid key}

	// Test member2 which should be 2nd
	rank, err = zs.ZRank("myKey", "member2")
	assert.NoError(t, err)   // there should be no errors
	assert.Equal(t, 2, rank) // rank should be 2 for key `member2`

	// Test member3 which should be 3rd
	rank, err = zs.ZRank("myKey", "member3")
	assert.NoError(t, err) // there should be no errors
	assert.Equal(t, 3, rank)

	// remove member2 and test `member3` which should become 2
	err = zs.ZRem("myKey", "member2")
	assert.NoError(t, err) // there should be no errors
	rank, err = zs.ZRank("myKey", "member3")
	assert.NoError(t, err)   // there should be no errors
	assert.Equal(t, 2, rank) // now `member3` should become 2nd
}
func TestZRevRank(t *testing.T) {
	zs, _ := initZSetDB()

	// Assume that ZAdd adds a member to a set and assigns the member a score.
	// Here the score does not matter
	err := zs.ZAdd("myKey", 1, "member1", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 2, "member2", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 3, "member3", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 4, "member4", "")
	assert.NoError(t, err)

	// Test when member is present in the set
	rank, err := zs.ZRevRank("myKey", "member3")
	assert.NoError(t, err)   // no error should occur
	assert.Equal(t, 2, rank) // as we inserted 'member1' first, its rank should be 1

	// Test when member is not present in the set
	rank, err = zs.ZRevRank("myKey", "unavailableMember")
	assert.Error(t, err)     // an error should occur
	assert.Equal(t, 0, rank) // as 'unavailableMember' is not part of set, rank should be 0

	// Test with an empty key
	rank, err = zs.ZRevRank("", "member")
	assert.Error(t, err)     // an error should occur
	assert.Equal(t, 0, rank) // rank should be 0 for invalid key}

	// Test member2 which should be 2nd
	rank, err = zs.ZRevRank("myKey", "member1")
	assert.NoError(t, err)   // there should be no errors
	assert.Equal(t, 4, rank) // rank should be 2 for key `member2`

	// Test member3 which should be 3rd
	rank, err = zs.ZRevRank("myKey", "member4")
	assert.NoError(t, err) // there should be no errors
	assert.Equal(t, 1, rank)

	// remove member2 and test `member3` which should become 2
	err = zs.ZRem("myKey", "member2")
	assert.NoError(t, err) // there should be no errors
	rank, err = zs.ZRevRank("myKey", "member3")
	assert.NoError(t, err)   // there should be no errors
	assert.Equal(t, 2, rank) // now `member3` should become 2nd
}
func TestZRevRange(t *testing.T) {
	zs, _ := initZSetDB()

	err := zs.ZAdd("myKey", 1, "member1", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 2, "member2", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 3, "member3", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 4, "member4", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 5, "member5", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 6, "member6", "")
	assert.NoError(t, err)

	var n []uint8
	tests := []struct {
		key     string
		start   int
		end     int
		want    []ZSetValue
		wantErr error
	}{
		{"myKey", 0, 3, []ZSetValue{
			{6, "member6", n},
			{5, "member5", n},
			{4, "member4", n},
		}, nil},
		{"", 0, 2, nil, _const.ErrKeyIsEmpty},
		{"fail", 0, 2, nil, _const.ErrKeyNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, err := zs.ZRevRange(tt.key, tt.start, tt.end)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ZRange() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestZRange(t *testing.T) {
	zs, _ := initZSetDB()

	err := zs.ZAdd("myKey", 1, "member1", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 2, "member2", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 3, "member3", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 4, "member4", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 5, "member5", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 6, "member6", "")
	assert.NoError(t, err)
	var n []uint8
	tests := []struct {
		key     string
		start   int
		end     int
		want    []ZSetValue
		wantErr error
	}{
		{"myKey", 0, 3, []ZSetValue{
			{1, "member1", n},
			{2, "member2", n},
			{3, "member3", n},
		}, nil},
		{"", 0, 2, nil, _const.ErrKeyIsEmpty},
		{"fail", 0, 2, nil, _const.ErrKeyNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, err := zs.ZRange(tt.key, tt.start, tt.end)

			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ZRange() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestZCard(t *testing.T) {
	zs, _ := initZSetDB()

	err := zs.ZAdd("myKey", 1, "member1", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 2, "member2", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 3, "member3", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 4, "member4", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 5, "member5", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 6, "member6", "")
	assert.NoError(t, err)
	tests := []struct {
		name    string
		key     string
		want    int
		wantErr error
	}{
		{"Empty Key", "", 0, _const.ErrKeyIsEmpty},
		{"Non-Existent Key", "nonExist", 0, _const.ErrKeyNotFound},
		{"Existing Key", "myKey", 6, nil},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got, err := zs.ZCard(tt.key)

			if tt.want != got {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}

			if tt.wantErr != nil {
				if err == nil || tt.wantErr.Error() != err.Error() {
					t.Fatalf("expected error '%v', got '%v'", tt.wantErr, err)
				}

			} else if err != nil {
				t.Fatalf("expected no error, got error '%v'", err)
			}
		})
	}
}
func TestZScore(t *testing.T) {
	zs, _ := initZSetDB()

	err := zs.ZAdd("myKey", 1, "member1", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 2, "member2", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 3, "member3", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 4, "member4", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 5, "member5", "")
	assert.NoError(t, err)

	err = zs.ZAdd("myKey", 6, "member6", "")
	assert.NoError(t, err)
	tests := []struct {
		expectError   error
		expectedScore int
		key           string
		member        string
	}{
		{_const.ErrKeyIsEmpty, 0, "", "member1"},
		{_const.ErrKeyNotFound, 0, "key1", "foo"},
		{nil, 1, "myKey", "member1"},
		{nil, 2, "myKey", "member2"},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			score, err := zs.ZScore(test.key, test.member)
			assert.Equal(t, test.expectError, err)
			assert.Equal(t, test.expectedScore, score)
		})
	}
}
func TestNewSkipList(t *testing.T) {
	s := newSkipList()

	assertions := assert.New(t)
	assertions.Equal(1, s.level)
	assertions.Nil(s.head.prev)
	assertions.Equal(0, s.head.value.score)
	assertions.Equal("", s.head.value.member)
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

func populateSkipListFromSlice(nodes *FZSet, zSetNodeValues []testZSetNodeValue) {
	// Iterate over the zsetNodes array
	for _, zSetNode := range zSetNodeValues {
		_ = nodes.InsertNode(zSetNode.score, zSetNode.member, zSetNode.value)
	}
}
func TestRandomLevel(t *testing.T) {
	for i := 0; i < 1000; i++ {
		level := randomLevel()
		if level < 1 || level > SKIPLIST_MAX_LEVEL {
			t.Errorf("Generated level out of range: %v", level)
		}
	}
}
func TestZSetNodes_InsertNode(t *testing.T) {
	pq := &FZSet{}

	// Case 1: Insert new node
	err := pq.InsertNode(1, "test", "value")
	if err != nil {
		t.Error("Failed when inserting a new node")
	}

	if _, ok := pq.dict["test"]; !ok {
		t.Error("Insert node failed, expected key to exist in dictionary")
	}

	// Case 2: Update existing node with same score
	err = pq.InsertNode(1, "test", "newvalue")
	if err != nil {
		t.Error("Failed when updating a score with same value")
	}

	if v, ok := pq.dict["test"]; !ok || v.value != "newvalue" {
		t.Error("Update node failed, expected value to be updated")
	}

	// Case 3: Insert node with existing key but different score
	err = pq.InsertNode(2, "test", "newvalue")
	if err != nil {
		t.Error("Failed when updating a score with a new value")
	}

	if v, ok := pq.dict["test"]; !ok || v.score != 2 {
		t.Error("Update node failed, expected score to be updated")
	}
}
func TestZCount(t *testing.T) {
	zs, _ := initZSetDB()

	tests := []struct {
		key   string
		input []testZSetNodeValue
		min   int
		max   int
		want  int
		err   error
	}{
		{
			"test1",
			[]testZSetNodeValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
				{score: 6, member: "mem6", value: ""},
			},
			1, 5, 5, nil,
		},
		{
			"test2",
			[]testZSetNodeValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
				{score: 6, member: "mem6", value: ""},
			},
			0, 5, 6, nil,
		},
		{
			"test3",
			[]testZSetNodeValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
				{score: 6, member: "mem6", value: ""},
			},
			1, 3, 3, nil,
		},
		{
			"test4",
			[]testZSetNodeValue{
				{score: 0, member: "mem0", value: ""},
				{score: 1, member: "mem1", value: ""},
				{score: 2, member: "mem2", value: ""},
				{score: 3, member: "mem3", value: ""},
				{score: 4, member: "mem4", value: ""},
				{score: 5, member: "mem5", value: ""},
				{score: 6, member: "mem6", value: ""},
			},
			2, 2, 1, nil,
		},
		{
			"test5",
			[]testZSetNodeValue{
				{score: 3, member: "mem3", value: ""},
			},
			10, 20, 0, nil,
		},
		{
			"test6",
			[]testZSetNodeValue{
				{score: 3, member: "mem3", value: ""},
			},
			10, 5, 0, ErrInvalidArgs,
		},
		{
			"",
			[]testZSetNodeValue{
				{score: 3, member: "mem3", value: ""},
			},
			10, 5, 0, _const.ErrKeyIsEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			for _, value := range tt.input {
				_ = zs.ZAdd(tt.key, value.score, value.member, value.value.(string))
			}
			got, err := zs.ZCount(tt.key, tt.min, tt.max)
			if got != tt.want {
				t.Errorf("TestZCount(%v, %v, %v) = %v, want: %v", tt.key, tt.min, tt.max, got, tt.want)
			}
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("TestZCount(%v, %v, %v) returned unexpected error: got %v, want: %v", tt.key, tt.min, tt.max, err, tt.err)
			}
		})
	}
}

func TestFZSetMinMax(t *testing.T) {
	fzs := &FZSet{
		skipList: newSkipList(),
	}
	_ = fzs.InsertNode(1, "mem1", "")
	_ = fzs.InsertNode(100, "mem2", "")

	minScore, maxScore := fzs.getMinMaxScore()

	if minScore != 1 || maxScore != 100 {
		t.Errorf("getMinMaxScore() = %d, %d, want: 1, 100", minScore, maxScore)
	}

	if min := fzs.min(5, 10); min != 5 {
		t.Errorf("min(5, 10) = %d, want: 5", min)
	}

	if max := fzs.max(5, 10); max != 10 {
		t.Errorf("max(5, 10) = %d, want: 10", max)
	}
}

func Test_exists(t *testing.T) {
	zs, _ := initZSetDB()

	tt := []struct {
		key    string
		score  int
		member string
		want   bool
	}{
		{
			key:    "",
			score:  1,
			member: "",
			want:   false,
		},
	}

	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			got := zs.exists(tc.key, tc.score, tc.member)

			if got != tc.want {
				t.Errorf("exists() = %v, want %v", got, tc.want)
			}
		})
	}
}
