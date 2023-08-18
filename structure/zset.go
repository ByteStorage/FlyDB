package structure

import (
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/encoding"
	"math"
	"math/rand"
	"time"
)

const (
	// SKIPLIST_MAX_LEVEL is better to be log(n) for the best performance.
	SKIPLIST_MAX_LEVEL = 10   //
	SKIPLIST_PROB      = 0.25 // SkipList Probability
)

/**
ZSet or Sorted Set structure is borrowed from Redis' implementation, the Redis implementation
utilizes a SkipList and a dictionary.
*/

// ZSetStructure is a structure for ZSet or SortedSet
type ZSetStructure struct {
	db *engine.DB
}

// FZSet represents a specific data structure in the database, which is key to handling sorted sets (ZSets).
// This struct facilitates interactions with data stored in the sorted set, allowing for both complex and simple operations.
//
// It contains three struct fields:
//
//   - 'dict': A Go map with string keys and pointers to ZSetValue values. This map aims to provide quick access to
//     individual values in the sorted set based on the provided key.
//
//   - 'size': An integer value representing the current size (number of elements) in the FZSet struct. This information is efficiently
//     kept track of whenever elements are added or removed from the set, so no separate computation is needed to retrieve this information.
//
//   - 'skipList': A pointer towards a SkipList struct. SkipLists perform well under numerous operations, such as insertion, deletion, and searching. They are
//     a crucial component in maintaining the sorted set in a practical manner. In this context, the SkipList is used to keep an ordered track of the elements
//     in the FZSet struct.
type FZSet struct {
	// dict field is a map where the key is a string and
	// the value is a pointer to ZSetValue instances,
	// codified with the tag "dict".
	dict map[string]*ZSetValue `codec:"dict"`

	// size field represents the quantity of elements within
	// the structure, codified with the tag "size".
	size int `codec:"size"`

	// skipList field is a pointer to an object of type SkipList,
	// codified with the tag "skip_list".
	skipList *SkipList `codec:"skip_list"`
}

// SkipList represents a skip list data structure, an ordered list with a hierarchical
// structure that allows for fast search and insertion of elements.
type SkipList struct {
	// level represents the highest level of the skip list.
	level int

	// head refers to the first node in the skip list.
	head *SkipListNode

	// tail refers to the last node in the skip list.
	tail *SkipListNode

	// size represents the total number of nodes in the skip list (excluding head and tail nodes).
	size int
}

// SkipListLevel is a structure encapsulating a single level in a skip list data structure.
// It contains two struct fields:
// - 'next': A pointer to the next SkipListNode in the current level.
// - 'span': An integer representing the span size of this SkipListLevel. The span is the number of nodes between the current node
// and the node to which the next pointer is pointing in the skip list.
type SkipListLevel struct {
	next *SkipListNode
	span int
}

// SkipListNode represents a single node in a SkipList structure.
// It is built with three elements:
//   - 'prev': This is a pointer to the previous node in the skip list. Together with the 'next' pointers in the SkipListNodeLevel,
//     it forms a network of nodes, where traversal of the skip list is possible both forwards and backwards.
//   - 'level': This is an array (slice) of pointers towards SkipListLevel structures. Each element corresponds to a level of the skip list,
//     embedding the 'next' node at that same level, and the span between the current node and that 'next' node.
//   - 'value': This is a pointer towards a single ZSetValue structure. It holds the actual payload of the node
//     (namely the 'score', 'key', and 'value' properties used in the context of Redis Sorted Sets), as well as provides the basis for ordering of nodes in the skip list.
type SkipListNode struct {
	// prev is a pointer to the previous node in the skip list.
	prev *SkipListNode

	// level is a slice of pointers to SkipListLevel.
	// Each level represents a forward pointer to the next node in the current list level.
	level []*SkipListLevel

	// value is a pointer to the ZSetValue.
	// This represents the value that this node holds.
	value *ZSetValue
}

// ZSetValue is a struct used in the SkipList data structure. In the context of Redis Sorted Set (ZSet) implementation,
// it represents a single node value in the skip list. A ZSetValue has three members:
// - 'score' which is an integer representing the score of the node. Nodes in a skip list are ordered by this score in ascending order.
// - 'member' which is a string defining the key of the node. For nodes with equal scores, order is determined with lexicographical comparison of keys.
// - 'value' which is an interface{}, meaning it can hold any data type. This represents the actual value of the node in the skip list.
type ZSetValue struct {
	// Score is typically used for sorting purposes. Nodes with higher scores will be placed higher in the skip list.
	score int

	// member represents the unique identifier for each node.
	member string

	// value is the actual content/data that is being stored in the node.
	value interface{}
}

// randomLevel is a function that generates a probabilistic level for a node in a SkipList data structure.
// The goal is to diversify the level distribution and contribute to achieving an ideal skiplist performance.
// Function has no parameters.
// The process starts with two initial variables:
//   - 'level' which starts from 1,
//   - 'thresh' which is a product of the constant skiplist probability 'SKIPLIST_PROB' and bitwise mask: 0xFFF, taken to the nearest integer.
//
// In an infinite loop, a random 31-bit integer value is generated, bitwise-and is computed with 0xFFF and compared with 'thresh'.
// If the result is smaller, 'level' is incremented by one. Otherwise, the loop is exited.
// Finally, the function checks the calculated level against the maximum allowed skiplist level 'SKIPLIST_MAX_LEVEL'.
// If 'level' is greater, 'SKIPLIST_MAX_LEVEL' is returned, otherwise the calculated 'level' value is returned.
// The function returns an integer which will be the level of new node in skiplist.
func randomLevel() int {
	// Initialize level to 1
	level := 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Calculate the threshold for level. It's derived from the probability constant of the skip list.
	thresh := int(math.Round(SKIPLIST_PROB * 0xFFF))

	// While a randomly generated number is less than this threshold, increment the level.
	for int(r.Int31()&0xFFF) < thresh {
		level++
	}

	// Check if the level is more than the maximum allowed level for the skip list
	// If it is, return the maximum level. Otherwise, return the generated level.
	if level > SKIPLIST_MAX_LEVEL {
		return SKIPLIST_MAX_LEVEL
	} else {
		return level
	}
}

// NewZSetStructure Returns a new ZSetStructure
func NewZSetStructure(options config.Options) (*ZSetStructure, error) {
	db, err := engine.NewDB(options)
	if err != nil {
		return nil, err
	}
	return &ZSetStructure{db: db}, nil
}

// newZSetNodes is a function that creates a new FZSet object and returns a pointer to it.
// It initializes the dictionary member dict of the newly created object to an empty map.
// The map is intended to map strings to pointers of ZSetValue objects.
// size member of the object is set to 0, indicating that the FZSet object is currently empty.
// The skipList member of the object is set to a new SkipList object created by calling `newSkipList()` function.
func newZSetNodes() *FZSet {
	return &FZSet{
		dict:     make(map[string]*ZSetValue),
		size:     0,
		skipList: newSkipList(),
	}
}

// newSkipList is a function that creates an instance of a SkipList struct object and returns a pointer to it.
// This involves initializing the level of the SkipList to 1 and creating a new SkipListNode object as the head of the list.
// The head node is constructed with a level set to SKIPLIST_MAX_LEVEL, key and value as empty string and value as nil respectively.
func newSkipList() *SkipList {
	return &SkipList{
		level: 1,
		head:  newSkipListNode(SKIPLIST_MAX_LEVEL, 0, "", nil),
	}
}

// newSkipListNode is a function that takes integer as level, score and a string as key along with a value of any type.
// It returns a pointer to a SkipListNode. This function is responsible for creating a new SkipListNode with provided level, score,
// key, and value. After creating the node, it initializes every level of the node with an empty SkipListLevel object.
// In the context of a skip list data structure, this function serves as a helper function for creating new nodes to be inserted to the list.
func newSkipListNode(level int, score int, key string, value interface{}) *SkipListNode {
	// Create a new SkipListNode with specified score, key, value and a slice of
	// SkipListLevel with length equal to specified level
	node := &SkipListNode{
		value: newSkipListNodeValue(score, key, value),
		level: make([]*SkipListLevel, level),
	}

	// Initialize each SkipListLevel in the level slice
	for i := range node.level {
		node.level[i] = new(SkipListLevel)
	}
	// Returning the pointer to the created node
	return node
}

// newSkipListNodeValue is a function that constructs and returns a new ZSetValue.
// It takes a score (int), a key (string), and a value (interface{}) as parameters.
// These parameters serve as the initial state of the ZSetValue upon its creation.
func newSkipListNodeValue(score int, member string, value interface{}) *ZSetValue {
	// Create a new instance of a ZSetValue with the provided score, key, and value.
	node := &ZSetValue{
		score:  score,
		member: member,
		value:  value,
	}

	// Return the newly created ZSetValue.
	return node
}

// insert is a method of the SkipList type that is used to insert a new node into the skip list. It takes as arguments
// the score (int), key (string) and a value (interface{}), and returns a pointer to the ZSetValue struct. The method
// organizes nodes in the list based on the score in ascending order. If two nodes have the same score, they will be arranged
// based on the key value. The method also assigns span values to the levels in the skip list.
func (sl *SkipList) insert(score int, key string, value interface{}) *ZSetValue {
	update := make([]*SkipListNode, SKIPLIST_MAX_LEVEL)
	rank := make([]int, SKIPLIST_MAX_LEVEL)
	node := sl.head

	// Go from highest level to lowest
	for i := sl.level - 1; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		if sl.level-1 == i {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		if node.level[i] != nil {
			for node.level[i].next != nil &&
				(node.level[i].next.value.score < score ||
					(node.level[i].next.value.score == score && // score is the same but the key is different
						node.level[i].next.value.member < key)) {
				rank[i] += node.level[i].span
				node = node.level[i].next
			}
		}
		update[i] = node
	}
	level := randomLevel()
	// add a new level
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			rank[i] = 0
			update[i] = sl.head
			update[i].level[i].span = sl.size
		}
		sl.level = level
	}
	node = newSkipListNode(level, score, key, value)

	for i := 0; i < level; i++ {
		node.level[i].next = update[i].level[i].next
		update[i].level[i].next = node
		// update span covered by update[i] as newNode is inserted here
		node.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}
	// increment span for untouched levels
	for i := level; i < sl.level; i++ {
		update[i].level[i].span++
	}
	// update info
	if update[0] == sl.head {
		node.prev = nil
	} else {
		node.prev = update[0]
	}
	if node.level[0].next != nil {
		node.level[0].next.prev = node
	} else {
		sl.tail = node
	}
	sl.size++
	return node.value
}

// SkipList is a data structure that allows fast search, insertion, and removal operations.
// Here we define a method delete on it.
//
// The delete method in the skip list will remove nodes that have a given score and key from the skip list.
// If no such nodes are found, the function does nothing.
//
// Parameters:
//
//	score: the score of the node to delete.
//	key: the key of the node to delete.
func (sl *SkipList) delete(score int, member string) {

	// update: an array of pointers to SkipListNodes; holds the nodes that will have their next pointers updated.
	update := make([]*SkipListNode, SKIPLIST_MAX_LEVEL)

	// node: start from the head of our SkipList sl
	node := sl.head

	// The code block of "for" loop populates the "update" variable with nodes which reference will change
	// due to the removal of the target node.
	for i := sl.level; i >= 0; i-- {
		// This loop is traversing the SkipList horizontally until it finds a node with a score greater
		// than or equal to our target score or if the scores are equal it also checks the member.
		for node.level[i].next != nil &&
			(node.level[i].next.value.score < score ||
				(node.level[i].next.value.score == score &&
					node.level[i].next.value.member < member)) {
			node = node.level[i].next
		}
		update[i] = node
	}

	// After the traversal, we set the node to point to the possibly (to be) deleted node.
	node = node.level[0].next

	// If the possibly deleted node is the target node (it has the same score and member), then remove it.
	if node != nil && node.value.score == score && node.value.member == member {
		sl.deleteNode(node, update)
	}
}
func (sl *SkipList) getRange(start int, end int, reverse bool) (nv []ZSetValue) {
	if end > sl.size {
		end = sl.size - 1
	}
	if start > end {
		return
	}
	if end <= 0 {
		return nil // todo unexpected behavior, we can set it to zero as well
	}

	node := sl.head
	if reverse {
		node = sl.tail
		if start > 0 {
			node = sl.getNodeByRank(sl.size - start)
		}
	} else {
		node = node.level[0].next
		if start > 0 {
			node = sl.getNodeByRank(start + 1)
		}
	}

	for i := start; i < end; i++ {
		if reverse {
			nv = append(nv, *node.value)
			node = node.prev
		} else {
			nv = append(nv, *node.value)
			node = node.level[0].next
		}
	}
	return nv
}

// deleteNode is a method linked to the SkipList struct that allows to remove nodes from the SkipList instance.
// It takes two parameters: a pointer to the node to be deleted, and a slice of pointers to SkipListNode which are required for node updates.
// deleteNode performs the deletion through a two-step process:
// - First, it loops over every level in the SkipList, updating level spans and next node pointers accordingly.
// - Then, it sets the pointers back to the previous node in the data structure and updates the tail and level of the whole list.
// Finally, it decreases the size of the list by one, as a node is being removed from it.
// It doesn't return any value and modifies the SkipList directly.
func (sl *SkipList) deleteNode(node *SkipListNode, updates []*SkipListNode) {
	for i := 0; i < sl.level; i++ {
		if updates[i].level[i].next == node {
			updates[i].level[i].span += node.level[i].span - 1
			updates[i].level[i].next = node.level[i].next
		} else {
			updates[i].level[i].span--
		}
	}
	//update backwards
	if node.level[0].next != nil {
		node.level[0].next.prev = node.prev
	} else {
		sl.tail = node.prev
	}

	for sl.level > 1 && sl.head.level[sl.level-1].next == nil {
		sl.level--
	}
	sl.size--
}

// getRank method receives a SkipList pointer and two parameters: an integer 'score' and a string 'key'.
// It then calculates the rank of an element in the SkipList. The rank is determined based on two conditions:
// - the score of the next node is less than the provided score
// - or, the score of the next node equal to the provided score and the key of the next node is less than or equal to the provided key.
//
// Parameters:
// sl: A pointer to the SkipList object.
// score: The score that we are comparing with the scores in the skiplist.
// key: The key that we are comparing with the keys in the skiplist.
//
// Return:
// Returns the rank of the element in the SkipList if it's found, otherwise returns 0.
func (sl *SkipList) getRank(score int, key string) int {
	var rank int
	h := sl.head // Start at the head node of the SkipList

	// For loop starts from the top level and goes down to the level 0
	for i := sl.level; i >= 0; i-- {
		// While loop advances the 'h' pointer as long as the next node exists and the conditions are fulfilled
		for h.level[i].next != nil &&
			(h.level[i].next.value.score < score ||
				(h.level[i].next.value.score == score &&
					h.level[i].next.value.member <= key)) {

			// Increase the rank by the span of the current level
			rank += h.level[i].span
			// Move to the next node
			h = h.level[i].next
		}
		// If the key of the current node is equal to the provided key, return the rank
		if h.value.member == key {
			return rank
		}
	}
	// If the element is not found in the SkipList, return 0
	return 0
}

// getNodeByRank is a method of the SkipList type that is used to retrieve a node based on its rank within the list.
// The method takes as argument an integer rank and returns a pointer to the SkipListNode at the specified rank,
// or nil if there is no such node.
//
// First, the method initializes a variable traversed to store the cumulative span of nodes traversed thus far in the search.
// It sets a helper variable h to the head of the SkipList, to begin the traversal.
//
// The method then enter a loop that iterates through the levels of the SkipList from the highest down to the base level.
// On each level, while the next node exists and the total span traversed plus the span of the next node doesn't exceed the target rank,
// the method moves to the next node and adds its span to the cumulative span traversed.
//
// If during the traversal the cumulative span equals the target rank, the method returns the getNodeByRank node.
// If the end of the SkipList is reached, or the target rank isn't found on any level, the method returns nil.
func (sl *SkipList) getNodeByRank(rank int) *SkipListNode {
	// This variable is used to keep track of the number of nodes we have
	// traversed while going through the levels of the SkipList.
	var traversed int

	// Define a SkipListNode pointer h, initialized with sl.head
	// At the start, this pointer is set to head node of the SkipList.
	h := sl.head

	// The outer loop decrements levels from highest level to lowest.
	for i := sl.level - 1; i >= 0; i-- {

		// The inner loop traverses the nodes at current level while the next node isn't null and we haven't traversed beyond the 'rank'.
		// The traversed variable is also updated to include the span of the current level.
		for h.level[i].next != nil && (traversed+h.level[i].span) <= rank {
			traversed += h.level[i].span
			h = h.level[i].next
		}

		// If traversed equals 'rank', it means we've found the node at the rank we are looking for.
		// So, return the node.
		if traversed == rank {
			return h
		}
	}

	// If the node at 'rank' wasn't found in the SkipList, return nil.
	return nil
}

// ZAdd adds a value with its given score and member to a sorted set (ZSet), associated with
// the provided key. It is a method on the ZSetStructure type.
//
// Parameters:
//
//	key:    a string that represents the key of the sorted set.
//	score:  an integer value that determines the order of the added element in the sorted set.
//	member: a string used for identifying the added value within the sorted set.
//	value:  the actual value to be stored within the sorted set.
//
// If the key is an empty string, an error will be returned
func (zs *ZSetStructure) ZAdd(key string, score int, member string, value string) error {
	return zs.ZAdds(key, []ZSetValue{{score: score, member: member, value: value}}...)
}

// ZAdds adds a value with its given score and member to a sorted set (ZSet), associated with
// the provided key. It is a method on the ZSetStructure type.
//
// Parameters:
//
//	values:    ...ZSetValue multiple values of ZSetValue.
func (zs *ZSetStructure) ZAdds(key string, vals ...ZSetValue) error {
	if err := checkKey(key); err != nil {
		return err
	}
	zSet, err := zs.getOrCreateZSet(key)
	if err != nil {
		return fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}

	for _, val := range vals {
		// if values didn't change, do nothing
		if zs.valuesDidntChange(zSet, val.score, val.member, val.value) {
			continue
		}
		if err := zs.updateZSet(zSet, key, val.score, val.member, val.value); err != nil {
			return fmt.Errorf("failed to set ZSet to DB with key '%v': %w", key, err)
		}
	}

	return zs.setZSetToDB(stringToBytesWithKey(key), zSet)
}

// Keys returns all keys of the ZSetStructure.
//
// Returns:
//
//	[]string: all keys of the ZSetStructure.
func (zs *ZSetStructure) Keys() ([]string, error) {
	var keys []string
	byte_keys := zs.db.GetListKeys()
	for _, key := range byte_keys {
		keys = append(keys, string(key))
	}
	return keys, nil
}

// exists checks if a given member with a specific score exists in a ZSet. It
// also verifies if the provided key is valid. The function returns a boolean
// value indicating whether the member with the specified score exists in the
// ZSet or not.
//
// Parameters:
//
//	key (string): Specifies the key of the ZSet.
//	score (int): The score of the member to be checked.
//	member (string): The specific member to check for in the ZSet.
//
// Returns:
//
//	bool: A boolean value indicating whether a member with the specified score
//
// exists in the ZSet or not. Returns false if the ZSet does not exist or if
// the key is invalid.
func (zs *ZSetStructure) exists(key string, score int, member string) bool {
	if err := checkKey(key); err != nil {
		return false
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)

	if err != nil {
		return false
	}
	return zSet.exists(score, member)
}

/*
ZRem is a method belonging to ZSetStructure that removes a member from a ZSet.

Parameters:
  - key (string): The key of the ZSet.
  - member (string): The member to be removed.

Returns:
  - error: An error if the operation fails.

The ZRem method checks for a non-empty key, retrieves the corresponding ZSet
from the database, removes the specified member, and then updates
the ZSet in the database. If any point of this operation fails,
the function will return the corresponding error.
*/
func (zs *ZSetStructure) ZRem(key string, member string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return err
	}
	if err = zSet.RemoveNode(member); err != nil {
		return err
	}
	return zs.setZSetToDB(keyBytes, zSet)
}

// ZRems method removes one or more specified members from the sorted set that's stored under the provided key.
// Params:
//   - key string: the identifier for storing the sorted set in the database.
//   - member ...string: a variadic parameter where each argument is a member string to remove.
//
// Returns: error
//
// The function will return an error if it fails at any point, if not it will return nil indicating a successful operation.
func (zs *ZSetStructure) ZRems(key string, member ...string) error {
	if err := checkKey(key); err != nil {
		return err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)

	if err != nil {
		return fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	for _, s := range member {
		if err = zSet.RemoveNode(s); err != nil {
			return err
		}
	}
	return zs.setZSetToDB(keyBytes, zSet)
}

// ZScore method retrieves the score associated with the member in a sorted set stored at the key
func (zs *ZSetStructure) ZScore(key string, member string) (int, error) {
	if err := checkKey(key); err != nil {
		return 0, err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, err
	}
	// if the member in the sorted set is found, return the score associated with it
	if v, ok := zSet.dict[member]; ok {
		return v.score, nil
	}

	// if the member doesn't exist in the set, return score of zero and an error
	return 0, _const.ErrKeyNotFound
}

/*
ZRank is a method belonging to the ZSetStructure type. This method retrieves the rank of an element within a sorted set identified by a key. The rank is an integer corresponding to the element's 0-based position in the sorted set when it is arranged in ascending order.

Parameters:
key (string): The key that identifies the sorted set.
member (string): The element for which you want to find the rank.

Returns:
int: An integer indicating the rank of the member in the set.

	Rank zero means the member is not found in the set.

error: If an error occurs, an error object will be returned.

	Possible errors include:
	- key is empty
	- failure to get or create the ZSet from the DB
	- the provided key does not exist in the DB

Example:
rank, err := zs.ZRank("myKey", "memberName")

	if err != nil {
	   log.Fatal(err)
	}

fmt.Printf("The rank of '%s' in the set '%s' is %d\n", "memberName", "myKey", rank)
*/
func (zs *ZSetStructure) ZRank(key string, member string) (int, error) {
	if err := checkKey(key); err != nil {
		return 0, err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	if v, ok := zSet.dict[member]; ok {
		return zSet.skipList.getRank(v.score, member), nil
	}

	// rank zero means no rank found
	return 0, _const.ErrKeyNotFound
}

// ZRevRank calculates the reverse rank of a member in a ZSet (Sorted Set) associated with a given key.
// ZSet exploits the Sorted Set data structure of Redis with O(log(N)) time complexity for Fetching the rank.
//
// Parameters:
//
//	key:    This is a string that serves as the key of a ZSet stored in the database.
//	member: This is a string that represents a member of a ZSet whose rank needs to be obtained.
//
// Returns:
//
//	int:    The integer represents the reverse rank of the member in the ZSet. It returns 0 if the member is not found in the ZSet.
//	        On successful execution, it returns the difference of the ZSet size and the member's rank.
//	error:  The error which will be null if no errors occurred. If the key provided is empty, an ErrKeyIsEmpty error is returned.
//	        If there's a problem getting or creating the ZSet from the database, an error message is returned with the format
//	        "failed to get or create ZSet from DB with key '%v': %w", where '%v' is the key and '%w' shows the error detail.
//	        If the member is not found in the ZSet, it returns an ErrKeyNotFound error.
//
// Note: The reverse rank is calculated as 'size - rank', and the ranks start from 1.
func (zs *ZSetStructure) ZRevRank(key string, member string) (int, error) {
	if err := checkKey(key); err != nil {
		return 0, err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	if v, ok := zSet.dict[member]; ok {
		rank := zSet.skipList.getRank(v.score, member)
		return (zSet.size) - rank + 1, nil
	}

	// rank zero means no rank found
	return 0, _const.ErrKeyNotFound
}

// ZRange retrieves a specific range of elements from a sorted set (ZSet) denoted by a specific key.
// It returns a slice of ZSetValue containing the elements within the specified range (inclusive), and a nil error when successful.
//
// The order of the returned elements is based on their rank in the set, not their score.
//
// Parameters:
//
//	key: A string identifier representing the ZSet. The key shouldn't be an empty string.
//	start: A zero-based integer representing the first index of the range.
//	end: A zero-based integer representing the last index of the range.
//
// Returns:
//
//	 []ZSetValue:
//			Slice of ZSetValue containing elements within the specified range.
//	 error:
//			An error if it occurs during execution, such as:
//	     		1. The provided key string is empty.
//	     		2. An error occurs while fetching the ZSet from the database, i.e., the ZSet represented by the given key doesn't exist.
//	     		In the case of an error, an empty slice and the actual error encountered will be returned.
//
// Note:
// On successful execution, ZRange returns the elements starting from 'start' index up to 'end' index inclusive.
// If the set doesn't exist or an error occurs during execution, ZRange returns an empty slice and the error.
//
// Example:
// Assume we have ZSet with the following elements: ["element1", "element2", "element3", "element4"]
// ZRange("someKey", 0, 2) will return ["element1", "element2", "element3"] and nil error.
//
// This method is part of the ZSetStructure type.
func (zs *ZSetStructure) ZRange(key string, start int, end int) ([]ZSetValue, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return nil, err
	}
	r := zSet.skipList.getRange(start, end, false)

	// rank zero means no rank found
	return r, nil
}

// ZCount traverses through the elements of the ZSetStructure based on the given key.
// The count of elements between the range of min and max scores is determined.
//
// The method takes a string as the key and two integers as min and max ranges.
// The range values are inclusive: [min, max]. If min is greater than max, an error is returned.
// The function ignores scores that fall out of the specified min and max range.
//
// It returns the count of elements within the range, and an error if any occurs during the process.
//
// For example, use as follows:
// count, err := zs.ZCount("exampleKey", 10, 50)
// This will count the number of elements that have the scores between 10 and 50 in the ZSetStructure associated with "exampleKey".
//
// Returns:
//  1. int: The total count of elements based on the score range.
//  2. error: Errors that occurred during execution, if any.
func (zs *ZSetStructure) ZCount(key string, min int, max int) (count int, err error) {
	if err = checkKey(key); err != nil {
		return 0, err
	}
	keyBytes := stringToBytesWithKey(key)
	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, err
	}
	if min > max {
		return 0, ErrInvalidArgs
	}
	min, max, err = zs.adjustMinMax(zSet, min, max)
	if err != nil {
		return 0, err
	}
	x := zSet.skipList.head
	// Node traversal loop. We keep moving to the next node at current level
	// as long as the score of the next node's value is less than 'min'.
	for i := zSet.skipList.level - 1; i >= 0; i-- {
		for x.level[i].next != nil && x.level[i].next.value.score < min {
			x = x.level[i].next
		}
	}

	x = x.level[0].next
	// Score range check loop. We traverse nodes and increment 'count'
	// as long as node value's score is in the range ['min', 'max']
	for x != nil {
		if x.value.score > max {
			break
		}
		count++
		x = x.level[0].next
	}
	return count, nil
}

// ZRevRange retrieves a range of elements from a sorted set (ZSet) in descending order.
// Inputs:
//   - key: Name of the ZSet
//   - startRank: Initial rank of the desired range
//   - endRank: Final rank of the desired range
//
// Output:
//   - An array of ZSetValue, representing elements from the range [startRank, endRank] in descending order
//   - Error if an issue occurs, such as when the key is empty or ZSet retrieval fails
//     error
func (zs *ZSetStructure) ZRevRange(key string, startRank int, endRank int) ([]ZSetValue, error) {
	if err := checkKey(key); err != nil {
		return nil, err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return nil, err
	}
	r := zSet.skipList.getRange(startRank, endRank, true)

	// rank zero means no rank found
	return r, nil
}

// The ZCard function returns the size of the dictionary of the sorted set stored at key in the database.
// It takes a string key as an argument.
func (zs *ZSetStructure) ZCard(key string) (int, error) {
	if err := checkKey(key); err != nil {
		return 0, err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, err
	}
	// get the size of the dictionary
	return zSet.size, nil
}

// ZIncrBy increases the score of an existing member in a sorted set stored at specified key by
// the increment `incBy` provided. If member does not exist, ErrKeyNotFound error is returned.
// If the key does not exist, it treats it as an empty sorted set and returns an error.
//
// The method accepts three parameters:
// `key`: a string type parameter that identifies the sorted set
// `member`: a string type parameter representing member in the sorted set
// `incBy`: an int type parameter provides the increment value for a member score
//
// The method throws error under following circumstances -
// if provided key is empty (ErrKeyIsEmpty error),
// if provided key or member is not present in the database (ErrKeyNotFound error),
// if it's unable to fetch or create ZSet from DB,
// if there's an issue with node insertion,
// if unable to set ZSet to DB post increment operation
func (zs *ZSetStructure) ZIncrBy(key string, member string, incBy int) error {
	if err := checkKey(key); err != nil {
		return err
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}

	if v, ok := zSet.dict[member]; ok {
		if err = zSet.InsertNode(v.score+incBy, member, v.value); err != nil {
			return err
		}
		if err = zs.setZSetToDB(keyBytes, zSet); err != nil {
			return err
		}
		return zs.setZSetToDB(keyBytes, zSet)
	}

	return _const.ErrKeyNotFound
}

// getOrCreateZSet attempts to retrieve a sorted set by a key, or creates a new one if it doesn't exist.
func (zs *ZSetStructure) getOrCreateZSet(key string) (*FZSet, error) {
	keyBytes := stringToBytesWithKey(key)
	zSet, err := zs.getZSetFromDB(keyBytes)
	// if key is not in the DB, create it.
	if errors.Is(err, _const.ErrKeyNotFound) {
		return newZSetNodes(), nil
	}

	return zSet, err
}

// valuesDidntChange checks if the data of a specific member in a sorted set remained the same.
func (zs *ZSetStructure) valuesDidntChange(zSet *FZSet, score int, member string, value interface{}) bool {
	if v, ok := zSet.dict[member]; ok {
		return v.score == score && v.member == member && v.value == value
	}

	return false
}

// updateZSet updates or inserts a member in a sorted set and saves the change in storage.
func (zs *ZSetStructure) updateZSet(zSet *FZSet, key string, score int, member string, value interface{}) error {
	return zSet.InsertNode(score, member, value)
}

// InsertNode is a method on the FZSet structure. It inserts a new node
// or updates an existing node in the skip list and the dictionary.
// It takes three parameters: score (an integer), key (a string),
// and value (of any interface type).
//
// If key already exists in the dictionary and the score equals the existing
// score, it updates the value and score in the skip list and the dictionary.
// If the score is different, it only updates the value in the dictionary
// because the ranking doesn't change and there is no need for an update in the
// skip list.
//
// If the key doesn't exist in the dictionary, it adds the new key, value and score
// to the dictionary, increments the size of the dictionary by 1, and also adds
// the node to the skip list.
func (fzs *FZSet) InsertNode(score int, member string, value interface{}) error {
	// Instantiate dictionary if it's not already
	if fzs.dict == nil {
		fzs.dict = make(map[string]*ZSetValue)
	}
	if fzs.skipList == nil {
		fzs.skipList = newSkipList()
	}

	// Check if key exists in dictionary
	if v, ok := fzs.dict[member]; ok {
		if v.score != score {
			// Update value and score as the score remains the same
			fzs.skipList.delete(score, member)
			fzs.dict[member] = fzs.skipList.insert(score, member, value)
		} else {
			// Ranking isn't altered, only update value
			v.value = value
		}
	} else { // Key doesn't exist, create new key
		fzs.dict[member] = fzs.skipList.insert(score, member, value)
		fzs.size++ // Increase size count by 1
		// Node is also added to the skip list
	}

	// Returns nil as no specific error condition is checked in this function
	return nil
}
func (zs *ZSetStructure) adjustMinMax(zSet *FZSet, min int, max int) (adjustedMin int, adjustedMax int, err error) {
	if min > max {
		return min, max, ErrInvalidArgs
	}
	minScore, maxScore := zSet.getMinMaxScore()
	return zSet.max(min, minScore), zSet.min(max, maxScore), nil
}
func (fzs *FZSet) getMinMaxScore() (minScore int, maxScore int) {
	if fzs == nil || fzs.skipList == nil || fzs.skipList.head == nil || len(fzs.skipList.head.level) < 1 || fzs.skipList.head.level[0].next == nil || fzs.skipList.tail == nil {
		return 0, 0
	}

	if fzs.skipList.head.level[0].next.value == nil || fzs.skipList.tail.value == nil {
		return 0, 0
	}
	return fzs.skipList.head.level[0].next.value.score,
		fzs.skipList.tail.value.score
}
func (fzs *FZSet) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (fzs *FZSet) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// RemoveNode is a method for FZSet structure.
// This method aims to delete a node from both
// the dictionary (dict) and the skip list (skipList).
//
// The method receives one parameter:
//   - member: a string that represents the key of the node
//     to be removed from the FZSet structure.
//
// The method follows these steps:
//  1. Check if a node with key 'member' exists in the dictionary.
//     If not, or if the dictionary itself is nil, it returns an error
//     (_const.ErrKeyNotFound) indicating that the node cannot be found.
//  2. If the node exists, it proceeds to remove the node from both the
//     skip list and dictionary.
//  3. After the successful removal of the node, it returns nil indicating
//     the success of the operation.
//
// The RemoveNode's primary purpose is to provide a way to securely and
// efficiently remove a node from the FZSet structure.
func (fzs *FZSet) RemoveNode(member string) error {
	// Check for existence of key in dictionary
	v, ok := fzs.dict[member]
	if !ok || fzs.dict == nil {
		return _const.ErrKeyNotFound
	}

	// Delete Node from the skip list and dictionary
	fzs.skipList.delete(v.score, member)
	delete(fzs.dict, member)
	fzs.size--

	return nil
}

func (fzs *FZSet) exists(score int, member string) bool {
	v, ok := fzs.dict[member]

	return ok && v.score == score
}

// Bytes encodes the FZSet instance into bytes using MessagePack
// binary serialization format. The encoded bytes can be used for
// storage or transmission. If the encoding operation fails, an
// error is returned.
func (fzs *FZSet) Bytes() ([]byte, error) {
	var msgPack = encoding.NewMessagePackEncoder()
	if encodingError := msgPack.Encode(fzs); encodingError != nil {
		return nil, encodingError
	}
	return msgPack.Bytes(), nil
}

// FromBytes decodes the input byte slice into the FZSet object using MessagePack.
// Returns an error if decoding fails, otherwise nil.
func (fzs *FZSet) FromBytes(b []byte) error {
	return encoding.NewMessagePackDecoder(b).Decode(fzs)
}

// getZSetFromDB fetches and deserializes FZSet from the database.
//
// Returns a pointer to the FZSet and error, if any.
// If the key doesn't exist, both the pointer and the error will be nil.
// In case of deserialization errors, returns nil and the error.
func (zs *ZSetStructure) getZSetFromDB(key []byte) (*FZSet, error) {
	dbData, err := zs.db.Get(key)

	// If key is not found, return nil for both; otherwise return the error.
	if err != nil {

		return nil, err
	}
	dec := encoding.NewMessagePackDecoder(dbData)
	// Deserialize the data.
	var zSetValue FZSet
	if err = dec.Decode(&zSetValue); err != nil {
		return nil, err
	}

	// return a pointer to the deserialized FZSet, nil for the error
	return &zSetValue, nil
}

// checkKey function that accepts a string parameter key
// and returns error if key is empty.
//
// # It returns nil otherwise
//
// Parameters:
//
//	key : A string that is checked if empty
//
// Returns:
//
//	error : _const.ErrKeyIsEmpty if key is empty, nil otherwise
func checkKey(key string) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	return nil
}

// setZSetToDB writes a FZSet object to the database.
//
//	parameters:
//	key: This is a byte slice that is used as a key in the database.
//	zSetValue: This is a pointer to a FZSet object that needs to be stored in the database.
//
// The function serializes the FZSet object into MessagePack format. If an error occurs
// either during serialization or when writing to the database, that specific error is returned.
// If the process is successful, it returns nil.
func (zs *ZSetStructure) setZSetToDB(key []byte, zSetValue *FZSet) error {
	val := encoding.NewMessagePackEncoder()
	err := val.Encode(zSetValue)
	if err != nil {
		return err
	}
	return zs.db.Put(key, val.Bytes())
}

// UnmarshalBinary de-serializes the given byte slice into FZSet instance
// it uses MessagePack format for de-serialization
// Returns an error if the decoding of size or insertion of node fails.
//
// Parameters:
// data : a slice of bytes to be decoded
//
// Returns:
// An error that will be nil if the function succeeds.
func (fzs *FZSet) UnmarshalBinary(data []byte) (err error) {
	// NewMessagePackDecoder creates a new MessagePack decoder with the provided data
	dec := encoding.NewMessagePackDecoder(data)

	var size int
	// Decode the size of the data structure
	if err = dec.Decode(&size); err != nil {
		return err // error handling if something goes wrong with decoding
	}

	// Iterate through each node in the data structure
	for i := 0; i < size; i++ {
		// Create an empty instance of ZSetValue for each node
		slValue := ZSetValue{}

		// Decode each node onto the empty ZSetValue instance
		if err = dec.Decode(&slValue); err != nil {
			return err // error handling if something goes wrong with decoding
		}

		// Insert the decoded node into the FZSet instance
		if err = fzs.InsertNode(slValue.score, slValue.member, slValue.value); err != nil {
			return err
		}
	}
	return // if all nodes are correctly decoded and inserted, return with nil error
}

// MarshalBinary serializes the FZSet instance into a byte slice.
// It uses MessagePack format for serialization
// Returns the serialized byte slice and an error if the encoding fails.
func (fzs *FZSet) MarshalBinary() (_ []byte, err error) {

	// Initializing the MessagePackEncoder
	enc := encoding.NewMessagePackEncoder()

	// Encoding the size attribute of d (i.e., d.size). The operation could fail, thus we check for an error.
	// An error, if occurred, will be returned immediately, hence the flow of execution stops here.
	err = enc.Encode(fzs.size)
	if err != nil {
		return nil, err
	}

	// This is the start of a loop going over all the nodes in d's skip list from the tail of the
	// list to the head.
	// The tail and head pointers refer to the last and first element of the list, respectively,
	// and are maintained for efficient traversing of the list.
	// we do that to get the elements in reverse order from biggest to the smallest for the best
	// insertion efficiency as it makes the insertion O(1), because each new element to be inserted is
	// the smallest yet.
	x := fzs.skipList.tail
	// as long as there are elements in the SkipList continue
	for x != nil {
		// Encoding the value of the current node in the skip list
		// Again, if an error occurs it gets immediately returned, thus breaking the loop.
		err = enc.Encode(x.value)
		if err != nil {
			return nil, err
		}

		// Move to the previous node in the skip list.
		x = x.prev
	}

	// After the traversal of the skip list, the encoder should now hold the serialized representation of the
	// FZSet. Now, we return the bytes from the encoder along with any error that might have occurred
	// during the encoding (should be nil if everything went fine).
	return enc.Bytes(), err
}

// UnmarshalBinary de-serializes the given byte slice into ZSetValue instance
// It uses the MessagePack format for de-serialization
// Returns an error if the decoding of Key, Score, or Value fails.
func (p *ZSetValue) UnmarshalBinary(data []byte) (err error) {
	dec := encoding.NewMessagePackDecoder(data)
	if err = dec.Decode(&p.member); err != nil {
		return
	}
	if err = dec.Decode(&p.score); err != nil {
		return err
	}
	if err = dec.Decode(&p.value); err != nil {
		return
	}
	return
}

// MarshalBinary uses MessagePack as the encoding format to serialize
// the ZSetValue object into a byte array.
func (d *ZSetValue) MarshalBinary() (_ []byte, err error) {

	// The NewMessagePackEncoder function is called to create a new
	// MessagePack encoder.
	enc := encoding.NewMessagePackEncoder()

	// Then, we try to encode the 'key' field of the ZSetValue
	// If an error occurs, it is returned immediately along with the
	// currently encoded byte slice.
	if err = enc.Encode(d.member); err != nil {
		return enc.Bytes(), err
	}

	// We do the same for the 'score' field.
	if err = enc.Encode(d.score); err != nil {
		return enc.Bytes(), err
	}

	// Lastly, the 'value' field is encoded in the same way.
	if err = enc.Encode(d.value); err != nil {
		return enc.Bytes(), err
	}

	// If everything goes well and we're done encoding, we return the
	// final byte slice which represents the encoded ZSetValue
	// and a nil error.
	return enc.Bytes(), err
}

func (d *ZSetStructure) Stop() error {
	err := d.db.Close()
	return err
}
