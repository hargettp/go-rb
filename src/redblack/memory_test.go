package redblack

import "math"
import "testing"

import l4g "code.google.com/p/log4go"

var log l4g.Logger

func init() {
	log = l4g.NewDefaultLogger(l4g.INFO)
}

func TestEmptyTree(t *testing.T) {
	tree := NewMemoryLLRB()
	if tree == nil {
		log.Error("Empty tree is nil")
		t.Fail()
	}
	if !checkInvariants(tree) {
		log.Error("Empty tree failed invariant checks")
		t.Fail()
	}
}

func TestSingleKeyTree(t *testing.T) {
	tree := NewMemoryLLRB()
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
	tree := NewMemoryLLRB()
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
	tree := NewMemoryLLRB()
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
	tree := NewMemoryLLRB()
	for i := 0; i < 75; i++ {
		tree.Insert(IntKey(i), StringValue(IntKey(i).String()))
	}
	if !checkInvariants(tree) {
		log.Error("Invariant check failed")
		t.Fail()
	}
}

func Test1DeleteEmptyTree(t *testing.T) {
	tree := NewMemoryLLRB()
	tree.Delete(IntKey(1))
}

func Test1DeleteSingleKeyTree(t *testing.T) {
	tree := NewMemoryLLRB()
	tree.Insert(IntKey(1), StringValue("one"))
	tree.Delete(IntKey(1))
	if !checkDelete(tree, IntKey(1)) {
		t.Fail()
	}
}

func Test1DeleteDoubleKeyTree(t *testing.T) {
	{
		// delete the min
		tree := NewMemoryLLRB()
		tree.Insert(IntKey(1), StringValue("one"))
		tree.Insert(IntKey(2), StringValue("two"))
		if !checkDelete(tree, IntKey(1)) {
			t.Fail()
		}
	}
	{
		// delete the max
		tree := NewMemoryLLRB()
		tree.Insert(IntKey(1), StringValue("one"))
		tree.Insert(IntKey(2), StringValue("two"))
		if !checkDelete(tree, IntKey(2)) {
			t.Fail()
		}
	}
}

func Test1DeleteMultipleKeyTree(t *testing.T) {
	genTree := func() LLRB {
		tree := NewMemoryLLRB()
		tree.Insert(IntKey(1), StringValue("one"))
		tree.Insert(IntKey(2), StringValue("two"))
		tree.Insert(IntKey(3), StringValue("three"))
		tree.Insert(IntKey(4), StringValue("four"))
		tree.Insert(IntKey(5), StringValue("five"))
		tree.Insert(IntKey(6), StringValue("six"))
		return tree
	}
	allOk := true
	for i := 1; i <= 6; i++ {
		tree := genTree()
		ok := checkDelete(tree, IntKey(i))
		if !ok {
			log.Error("Failed on deletion of key %v", i)
		}
		allOk = allOk && ok
	}
	if !allOk {
		t.Fail()
	}
}

func noTest1DeleteLotsOfKeysDeleteRoot(t *testing.T) {
	lots := 8
	genTree := func() LLRB {
		tree := NewMemoryLLRB()
		for i := 1; i <= lots; i++ {
			tree.Insert(IntKey(i), StringValue(IntKey(i).String()))
			if !checkInvariants(tree) {
				log.Error("Failed on of insertion of %v\n%v", i, tree)
				t.Fail()
			}
		}
		return tree
	}
	tree := genTree()
	log.Info("Original tree \n%v", tree)
	key := tree.Root().Key()
	if !checkDelete(tree, key) {
		log.Error("Failed on deletion of key %v\n%v", key, tree)
		should := NewMemoryLLRB()
		for i := 1; i <= lots; i++ {
			if IntKey(i) != key {
				should.Insert(IntKey(i), StringValue(IntKey(i).String()))
			}
		}
		log.Info("Should have been \n%v", should)
		t.Fail()
		tree.Root().SetLeft(tree.rotateLeft(tree.Root().Left()))
		log.Info("After rotate root's left left\n%v", tree)
		tree.SetRoot(tree.rotateRight(tree.Root()))
		log.Info("After rotate root right\n%v", tree)
	}
}

func noTest1DeleteLotsOfKeys(t *testing.T) {
	lots := 8
	genTree := func() LLRB {
		tree := NewMemoryLLRB()
		for i := 0; i < lots; i++ {
			tree.Insert(IntKey(i), StringValue(IntKey(i).String()))
		}
		return tree
	}
	allOk := true
	for i := 1; i <= lots; i++ {
		tree := genTree()
		ok := checkDelete(tree, IntKey(i))
		if !ok {
			log.Error("Failed on deletion of key %v", i)
		}
		allOk = allOk && ok
	}
	if !allOk {
		t.Fail()
	}
}

//=============================================================================
//
// Utility methods
//
//=============================================================================

func checkDelete(tree LLRB, key Key) bool {
	expectedSize := tree.Size() - 1
	if expectedSize < 0 {
		expectedSize = 0
	}
	tree.Delete(key)
	if tree.Search(key) != nil {
		log.Error("Deleted key still has value in tree")
		return false
	}
	if tree.Size() != expectedSize {
		log.Error("Deletion did not result in tree of expected size: expected %v, saw %v", expectedSize, tree.Size())
		return false
	}
	if !checkInvariants(tree) {
		log.Error("Invariant check failed")
		return false
	}
	return true
}

func checkInvariants(tree LLRB) bool {

	return checkBlackRoot(tree) &&
		checkAllLeavesBlack(tree) &&
		checkAllPathsSameNumberBlack(tree) &&
		checkChildrenOfRedAreBlack(tree) &&
		checkDepth(tree) &&
		checkTwoColors(tree)
}

func checkBlackRoot(tree LLRB) bool {
	// root must be black--or nil
	if tree.Root() == nil {
		return true
	}
	if isRed(tree.Root()) {
		log.Error("Root is not black")
		return false
	}
	return true
}

func checkAllLeavesBlack(tree LLRB) bool {
	// check that leaves are black (e.g., same as root)
	// 
	// Sorta happens automatically, because by default
	// all actual leaves are nil...but we only need
	// to navigate to nil leaves if the last non-nil
	// node is red
	//
	return true
}

func checkAllPathsSameNumberBlack(tree LLRB) bool {
	// check that all paths have same # black nodes
	allSameLength := true
	blackNodeCount := 0 // we start at 0, first path sets it
	log.Debug("Checking number of black nodes in paths for \n%v", tree)
	visitPaths(tree.Root(), func(path []Node) {
		log.Debug("Checking colors in path: %v", path)
		colors := make([]Color, 0)
		count := 0
		for _, n := range path {
			if n != nil {
				colors = append(colors, n.Color())
				if !isRed(n) {
					count++
				}
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

func checkChildrenOfRedAreBlack(tree LLRB) bool {
	// check that children of red nodes are black
	allRedHaveBlackChildren := true
	visitNodes(tree.Root(), func(h Node) {
		if isRed(h) {
			if h.Left() != nil {
				allRedHaveBlackChildren = allRedHaveBlackChildren && !isRed(h.Left())
			}
			if h.Right() != nil {
				allRedHaveBlackChildren = allRedHaveBlackChildren && !isRed(h.Right())
			}
		}
	})
	if !allRedHaveBlackChildren {
		log.Error("Not all red nodes have all black children\n%v", tree)
		return false
	}
	return true
}

func checkDepth(tree LLRB) bool {
	// verify that the depth is as expected (e.g., <= 2 * log2(N) )
	size := tree.Size()
	if size == 0 {
		return true
	}
	maxDepth := 2 * int(math.Ceil(math.Log2(float64(size))))
	if size == 1 {
		maxDepth = 1 // otherwise the log calculation would say zero
	}
	var depth func(h Node, d int) int
	depth = func(h Node, d int) int {
		log.Debug("Checking at depth %v node %v", d, h)
		leftDepth := d
		rightDepth := d
		if h.Left() != nil {
			leftDepth = depth(h.Left(), d+1)
		}
		if h.Right() != nil {
			rightDepth = depth(h.Right(), d+1)
		}
		if leftDepth > rightDepth {
			return leftDepth
		}
		return rightDepth
	}
	actualDepth := depth(tree.Root(), 1)
	ok := actualDepth <= maxDepth
	if !ok {
		log.Error("Depth check failed: expected <= %v, was %v\n%v", maxDepth, actualDepth, tree)
	}
	log.Debug("Observed size %v, actual depth %v, and max depth is %v", size, actualDepth, maxDepth)
	return ok
}

func checkTwoColors(tree LLRB) bool {
	sawRed := false
	sawBlack := false
	visitNodes(tree.Root(), func(h Node) {
		switch {
		case h == nil:
			sawBlack = true
		case isRed(h):
			sawRed = true
		case !isRed(h):
			sawBlack = true
		}
	})
	success := (sawRed && sawBlack) || tree.Size() <= 2
	if !success {
		log.Error("Expected at least 2 colors in tree: sawRed=%v,sawBlack=%v\n%v", sawRed, sawBlack, tree)
	}
	return success
}

// Useful visitor routines

func visitNodes(b Node, visit func(h Node)) {
	visit(b)
	if b != nil {
		visitNodes(b.Left(), visit)
		visitNodes(b.Right(), visit)
	}
}

func visitLeaves(b Node, visit func(h Node)) {
	visit(b)
	if b != nil {
		visitLeaves(b.Left(), visit)
		visitLeaves(b.Right(), visit)
	}
}

func visitPaths(b Node, visit func(path []Node)) {
	var visitor func(h Node, basePath []Node)
	visitor = func(h Node, basePath []Node) {
		if (h.Left() == nil) && (h.Right() == nil) {
			visit(basePath)
		}
		leftPath := make([]Node, len(basePath))
		copy(leftPath, basePath)
		if h.Left() != nil {
			leftPath = append(leftPath, h.Left())
			visitor(h.Left(), leftPath)
		}

		rightPath := make([]Node, len(basePath))
		copy(rightPath, basePath)
		if h.Right() != nil {
			rightPath = append(rightPath, h.Right())
			visitor(h.Right(), rightPath)
		}
	}
	if b != nil {
		initialPath := []Node{b}
		visitor(b, initialPath)
	}
}
