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

func (h *memoryNode) Right() Node {
	if h == nil {
		return nil
	}
	return h.right
}

func (h *memoryNode) FlipColors() {
	h.SetColor(!h.Color())
	h.left.SetColor(!h.left.Color())
	h.right.SetColor(!h.right.Color())
}

func (h *memoryNode) Color() Color {
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

func (tree *memoryLLRB) Search(key Key) Value {
	x := tree.root
	// NOTE this is a check for the sentinel
	for x != nil {
		// cmp := key.CompareTo(x.key)
		cmp := key.Compare(x.key)
		if cmp == 0 {
			return x.value
		} else if cmp < 0 {
			x = x.left
		} else if cmp > 0 {
			x = x.right
		}
	}
	return nil
}

func (tree *memoryLLRB) Insert(key Key, value Value) {
	tree.root = tree.insert(tree.root, key, value)
	tree.root.SetColor(BLACK)
}

func (tree *memoryLLRB) Delete(key Key) {
	tree.root = tree.delete(tree.root, key)
	tree.root.SetColor(BLACK)
}

func (tree *memoryLLRB) DeleteMin() {
	tree.root = tree.deleteMin(tree.root)
	tree.root.SetColor(BLACK)
}

func (tree *memoryLLRB) Size() int {
	var count func(h Node) int
	count = func(h Node) int {
		if h != nil {
			return 1 + count(h.Left()) + count(h.Right())
		}
		return 0
	}
	return count(tree.root)
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

func (tree *memoryLLRB) insert(h *memoryNode, key Key, value Value) *memoryNode {
	// NOTE this is a check for the sentinel
	if h == nil {
		return &memoryNode{key: key, value: value, color: RED}
	}
	if h.left.IsRed() && h.right.IsRed() {
		h.FlipColors()
	}
	cmp := key.Compare(h.key)
	if cmp == 0 {
		h.value = value
	} else if cmp < 0 {
		h.left = tree.insert(h.left, key, value)
	} else if cmp > 0 {
		h.right = tree.insert(h.right, key, value)
	}
	if h.right.IsRed() && !h.left.IsRed() {
		h = tree.rotateLeft(h)
	}
	if h.left.IsRed() && h.left.left.IsRed() {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *memoryLLRB) rotateLeft(h *memoryNode) *memoryNode {
	x := h.right
	h.right = x.left
	x.left = h
	x.SetColor(h.Color())
	h.SetColor(RED)
	return x
}

func (tree *memoryLLRB) rotateRight(h *memoryNode) *memoryNode {
	x := h.left
	h.left = x.right
	x.right = h
	x.SetColor(h.Color())
	h.SetColor(RED)
	return x
}

func (tree *memoryLLRB) delete(h *memoryNode, key Key) *memoryNode {
	return nil
}

func (tree *memoryLLRB) moveRedLeft(h *memoryNode) *memoryNode {
	h.FlipColors()
	if h.right.left.IsRed() {
		h.right = tree.rotateRight(h.right)
		h = tree.rotateLeft(h)
		h.FlipColors()
	}
	return h
}

func (tree *memoryLLRB) moveRedRight(h *memoryNode) *memoryNode {
	h.FlipColors()
	if h.left.left.IsRed() {
		h = tree.rotateRight(h)
		h.FlipColors()
	}
	return h
}

func (tree *memoryLLRB) fixUp(h *memoryNode) *memoryNode {
	if h.right.IsRed() && !h.left.IsRed() {
		h = tree.rotateLeft(h)
	}
	if h.left.IsRed() && h.left.left.IsRed() {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *memoryLLRB) deleteMin(h *memoryNode) *memoryNode {
	// NOTE this is a check for the sentinel
	if h.left == nil {
		return nil
	}
	if !h.left.IsRed() && !h.left.left.IsRed() {
		h = tree.moveRedLeft(h)
	}
	h.left = tree.deleteMin(h.left)
	return tree.fixUp(h)
}
