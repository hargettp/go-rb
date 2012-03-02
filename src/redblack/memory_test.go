package redblack

import "testing"

func TestEmptyTree(t *testing.T) {
	tree := NewMemoryLLRB().(*memoryLLRB)
	if tree == nil {
		t.Fail()
	}
	if !checkInvariants(tree) {
		t.Fail()
	}
}

func TestSingleKeyTree(t *testing.T) {
}

//
// Utility methods
//
func checkInvariants(tree *memoryLLRB) bool {

	// root must be black--or nil
	if tree.root == nil {
		return true
	}
	if tree.root.isRed() {
		return false
	}

	// check that leaves are black (e.g., same as root)
	// NOTE actually, funny thing, not sure we can't count
	// this way, since the "leaves" are just nil pointers
	allBlack := true
	visitLeaves(tree.root, func(h *memoryNode) {
		if h == nil {
			allBlack = allBlack && true
		} else {
			allBlack = allBlack && !h.isRed()
		}
	})
	if !allBlack {
		return false
	}

	// check that all paths have same # black nodes
	allSameLength := true
	blackNodeCount := 0 // we start at 0, first path sets it
	visitPaths(tree.root, func(path []*memoryNode) {
		count := 0
		for _, n := range path {
			if !n.isRed() {
				count++
			}
		}
		if blackNodeCount == 0 {
			blackNodeCount = count
		} else {
			allSameLength = allSameLength && (count == blackNodeCount)
		}
	})
	if !allSameLength {
		return false
	}

	// check that children of red nodes are black

	return true
}

func visitLeaves(b *memoryNode, visit func(h *memoryNode)) {
	visit(b)
	if b != nil {
		visitLeaves(b.left, visit)
		visitLeaves(b.right, visit)
	}
}

func visitPaths(b *memoryNode, visit func(path []*memoryNode)) {
	var visitor func(h *memoryNode, basePath []*memoryNode)
	visitor = func(h *memoryNode, basePath []*memoryNode) {
		if h != nil {
			visit(basePath)
			leftPath := make([]*memoryNode, len(basePath))
			copy(basePath, leftPath)
			if h.left != nil {
				leftPath = append(leftPath, h.left)
			}
			visitor(h.left, leftPath)

			rightPath := make([]*memoryNode, len(basePath))
			copy(basePath, rightPath)
			if h.right != nil {
				rightPath = append(rightPath, h.right)
			}
			visitor(h.right, rightPath)
		}
	}
	if b != nil {
		initialPath := []*memoryNode{b}
		visitor(b, initialPath)
	}
}
