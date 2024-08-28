package btree

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
		panic("Invalid degree, should be at leas 3")
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

func (t *BTree[K, V]) isFull(node *Node[K, V]) bool {
	return len(node.entries) == t.maxEntries()
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

// Less is a convinience function that perfomrs comparsion between two items
// using hte same "less" function provided to New
func (t *BTree[K, V]) Less(a, b K) int {
	return t.less(a, b)
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

//  1.When inserting into a leaf node,
//    we simply add the key-value pair to the node (maintaining order).
//  2.When inserting into an internal node,
//    we need to traverse down to a leaf node where the actual insertion will occur.

func (t *BTree[K, V]) insert(node *Node[K, V], entry *Item[K, V]) bool {
	if t.isLeaf(node) {
		return t.insertLeaf(node, entry)
	}
	return t.insertInternal(node, entry)
}

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

func (t *BTree[K, V]) insertLeaf(node *Node[K, V], entry *Item[K, V]) bool {
	insertIndex, found := t.searchKeyIndex(node, entry.Key)
	if found { // this is for the case when we change the value for a exisiting key
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
	if node == t.node {
		t.splitRoot()
	} else {
		t.splitNonRoot(node)
	}
}

func (t *BTree[K, V]) insertInternal(node *Node[K, V], entry *Item[K, V]) bool {
	return true
}
