package redblack

import "fmt"

type memoryLLRB struct {
	root *memoryNode
}

type memoryNode struct {
	key         Key
	value       Value
	left, right *memoryNode
	color       Color
}

// Node implementation

func (h *memoryNode) Key() Key {
	return h.key
}

func (h *memoryNode) Value() Value {
	return h.value
}

func (h *memoryNode) SetValue(value Value) {
	h.value = value
}

func (h *memoryNode) IsRed() bool {
	if h == nil {
		return false
	}
	return h.Color() == RED
}

func (h *memoryNode) Left() Node {
	if h == nil {
		return nil
	}
	return h.left
}

func (n *memoryNode) SetLeft(h Node) {
	n.left = h.(*memoryNode)
}

func (h *memoryNode) Right() Node {
	if h == nil {
		return nil
	}
	return h.right
}

func (n *memoryNode) SetRight(h Node) {
	n.right = h.(*memoryNode)
}

func (h *memoryNode) FlipColors() {
	h.SetColor(!h.Color())
	h.left.SetColor(!h.left.Color())
	h.right.SetColor(!h.right.Color())
}

func (h *memoryNode) Color() Color {
	if h == nil {
		return BLACK
	}
	return h.color
}

func (h *memoryNode) SetColor(c Color) {
	h.color = c
}

func (h *memoryNode) String() string {
	var leftKey, rightKey Key
	if h.left != nil {
		leftKey = h.left.key
	} else {
		leftKey = nil
	}
	if h.right != nil {
		rightKey = h.right.key
	} else {
		rightKey = nil
	}
	return fmt.Sprintf("key=%v,left={%v},right={%v},color=%v,value=%v",
		h.key, leftKey, rightKey, h.Color(), h.value)
}

// LLRB implementation

func NewMemoryLLRB() LLRB {
	return &memoryLLRB{root: nil}
}

func (tree *memoryLLRB) NewNode(key Key, value Value) Node {
	return &memoryNode{key: key, value: value, color: RED}
}

func (tree *memoryLLRB) Search(key Key) Value {
	x := tree.Root()
	// NOTE this is a check for the sentinel
	for x != nil {
		// cmp := key.CompareTo(x.key)
		cmp := key.Compare(x.Key())
		if cmp == 0 {
			return x.Value()
		} else if cmp < 0 {
			x = x.Left()
		} else if cmp > 0 {
			x = x.Right()
		}
	}
	return nil
}

func (tree *memoryLLRB) Insert(key Key, value Value) {
	tree.SetRoot(tree.insert(tree.root, key, value))
	tree.Root().SetColor(BLACK)
}

func (tree *memoryLLRB) Delete(key Key) {
	tree.SetRoot(tree.delete(tree.root, key))
	tree.Root().SetColor(BLACK)
}

func (tree *memoryLLRB) DeleteMin() {
	tree.SetRoot(tree.deleteMin(tree.root))
	tree.Root().SetColor(BLACK)
}

func (tree *memoryLLRB) Size() int {
	var count func(h Node) int
	count = func(h Node) int {
		if h != nil {
			return 1 + count(h.Left()) + count(h.Right())
		}
		return 0
	}
	return count(tree.Root())
}

func (tree *memoryLLRB) Root() Node {
	return tree.root
}

func (tree *memoryLLRB) SetRoot(h Node) {
	tree.root = h.(*memoryNode)
}

func (tree *memoryLLRB) String() string {
	var visit func(h *memoryNode, nix int) string
	visit = func(h *memoryNode, nix int) string {
		if h == nil {
			return ""
		}
		var lix, rix int
		if h.left != nil {
			lix = nix + 1
		}
		if h.right != nil {
			rix = nix + 2
		}
		s := fmt.Sprintf("#%v: key=%v,color=%v,left=#%v,right=#%v,value=%v\n",
			nix, h.key, h.Color(), lix, rix, h.value)
		s += visit(h.left, lix)
		s += visit(h.right, rix)
		return s
	}
	return visit(tree.root, 1)
}

// LLRB internals

func (tree *memoryLLRB) insert(h Node, key Key, value Value) Node {
	// NOTE this is a check for the sentinel
	if h.(*memoryNode) == nil {
		return tree.NewNode(key, value)
	}
	if h.Left() != nil && h.Left().IsRed() && h.Right() != nil && h.Right().IsRed() {
		h.FlipColors()
	}
	cmp := key.Compare(h.Key())
	if cmp == 0 {
		h.SetValue(value)
	} else if cmp < 0 {
		h.SetLeft(tree.insert(h.Left(), key, value))
	} else if cmp > 0 {
		h.SetRight(tree.insert(h.Right(), key, value))
	}
	if h.Right().IsRed() && !h.Left().IsRed() {
		h = tree.rotateLeft(h)
	}
	if h.Left().IsRed() && h.Left().Left().IsRed() {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *memoryLLRB) rotateLeft(h Node) Node {
	x := h.Right()
	h.SetRight(x.Left())
	x.SetLeft(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	return x
}

func (tree *memoryLLRB) rotateRight(h Node) Node {
	x := h.Left()
	h.SetLeft(x.Right())
	x.SetRight(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	return x
}

func (tree *memoryLLRB) delete(h *memoryNode, key Key) *memoryNode {
	return nil
}

func (tree *memoryLLRB) moveRedLeft(h Node) Node {
	h.FlipColors()
	if h.Right().Left().IsRed() {
		h.SetRight(tree.rotateRight(h.Right()))
		h = tree.rotateLeft(h)
		h.FlipColors()
	}
	return h
}

func (tree *memoryLLRB) moveRedRight(h Node) Node {
	h.FlipColors()
	if h.Left().Left().IsRed() {
		h = tree.rotateRight(h)
		h.FlipColors()
	}
	return h
}

func (tree *memoryLLRB) fixUp(h Node) Node {
	if h.Right().IsRed() && !h.Left().IsRed() {
		h = tree.rotateLeft(h)
	}
	if h.Left().IsRed() && h.Left().Left().IsRed() {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *memoryLLRB) deleteMin(h Node) Node {
	// NOTE this is a check for the sentinel
	if h.Left() == nil {
		return nil
	}
	if !h.Left().IsRed() && !h.Left().Left().IsRed() {
		h = tree.moveRedLeft(h)
	}
	h.SetLeft(tree.deleteMin(h.Left()))
	return tree.fixUp(h)
}
