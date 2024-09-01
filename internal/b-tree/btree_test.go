package btree

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	minItems   = 3
	benchItems = 16
)

func cmpString(s1, s2 string) int {
	if s1 == s2 {
		return 0
	} else if s1 < s2 {
		return -1
	}
	return 1
}

func cmpInt(i1, i2 int) int {
	if i1 == i2 {
		return 0
	} else if i1 < i2 {
		return -1
	}
	return 1
}

func TestBtreeEmpty(t *testing.T) {
	tree := NewBTree[int, string](3, cmpInt)
	assert.True(t, tree.isEmpty())

	tree.Put(1, "uno")
	assert.False(t, tree.isEmpty())
	entrie := &Item[int, string]{1, "uno"}
	//assert that the unique entrie is matched
	assert.EqualValues(t, tree.root.entries[0], entrie)
}

func TestBTreePut1(t *testing.T) {
	tree := NewBTree[int, int](3, cmpInt)
	assertValidTree(t, tree, 0)

	tree.Put(1, 0)
	assertValidTree(t, tree, 1)
	assertValidTreeNode(t, tree.root, 1, 0, []int{1}, false)

	tree.Put(2, 1)
	assertValidTree(t, tree, 2)
	assertValidTreeNode(t, tree.root, 2, 0, []int{1, 2}, false)

	tree.Put(3, 2)
	assertValidTree(t, tree, 3)
	assertValidTreeNode(t, tree.root, 1, 2, []int{2}, false)
	assertValidTreeNode(t, tree.root.children[0], 1, 0, []int{1}, true)
	//assertValidTreeNode(t, tree.root.children[1], 1, 0, []int{3}, true)

	//tree.Put(4, 2)
	//assertValidTree(t, tree, 4)
	//assertValidTreeNode(t, tree.root, 1, 2, []int{2}, false)
	//assertValidTreeNode(t, tree.root.children[0], 1, 0, []int{1}, true)
	//assertValidTreeNode(t, tree.root.children[1], 2, 0, []int{3, 4}, true)
	//
	//tree.Put(5, 2)
	//assertValidTree(t, tree, 5)
	//assertValidTreeNode(t, tree.root, 2, 3, []int{2, 4}, false)
	//assertValidTreeNode(t, tree.root.children[0], 1, 0, []int{1}, true)
	//assertValidTreeNode(t, tree.root.children[1], 1, 0, []int{3}, true)
	//assertValidTreeNode(t, tree.root.children[2], 1, 0, []int{5}, true)
	//
	//tree.Put(6, 2)
	//assertValidTree(t, tree, 6)
	//assertValidTreeNode(t, tree.root, 2, 3, []int{2, 4}, false)
	//assertValidTreeNode(t, tree.root.children[0], 1, 0, []int{1}, true)
	//assertValidTreeNode(t, tree.root.children[1], 1, 0, []int{3}, true)
	//assertValidTreeNode(t, tree.root.children[2], 2, 0, []int{5, 6}, true)
	//
	//tree.Put(7, 2)
	//assertValidTree(t, tree, 7)
	//assertValidTreeNode(t, tree.root, 1, 2, []int{4}, false)
	//assertValidTreeNode(t, tree.root.children[0], 1, 2, []int{2}, true)
	//assertValidTreeNode(t, tree.root.children[1], 1, 2, []int{6}, true)
	//assertValidTreeNode(t, tree.root.children[0].children[0], 1, 0, []int{1}, true)
	//assertValidTreeNode(t, tree.root.children[0].children[1], 1, 0, []int{3}, true)
	//assertValidTreeNode(t, tree.root.children[1].children[0], 1, 0, []int{5}, true)
	//assertValidTreeNode(t, tree.root.children[1].children[1], 1, 0, []int{7}, true)
}

func assertValidTree[K comparable, V any](t *testing.T, tree *BTree[K, V], expectedSize int) {
	if actualValue, expectedValue := tree.size, expectedSize; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for tree size", actualValue, expectedValue)
	}
}

func assertValidTreeNode[K comparable, V any](t *testing.T, node *Node[K, V], expectedEntries int, expectedChildren int, keys []K, hasParent bool) {
	if actualValue, expectedValue := node.parent != nil, hasParent; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for hasParent", actualValue, expectedValue)
	}
	if actualValue, expectedValue := len(node.entries), expectedEntries; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for entries size", actualValue, expectedValue)
	}
	if actualValue, expectedValue := len(node.children), expectedChildren; actualValue != expectedValue {
		t.Errorf("Got %v expected %v for children size", actualValue, expectedValue)
	}
	for i, key := range keys {
		if actualValue, expectedValue := node.entries[i].Key, key; actualValue != expectedValue {
			t.Errorf("Got %v expected %v for key", actualValue, expectedValue)
		}
	}
}
