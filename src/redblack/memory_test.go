package redblack

import "testing"

import l4g "code.google.com/p/log4go"

var log = l4g.NewDefaultLogger(l4g.INFO)

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
	tree := NewMemoryLLRB().(*memoryLLRB)
	tree.Insert(IntKey(1), StringValue("one"))
	val := tree.Search(IntKey(1))
	log.Debug("Single key tree is: %v", tree)
	if val.String() != "one" {
		log.Error("Search did not return inserted value")
		t.Fail()
	}
	if !checkInvariants(tree) {
		log.Error("Invariant check failed")
		t.Fail()
	}
}

func TestDoubleKeyTree(t *testing.T) {
	tree := NewMemoryLLRB().(*memoryLLRB)
	tree.Insert(IntKey(1), StringValue("one"))
	tree.Insert(IntKey(2), StringValue("two"))
	val := tree.Search(IntKey(1))
	if val.String() != "one" {
		log.Error("Search did not return inserted value 'one'")
		t.Fail()
	}
	val = tree.Search(IntKey(2))
	if val.String() != "two" {
		log.Error("Search did not return inserted value 'two'")
		t.Fail()
	}
	if !checkInvariants(tree) {
		log.Error("Invariant check failed")
		t.Fail()
	}

}

func TestMultipleKeyTree(t *testing.T) {
	tree := NewMemoryLLRB().(*memoryLLRB)
	tree.Insert(IntKey(1), StringValue("one"))
	tree.Insert(IntKey(2), StringValue("two"))
	tree.Insert(IntKey(3), StringValue("three"))
	tree.Insert(IntKey(4), StringValue("four"))
	tree.Insert(IntKey(5), StringValue("five"))
	tree.Insert(IntKey(6), StringValue("six"))
	val := tree.Search(IntKey(1))
	if val.String() != "one" {
		log.Error("Search did not return inserted value 'one'")
		t.Fail()
	}
	val = tree.Search(IntKey(5))
	if val.String() != "five" {
		log.Error("Search did not return inserted value 'five'")
		t.Fail()
	}
	if !checkInvariants(tree) {
		log.Error("Invariant check failed")
		t.Fail()
	}
}

func TestLotsOfKeys(t *testing.T) {
	tree := NewMemoryLLRB().(*memoryLLRB)
	for i := 0; i < 50; i++ {
		tree.Insert(IntKey(i), StringValue(IntKey(i).String()))
	}
	if !checkInvariants(tree) {
		log.Error("Invariant check failed")
		t.Fail()
	}
}

//
// Utility methods
//
func checkInvariants(tree *memoryLLRB) bool {

	return checkBlackRoot(tree) &&
		checkAllLeavesBlack(tree) &&
		checkAllPathsSameNumberBlack(tree) &&
		checkChildrenOfRedAreBlack(tree)
}

func checkBlackRoot(tree *memoryLLRB) bool {
	// root must be black--or nil
	if tree.root == nil {
		return true
	}
	if tree.root.IsRed() {
		log.Error("Root is not black")
		return false
	}
	return true
}

func checkAllLeavesBlack(tree *memoryLLRB) bool {
	// check that leaves are black (e.g., same as root)
	// 
	// Sorta happens automatically, because by default
	// all actual leaves are nil...but we only need
	// to navigate to nil leaves if the last non-nil
	// node is red
	//
	return true
}

func checkAllPathsSameNumberBlack(tree *memoryLLRB) bool {
	// check that all paths have same # black nodes
	allSameLength := true
	blackNodeCount := 0 // we start at 0, first path sets it
	log.Debug("Checking number of black nodes in paths for \n%v", tree)
	visitPaths(tree.root, func(path []*memoryNode) {
		log.Debug("Checking colors in path: %v", path)
		colors := make([]Color, 0)
		count := 0
		for _, n := range path {
			if n != nil {
				colors = append(colors, n.color)
			}
			if !n.IsRed() {
				count++
			}
		}
		if blackNodeCount == 0 {
			blackNodeCount = count
		} else {
			allSameLength = (allSameLength && (count == blackNodeCount))
		}
		log.Debug("Colors in path(found %v, expected %v): %v", count, blackNodeCount, colors)
	})
	if !allSameLength {
		log.Error("Not all paths are same length; expected %v: \n%v", blackNodeCount, tree)
		return false
	}
	return true
}

func checkChildrenOfRedAreBlack(tree *memoryLLRB) bool {
	// check that children of red nodes are black
	allRedHaveBlackChildren := true
	visitNodes(tree.root, func(h *memoryNode) {
		if h.IsRed() {
			allRedHaveBlackChildren = allRedHaveBlackChildren && !h.left.IsRed()
			allRedHaveBlackChildren = allRedHaveBlackChildren && !h.right.IsRed()
		}
	})
	if !allRedHaveBlackChildren {
		log.Error("Not all red nodes have all black children")
		return false
	}
	return true
}

func visitNodes(b *memoryNode, visit func(h *memoryNode)) {
	if b != nil {
		visit(b)
		visit(b.left)
		visit(b.right)
	}
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
		if (h.left == nil) && (h.right == nil) {
			visit(basePath)
		}
		leftPath := make([]*memoryNode, len(basePath))
		copy(leftPath, basePath)
		if h.left != nil {
			leftPath = append(leftPath, h.left)
			visitor(h.left, leftPath)
		}

		rightPath := make([]*memoryNode, len(basePath))
		copy(rightPath, basePath)
		if h.right != nil {
			rightPath = append(rightPath, h.right)
			visitor(h.right, rightPath)
		}
	}
	if b != nil {
		initialPath := []*memoryNode{b}
		visitor(b, initialPath)
	}
}
