package btree

// Why B-Tree
// 1. keeps keys in sorted order for sequential traversing
// 2. uses a hierarchical index to minimize the number of disk reads
// 4. uses partially full blocks to speed up insertions and deletions
// 5. keeps the index balanced with a recursive algorithm

type node[K comparable, V any] struct {
	keys     []K           // Sorted array of keys
	children []*node[K, V] // Array of child pointers
	values   []V           // Array of values associated with keys
}

// return true if node is a leaf
func (n *node[K, V]) leaf() bool {
	return n.children == nil
}

// initalize a new node
func newNode[K comparable, V any](isLeaf bool) *node[K, V] {
	n := &node[K, V]{
		keys:   []K{},
		values: []V{},
	}
	if !isLeaf {
		n.children = []*node[K, V]{}
	}
	return n
}

// funcCmp determines how to order a type K
// Comparability: The keys must support comparison operations (e.g., <, <=, =, >, >=).
type funcCmp[K comparable] func(K, K) int

// definition of Btree structure
type BTree[K comparable, V any] struct {
	root   *node[K, V] // root node of the B-Tree
	degree int         // Minimum degree (minimum number of keys) of the B-tree
	less   funcCmp[K]
}

func New[K comparable, V any](degree int, less funcCmp[K]) *BTree[K, V] {
	return &BTree[K, V]{degree: degree, less: less}
}

//  1. If the B-Tree is empty:
//     a .Allocate a root node, and insert the key.
//  2. If the B-Tree is not empty:
//     a. Find the proper node for insertion.
//     b. If the node is not full:
//     i. Insert the key in ascending order.
//     c. If the node is full:
//     i. Split the node at the median.
//     ii. Push the median key upward, and make the left keys a left child node and the right keys a right child node.
func (t *BTree[K, V]) isFull(n *node[K, V]) bool {
	return len(n.keys) == 2*t.degree-1
}

func (t *BTree[K, V]) Insert(key K, value V) {
	if t.root == nil {
		t.root = newNode[K, V](true)
		t.root.keys = append(t.root.keys, key)
		t.root.values = append(t.root.values, value)
		return
	}
	if t.isFull(t.root) {
	}
}
