package btree

import "fmt"

// Why B-Tree
// 1. keeps keys in sorted order for sequential traversing
// 2. uses a hierarchical index to minimize the number of disk reads
// 4. uses partially full blocks to speed up insertions and deletions
// 5. keeps the index balanced with a recursive algorithm

// Node is a single element within the tree
type Node[K comparable, V any] struct {
	parent   *Node[K, V]   // this will be helpfull
	entries  []*Item[K, V] // Sorted array of keys
	children []*Node[K, V] // Array of child pointers
}

// Entry represents the key-value pair contained within nodes
type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// definition of Btree structure
type BTree[K comparable, V any] struct {
	root  *Node[K, V] // root node of the B-Tree
	order int         // Minimum degree (minimum number of keys) of the B-tree
	less  funcCmp[K]
	size  int
}

// NewBTree creates a new B-tree with the given degree
func NewBTree[K comparable, V any](order int, less funcCmp[K]) *BTree[K, V] {
	if order < 2 {
		panic("Invalid degree, should be at least 3")
	}
	tree := new(BTree[K, V])
	tree.less = less
	tree.order = order
	return tree
}

// SET OF HELPER FUNCITONS

// funcCmp determines how to order a type K
// Comparability: The keys must support comparison operations (e.g., <, <=, =, >, >=).
type funcCmp[K comparable] func(K, K) int

func (t *BTree[K, V]) isLeaf(node *Node[K, V]) bool {
	return len(node.children) == 0
}

func (t *BTree[K, V]) shouldSplit(node *Node[K, V]) bool {
	return len(node.entries) > t.maxEntries()
}

func (t *BTree[K, V]) maxChildren() int {
	return t.order
}

func (t *BTree[K, V]) minChildren() int {
	return (t.order + 1) / 2 // ceil(m/2)
}

func (t *BTree[K, V]) maxEntries() int {
	return t.maxChildren() - 1
}

func (t *BTree[K, V]) minEntries() int {
	return t.minChildren() - 1
}

func (t *BTree[K, V]) middle() int {
	return (t.order - 1) / 2 // "-1" to favor right nodes to have more keys when splitting
}

func (t *BTree[K, V]) isEmpty() bool {
	return t.size == 0
}

// Less is a convinience function that perfomrs comparsion between two items
// using hte same "less" function provided to New
func (t *BTree[K, V]) Less(a, b K) int {
	return t.less(a, b)
}

// helper function to avoid repetive code, the goal is search the index where the key is stored, also use a binary search to perform o(log n) searches
// and return the correct index and a boolean to indicate that it was found
// in case that the key is not there we return that is not found and the right place where it should be placed
func (t *BTree[K, V]) searchKeyIndex(node *Node[K, V], key K) (index int, found bool) {
	low, high := 0, len(node.entries)-1
	var mid int
	for low <= high {
		mid = (high + low) / 2
		if t.Less(key, node.entries[mid].Key) == 0 {
			return mid, true
		} else if t.Less(key, node.entries[mid].Key) > 0 {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return low, false
}

func (t *BTree[K, V]) Put(key K, value V) {
	entry := &Item[K, V]{Key: key, Value: value}
	if t.root == nil { // empty tree
		t.root = &Node[K, V]{entries: []*Item[K, V]{entry}, children: []*Node[K, V]{}}
		t.size++
		return
	}
	if t.insert(t.root, entry) {
		t.size++
	}
}

//	1.When inserting into a leaf node,
//	  we simply add the key-value pair to the node (maintaining order).
//	2.When inserting into an internal node,
//	  we need to traverse down to a leaf node where the actual insertion will occur.

// TODO
// 1. insertLeaf method
// 2. insertInternal method
func (t *BTree[K, V]) insert(node *Node[K, V], entry *Item[K, V]) bool {
	if t.isLeaf(node) {
		return t.insertLeaf(node, entry)
	}
	return t.insertInternal(node, entry)
}

func (t *BTree[K, V]) insertInternal(node *Node[K, V], entry *Item[K, V]) bool {
	insertIndex, found := t.searchKeyIndex(node, entry.Key)
	if found {
		node.entries[insertIndex] = entry
		return false
	}
	if insertIndex >= len(node.children) {
		fmt.Printf(
			"insert index:%d,len of node children slice = %d\n",
			insertIndex,
			len(node.children),
		)
		panic("insertIndex equals len of node children slice")
	}
	return t.insert(node.children[insertIndex], entry)
}

func (t *BTree[K, V]) insertLeaf(node *Node[K, V], entry *Item[K, V]) bool {
	insertIndex, found := t.searchKeyIndex(node, entry.Key)
	if found { // this is for the case when we change the value for an exisiting key
		node.entries[insertIndex] = entry
		return false
	}
	// Expand the slice to make space for the new entry
	node.entries = append(node.entries, nil)
	copy(node.entries[insertIndex+1:], node.entries[insertIndex:])
	node.entries[insertIndex] = entry

	// we need to check if after insertion is split and rebalacing needed
	t.split(node)

	return true
}

func (t *BTree[K, V]) split(node *Node[K, V]) {
	if !t.shouldSplit(node) {
		return
	}
	if node == t.root {
		t.splitRoot()
	} else {
		t.splitNonRoot(node)
	}
}

func (t *BTree[K, V]) splitNonRoot(node *Node[K, V]) {
	middle := t.middle()
	parent := node.parent

	left := &Node[K, V]{
		entries: append([]*Item[K, V](nil), node.entries[:middle]...),
		parent:  parent,
	}
	right := &Node[K, V]{
		entries: append([]*Item[K, V](nil), node.entries[middle+1:]...),
		parent:  parent,
	}

	// Move children from the node to be split into left and right nodes
	if !t.isLeaf(node) {
		left.children = append([]*Node[K, V](nil), node.children[:middle+1]...)
		right.children = append([]*Node[K, V](nil), node.children[middle+1:]...)
		setParent(left.children, left)
		setParent(right.children, right)
	}

	insertPosition, _ := t.searchKeyIndex(parent, node.entries[middle].Key)

	// Insert middle key into parent
	parent.entries = append(parent.entries, nil)
	copy(parent.entries[insertPosition+1:], parent.entries[insertPosition:])
	parent.entries[insertPosition] = node.entries[middle]

	// Set child left of inserted key in parent to the created left node
	parent.children[insertPosition] = left

	// Set child right of inserted key in parent to the created right node
	parent.children = append(parent.children, nil)
	copy(parent.children[insertPosition+2:], parent.children[insertPosition+1:])
	parent.children[insertPosition+1] = right

	t.split(parent)
}

func (t *BTree[K, V]) splitRoot() {
	middle := t.middle()

	left := &Node[K, V]{entries: append([]*Item[K, V](nil), t.root.entries[:middle]...)}
	right := &Node[K, V]{entries: append([]*Item[K, V](nil), t.root.entries[middle+1:]...)}

	// Move children from the node to be split into left and right nodes
	if !t.isLeaf(t.root) {
		left.children = append([]*Node[K, V](nil), t.root.children[:middle+1]...)
		right.children = append([]*Node[K, V](nil), t.root.children[middle+1:]...)
		setParent(left.children, left)
		setParent(right.children, right)
	}

	// Root is a node with one entry and two children (left and right)
	newRoot := &Node[K, V]{
		entries:  []*Item[K, V]{t.root.entries[middle]},
		children: []*Node[K, V]{left, right},
	}

	left.parent = newRoot
	right.parent = newRoot
	t.root = newRoot
}

func setParent[K comparable, V any](nodes []*Node[K, V], parent *Node[K, V]) {
	for _, node := range nodes {
		node.parent = parent
	}
}

// searchRecursively searches for a key starting from a specific node recursively
// It returns the node containing the key (if found), the index of the key in the node, and a boolean indicating if the key was found.
func (t *BTree[K, V]) searchRecursively(
	startNode *Node[K, V],
	key K,
) (node *Node[K, V], index int, found bool) {
	if t.isEmpty() {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = t.searchKeyIndex(node, key)
		if found {
			return node, index, true
		}
		if t.isLeaf(node) {
			return nil, -1, false
		}
		node = node.children[index]
	}
}

func (t *BTree[K, V]) Get(key K) (value V, found bool) {
	node, index, found := t.searchRecursively(t.root, key)
	if found {
		return node.entries[index].Value, true
	}
	return value, false
}

//func (t *BTree[K, V]) Delete(key K) error {
//	if t.root == nil {
//		return fmt.Errorf("Tree is empty")
//	}
//	node, index, found := t.searchRecursively(t.root, key)
//	if found {
//		t.remove(node, index)
//		t.size--
//		return nil
//	}
//	return fmt.Errorf("Key is not in the tree")
//}

func (t *BTree[K, V]) Delete(key K) error {
	if t.root == nil {
		return fmt.Errorf("Tree is empty")
	}
	node, index, found := t.searchRecursively(t.root, key)
	if !found {
		return fmt.Errorf("Key is not in the tree")
	}
	t.remove(node, index)
	t.size--

	// If the root node is empty after removal, make its only child the new root
	if len(t.root.entries) == 0 && len(t.root.children) > 0 {
		t.root = t.root.children[0]
		t.root.parent = nil
	}

	return nil
}

func (t *BTree[K, V]) remove(node *Node[K, V], index int) {
	if t.isLeaf(node) {
		t.removeFromLeaf(node, index)
	} else {
		t.removeFromNonLeaf(node, index)
	}
	t.rebalance(node) // reblance if necessary
}

func (t *BTree[K, V]) removeFromLeaf(node *Node[K, V], index int) {
	// Remove the entry at the given index
	copy(node.entries[index:], node.entries[index+1:])
	node.entries = node.entries[:len(node.entries)-1]
}

func (t *BTree[K, V]) removeFromNonLeaf(node *Node[K, V], index int) {
	if len(node.children[index].entries) >= t.minEntries() { // LST
		predecessor := t.predecessor(node, index)
		node.entries[index] = predecessor
		t.remove(node.children[index], len(node.entries)-1)
	} else if len(node.children[index+1].entries) >= t.minEntries() { // RST
		succesor := t.successor(node, index)
		node.entries[index] = succesor
		t.remove(node.children[index+1], 0)
	} else {
		t.mergeChildren(node, index)
		key, _ := t.searchKeyIndex(node, node.entries[index].Key)
		t.remove(node.children[index], key)
	}
}

func (t *BTree[K, V]) predecessor(node *Node[K, V], index int) *Item[K, V] {
	current := node.children[index]
	for !t.isLeaf(current) {
		current = current.children[len(current.children)-1]
	}
	return current.entries[len(current.entries)-1]
}

func (t *BTree[K, V]) successor(node *Node[K, V], index int) *Item[K, V] {
	current := node.children[index+1]
	for !t.isLeaf(current) {
		current = current.children[0]
	}
	return current.entries[0]
}

func (t *BTree[K, V]) mergeChildren(parent *Node[K, V], index int) {
	// primero borro del parent node y acomodo
	// despues hago el merge entre los nodos izq y derechos
	// actualizo punteros

	leftChild, rightChild := parent.children[index], parent.children[index+1]

	leftChild.entries = append(leftChild.entries, parent.entries[index])
	leftChild.entries = append(leftChild.entries, rightChild.entries...)
	leftChild.children = append(leftChild.children, rightChild.children...)

	setParent(rightChild.children, leftChild)

	copy(parent.entries[index:], parent.entries[index+1:])
	copy(parent.children[index+1:], parent.children[index+2:])
	parent.entries = parent.entries[:len(parent.entries)-1]
	parent.children = parent.children[:len(parent.children)-1]
}

func (t *BTree[K, V]) rebalance(node *Node[K, V]) {
	if len(node.entries) >= t.minEntries() {
		return
	}

	if node == t.root {
		return
	}

	parent := node.parent
	index := t.getChildIndex(parent, node)

	// Try to borrow from left sibling
	if index > 0 && len(parent.children[index-1].entries) > t.minEntries() {
		t.borrowFromLeft(node, index)
	} else if index < len(parent.children)-1 && len(parent.children[index+1].entries) > t.minEntries() {
		// Try to borrow from right sibling
		t.borrowFromRight(node, index)
	} else if index > 0 {
		// Merge with left sibling
		t.mergeChildren(parent, index-1)
	} else {
		// Merge with right sibling
		t.mergeChildren(parent, index)
	}

	t.rebalance(parent)
}

func (t *BTree[K, V]) borrowFromLeft(node *Node[K, V], index int) {
	parent := node.parent
	leftSibling := parent.children[index-1]

	// Move the separating key from the parent to the beginning of the node
	node.entries = append([]*Item[K, V]{parent.entries[index-1]}, node.entries...)

	// Move the last key from the left sibling to the parent
	parent.entries[index-1] = leftSibling.entries[len(leftSibling.entries)-1]
	leftSibling.entries = leftSibling.entries[:len(leftSibling.entries)-1]

	if !t.isLeaf(node) {
		// Move the last child pointer from the left sibling to the beginning of the node
		node.children = append(
			[]*Node[K, V]{leftSibling.children[len(leftSibling.children)-1]},
			node.children...)
		leftSibling.children = leftSibling.children[:len(leftSibling.children)-1]
		node.children[0].parent = node
	}
}

func (t *BTree[K, V]) borrowFromRight(node *Node[K, V], index int) {
	parent := node.parent
	rightSibling := parent.children[index+1]

	// Move the separating key from the parent to the end of the node
	node.entries = append(node.entries, parent.entries[index])

	// Move the first key from the right sibling to the parent
	parent.entries[index] = rightSibling.entries[0]
	rightSibling.entries = rightSibling.entries[1:]

	if !t.isLeaf(node) {
		// Move the first child pointer from the right sibling to the end of the node
		node.children = append(node.children, rightSibling.children[0])
		rightSibling.children = rightSibling.children[1:]
		node.children[len(node.children)-1].parent = node
	}
}

func (t *BTree[K, V]) getChildIndex(parent *Node[K, V], child *Node[K, V]) int {
	for i, c := range parent.children {
		if c == child {
			return i
		}
	}
	return -1
}
