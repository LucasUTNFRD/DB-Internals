package btree

type node[K comparable, V any] struct {
	items    []K
	children []*node[K, V]
	isLeaf   bool
	n        int // number of keys
}

// funcCmp determines how to order a type K
// Comparability: The keys must support comparison operations (e.g., <, <=, =, >, >=).
type funcCmp[K comparable] func(K, K) int

// definition of Btree structure
type BTree[K comparable, V any] struct {
	root  *node[K, V]
	order int
	less  funcCmp[K]
}

func New[K comparable, V any](order int, less funcCmp[K]) *BTree[K, V] {
	return &BTree[K, V]{order: order, less: less}
}

func (t *BTree[K, V]) Insert(key K, val V) {
}
