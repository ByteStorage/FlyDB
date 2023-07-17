package structure

import (
	"errors"
	"fmt"
	"github.com/ByteStorage/FlyDB/config"
	"math"
	"math/rand"

	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/ByteStorage/FlyDB/lib/encoding"
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

// ZSetNodes represents a specific data structure in the database, which is key to handling sorted sets (ZSets).
// This struct facilitates interactions with data stored in the sorted set, allowing for both complex and simple operations.
//
// It contains three struct fields:
//
//   - 'dict': A Go map with string keys and pointers to SkipListNodeValue values. This map aims to provide quick access to
//     individual values in the sorted set based on the provided key.
//
//   - 'size': An integer value representing the current size (number of elements) in the ZSetNodes struct. This information is efficiently
//     kept track of whenever elements are added or removed from the set, so no separate computation is needed to retrieve this information.
//
//   - 'skipList': A pointer towards a SkipList struct. SkipLists perform well under numerous operations, such as insertion, deletion, and searching. They are
//     a crucial component in maintaining the sorted set in a practical manner. In this context, the SkipList is used to keep an ordered track of the elements
//     in the ZSetNodes struct.
type ZSetNodes struct {
	// dict field is a map where the key is a string and
	// the value is a pointer to SkipListNodeValue instances,
	// codified with the tag "dict".
	dict map[string]*SkipListNodeValue `codec:"dict"`

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
//   - 'value': This is a pointer towards a single SkipListNodeValue structure. It holds the actual payload of the node
//     (namely the 'score', 'key', and 'value' properties used in the context of Redis Sorted Sets), as well as provides the basis for ordering of nodes in the skip list.
type SkipListNode struct {
	// prev is a pointer to the previous node in the skip list.
	prev *SkipListNode

	// level is a slice of pointers to SkipListLevel.
	// Each level represents a forward pointer to the next node in the current list level.
	level []*SkipListLevel

	// value is a pointer to the SkipListNodeValue.
	// This represents the value that this node holds.
	value *SkipListNodeValue
}

// SkipListNodeValue is a struct used in the SkipList data structure. In the context of Redis Sorted Set (ZSet) implementation,
// it represents a single node value in the skip list. A SkipListNodeValue has three members:
// - 'score' which is an integer representing the score of the node. Nodes in a skip list are ordered by this score in ascending order.
// - 'member' which is a string defining the key of the node. For nodes with equal scores, order is determined with lexicographical comparison of keys.
// - 'value' which is an interface{}, meaning it can hold any data type. This represents the actual value of the node in the skip list.
type SkipListNodeValue struct {
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

	// Calculate the threshold for level. It's derived from the probability constant of the skip list.
	thresh := int(math.Round(SKIPLIST_PROB * 0xFFF))

	// While a randomly generated number is less than this threshold, increment the level.
	for int(rand.Int31()&0xFFF) < thresh {
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

// newZSetNodes is a function that creates a new ZSetNodes object and returns a pointer to it.
// It initializes the dictionary member dict of the newly created object to an empty map.
// The map is intended to map strings to pointers of SkipListNodeValue objects.
// size member of the object is set to 0, indicating that the ZSetNodes object is currently empty.
// The skipList member of the object is set to a new SkipList object created by calling `newSkipList()` function.
func newZSetNodes() *ZSetNodes {
	return &ZSetNodes{
		dict:     make(map[string]*SkipListNodeValue),
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

// newSkipListNodeValue is a function that constructs and returns a new SkipListNodeValue.
// It takes a score (int), a key (string), and a value (interface{}) as parameters.
// These parameters serve as the initial state of the SkipListNodeValue upon its creation.
func newSkipListNodeValue(score int, member string, value interface{}) *SkipListNodeValue {
	// Create a new instance of a SkipListNodeValue with the provided score, key, and value.
	node := &SkipListNodeValue{
		score:  score,
		member: member,
		value:  value,
	}

	// Return the newly created SkipListNodeValue.
	return node
}

// insert is a method of the SkipList type that is used to insert a new node into the skip list. It takes as arguments
// the score (int), key (string) and a value (interface{}), and returns a pointer to the SkipListNodeValue struct. The method
// organizes nodes in the list based on the score in ascending order. If two nodes have the same score, they will be arranged
// based on the key value. The method also assigns span values to the levels in the skip list.
func (sl *SkipList) insert(score int, key string, value interface{}) *SkipListNodeValue {
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
func (sl *SkipList) getRange(start int, end int, reverse bool) (nv []SkipListNodeValue) {
	if end > sl.size {
		end = sl.size - 1
	}
	if start > end {
		return
	}
	if end < 0 {
		return nil // todo unexpected behavior, we can set it to zero as well
	}
	node := sl.head
	if reverse {
		node = sl.getNodeByRank(end)
	} else {
		node = sl.getNodeByRank(start)
	}
	if reverse {
		node = sl.getNodeByRank(end)
	} else {
		node = sl.getNodeByRank(start)
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
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}

	zSet, err := zs.getOrCreateZSet(key)

	if err != nil {
		return fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}

	// if values didn't change, do nothing
	if zs.valuesDidntChange(zSet, score, member, value) {
		return nil
	}

	if err := zs.updateZSet(zSet, key, score, member, value); err != nil {
		return fmt.Errorf("failed to set ZSet to DB with key '%v': %w", key, err)
	}

	return nil
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
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)

	if err != nil {
		return fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	if err = zSet.RemoveNode(member); err != nil {
		return err
	}
	return zs.setZSetToDB(keyBytes, zSet)
}

// ZScore method retrieves the score associated with the member in a sorted set stored at the key
func (zs *ZSetStructure) ZScore(key string, member string) (int, error) {
	if len(key) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
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
	if len(key) == 0 {
		return 0, _const.ErrKeyIsEmpty
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
func (zs *ZSetStructure) ZRevRank(key string, member string) (int, error) {
	if len(key) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	if v, ok := zSet.dict[member]; ok {
		rank := zSet.skipList.getRank(v.score, member)
		return zSet.size - rank, nil
	}

	// rank zero means no rank found
	return 0, _const.ErrKeyNotFound
}
func (zs *ZSetStructure) ZRange(key string, start int, end int) ([]SkipListNodeValue, error) {
	if len(key) == 0 {
		return nil, _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	r := zSet.skipList.getRange(start, end, false)

	// rank zero means no rank found
	return r, nil
}
func (zs *ZSetStructure) ZRevRange(key string, start int, end int) ([]SkipListNodeValue, error) {
	if len(key) == 0 {
		return nil, _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	r := zSet.skipList.getRange(start, end, true)

	// rank zero means no rank found
	return r, nil
}
func (zs *ZSetStructure) ZCard(key string) (int, error) {
	if len(key) == 0 {
		return 0, _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	// get the size of the dictionary
	return zSet.size, nil
}
func (zs *ZSetStructure) ZIncrBy(key string, member string, incBy int) error {
	if len(key) == 0 {
		return _const.ErrKeyIsEmpty
	}
	keyBytes := stringToBytesWithKey(key)

	zSet, err := zs.getZSetFromDB(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to get or create ZSet from DB with key '%v': %w", key, err)
	}
	if v, ok := zSet.dict[member]; ok {
		return zSet.InsertNode(v.score+incBy, member, v.value)
	}

	return _const.ErrKeyNotFound
}

// getOrCreateZSet attempts to retrieve a sorted set by a key, or creates a new one if it doesn't exist.
func (zs *ZSetStructure) getOrCreateZSet(key string) (*ZSetNodes, error) {
	keyBytes := stringToBytesWithKey(key)
	zSet, err := zs.getZSetFromDB(keyBytes)
	// if key is not in the DB, create it.
	if errors.Is(err, _const.ErrKeyNotFound) {
		return newZSetNodes(), nil
	}

	return zSet, err
}

// valuesDidntChange checks if the data of a specific member in a sorted set remained the same.
func (zs *ZSetStructure) valuesDidntChange(zSet *ZSetNodes, score int, member string, value string) bool {
	if v, ok := zSet.dict[member]; ok {
		return v.score == score && v.member == member && v.value == value
	}

	return false
}

// updateZSet updates or inserts a member in a sorted set and saves the change in storage.
func (zs *ZSetStructure) updateZSet(zSet *ZSetNodes, key string, score int, member string, value string) error {
	if err := zSet.InsertNode(score, member, value); err != nil {
		return err
	}

	return zs.setZSetToDB(stringToBytesWithKey(key), zSet)
}

// InsertNode is a method on the ZSetNodes structure. It inserts a new node
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
func (pq *ZSetNodes) InsertNode(score int, member string, value interface{}) error {
	// Instantiate dictionary if it's not already
	if pq.dict == nil {
		pq.dict = make(map[string]*SkipListNodeValue)
	}

	// Check if key exists in dictionary
	if v, ok := pq.dict[member]; ok {
		if v.score == score {
			// Update value and score as the score remains the same
			pq.skipList.delete(score, member)
			pq.dict[member] = pq.skipList.insert(score, member, value)
		} else {
			// Ranking isn't altered, only update value
			v.value = value
		}
	} else { // Key doesn't exist, create new key
		pq.dict[member] = pq.skipList.insert(score, member, value)
		pq.size++ // Increase size count by 1
		// Node is also added to the skip list
	}

	// Returns nil as no specific error condition is checked in this function
	return nil
}

// RemoveNode is a method for ZSetNodes structure.
// This method aims to delete a node from both
// the dictionary (dict) and the skip list (skipList).
//
// The method receives one parameter:
//   - member: a string that represents the key of the node
//     to be removed from the ZSetNodes structure.
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
// efficiently remove a node from the ZSetNodes structure.
func (pq *ZSetNodes) RemoveNode(member string) error {
	// Check for existence of key in dictionary
	v, ok := pq.dict[member]
	if !ok || pq.dict == nil {
		return _const.ErrKeyNotFound
	}

	// Delete Node from the skip list and dictionary
	pq.skipList.delete(v.score, member)
	delete(pq.dict, member)
	pq.size--

	return nil
}

func (pq *ZSetNodes) exists(score int, member string) bool {
	v, ok := pq.dict[member]

	return ok && v.score == score
}

// Bytes encodes the ZSetNodes instance into bytes using MessagePack
// binary serialization format. The encoded bytes can be used for
// storage or transmission. If the encoding operation fails, an
// error is returned.
func (pq *ZSetNodes) Bytes() ([]byte, error) {
	var msgPack = encoding.NewMessagePackEncoder()
	if encodingError := msgPack.Encode(pq); encodingError != nil {
		return nil, encodingError
	}
	return msgPack.Bytes(), nil
}

// FromBytes decodes the input byte slice into the ZSetNodes object using MessagePack.
// Returns an error if decoding fails, otherwise nil.
func (pq *ZSetNodes) FromBytes(b []byte) error {
	return encoding.NewMessagePackDecoder(b).Decode(pq)
}

// getZSetFromDB fetches and deserializes ZSetNodes from the database.
//
// Returns a pointer to the ZSetNodes and error, if any.
// If the key doesn't exist, both the pointer and the error will be nil.
// In case of deserialization errors, returns nil and the error.
func (zs *ZSetStructure) getZSetFromDB(key []byte) (*ZSetNodes, error) {
	dbData, err := zs.db.Get(key)

	// If key is not found, return nil for both; otherwise return the error.
	if err != nil {

		return nil, err
	}

	// Deserialize the data.
	var zSetValue ZSetNodes
	if err := encoding.DecodeMessagePack(dbData, zSetValue); err != nil {
		return nil, err
	}
	// return a pointer to the deserialized ZSetNodes, nil for the error
	return &zSetValue, nil
}

// setZSetToDB writes a ZSetNodes object to the database.
//
//	parameters:
//	key: This is a byte slice that is used as a key in the database.
//	zSetValue: This is a pointer to a ZSetNodes object that needs to be stored in the database.
//
// The function serializes the ZSetNodes object into MessagePack format. If an error occurs
// either during serialization or when writing to the database, that specific error is returned.
// If the process is successful, it returns nil.
func (zs *ZSetStructure) setZSetToDB(key []byte, zSetValue *ZSetNodes) error {
	val, err := encoding.EncodeMessagePack(zSetValue)
	if err != nil {
		return err
	}
	return zs.db.Put(key, val)
}

// UnmarshalBinary de-serializes the given byte slice into ZSetNodes instance
// it uses MessagePack format for de-serialization
// Returns an error if the decoding of size or insertion of node fails.
//
// Parameters:
// data : a slice of bytes to be decoded
//
// Returns:
// An error that will be nil if the function succeeds.
func (p *ZSetNodes) UnmarshalBinary(data []byte) (err error) {
	// NewMessagePackDecoder creates a new MessagePack decoder with the provided data
	dec := encoding.NewMessagePackDecoder(data)

	var size int
	// Decode the size of the data structure
	if err = dec.Decode(&size); err != nil {
		return err // error handling if something goes wrong with decoding
	}

	// Iterate through each node in the data structure
	for i := 0; i < size; i++ {
		// Create an empty instance of SkipListNodeValue for each node
		slValue := SkipListNodeValue{}

		// Decode each node onto the empty SkipListNodeValue instance
		if err = dec.Decode(&slValue); err != nil {
			return err // error handling if something goes wrong with decoding
		}

		// Insert the decoded node into the ZSetNodes instance
		if err = p.InsertNode(slValue.score, slValue.member, slValue.value); err != nil {
			return err
		}
	}
	return // if all nodes are correctly decoded and inserted, return with nil error
}

// MarshalBinary serializes the ZSetNodes instance into a byte slice.
// It uses MessagePack format for serialization
// Returns the serialized byte slice and an error if the encoding fails.
func (d *ZSetNodes) MarshalBinary() (_ []byte, err error) {

	// Initializing the MessagePackEncoder
	enc := encoding.NewMessagePackEncoder()

	// Encoding the size attribute of d (i.e., d.size). The operation could fail, thus we check for an error.
	// An error, if occurred, will be returned immediately, hence the flow of execution stops here.
	err = enc.Encode(d.size)
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
	x := d.skipList.tail
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
	// ZSetNodes. Now, we return the bytes from the encoder along with any error that might have occurred
	// during the encoding (should be nil if everything went fine).
	return enc.Bytes(), err
}

// UnmarshalBinary de-serializes the given byte slice into SkipListNodeValue instance
// It uses the MessagePack format for de-serialization
// Returns an error if the decoding of Key, Score, or Value fails.
func (p *SkipListNodeValue) UnmarshalBinary(data []byte) (err error) {
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
// the SkipListNodeValue object into a byte array.
func (d *SkipListNodeValue) MarshalBinary() (_ []byte, err error) {

	// The NewMessagePackEncoder function is called to create a new
	// MessagePack encoder.
	enc := encoding.NewMessagePackEncoder()

	// Then, we try to encode the 'key' field of the SkipListNodeValue
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
	// final byte slice which represents the encoded SkipListNodeValue
	// and a nil error.
	return enc.Bytes(), err
}
