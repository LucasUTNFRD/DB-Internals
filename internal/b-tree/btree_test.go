package btree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const minItems = 2

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

func TestBTreeInsertion(t *testing.T) {
	t.Log("Test for creating a Btree and adding an element to the root node")
	// Initialize a B-tree with a specific degree (t)
	btree := NewBTree[int, string](minItems, cmpInt) // Test case 1: Insert into an empty B-tree
	btree.Insert(10, "value10")
	require.EqualValues(t, 1, btree.NumOfItems())
}
