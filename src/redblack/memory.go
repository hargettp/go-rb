package redblack

type memoryLLRB struct {
	root *memoryNode
}

type memoryNode struct {
	key         Key
	value       Value
	left, right *memoryNode
	color       Color
}

// node implementation

func (h *memoryNode) getKey() Key {
	return h.key
}

func (h *memoryNode) getValue() Value {
	return h.value
}

func (h *memoryNode) isRed() bool {
	return h.color == RED
}

func (h *memoryNode) leftChild() node {
	return h.left
}

func (h *memoryNode) rightChild() node {
	return h.right
}

func (h *memoryNode) flipColors() {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

// LLRB implementation

func NewMemoryLLRB() LLRB {
	return &memoryLLRB{root: nil}
}

func (tree *memoryLLRB) Search(key Key) Value {
	x := tree.root
	// NOTE this is a check for the sentinel
	for x != nil {
		cmp := key.CompareTo(x.key)
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
}

func (tree *memoryLLRB) Delete(key Key) {
	tree.root = tree.delete(tree.root, key)
	tree.root.color = BLACK
}

func (tree *memoryLLRB) DeleteMin() {
	tree.root = tree.deleteMin(tree.root)
	tree.root.color = BLACK
}

// LLRB internals

func (tree *memoryLLRB) insert(h *memoryNode, key Key, value Value) *memoryNode {
	// NOTE this is a check for the sentinel
	if h == nil {
		return &memoryNode{key: key, value: value, color: RED}
	}
	if h.left.isRed() && h.right.isRed() {
		h.flipColors()
	}
	cmp := key.CompareTo(h.key)
	if cmp == 0 {
		h.value = value
	} else if cmp < 0 {
		h.left = tree.insert(h.left, key, value)
	} else if cmp > 0 {
		h.right = tree.insert(h.right, key, value)
	}
	if h.right.isRed() && !h.left.isRed() {
		h = tree.rotateLeft(h)
	}
	if h.left.isRed() && h.left.left.isRed() {
		h = tree.rotateRight(h)
	}
	return nil
}

func (tree *memoryLLRB) rotateLeft(h *memoryNode) *memoryNode {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = RED
	return x
}

func (tree *memoryLLRB) rotateRight(h *memoryNode) *memoryNode {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = RED
	return x
}

func (tree *memoryLLRB) delete(h *memoryNode, key Key) *memoryNode {
	return nil
}

func (tree *memoryLLRB) moveRedLeft(h *memoryNode) *memoryNode {
	h.flipColors()
	if h.right.left.isRed() {
		h.right = tree.rotateRight(h.right)
		h = tree.rotateLeft(h)
		h.flipColors()
	}
	return h
}

func (tree *memoryLLRB) moveRedRight(h *memoryNode) *memoryNode {
	h.flipColors()
	if h.left.left.isRed() {
		h = tree.rotateRight(h)
		h.flipColors()
	}
	return h
}

func (tree *memoryLLRB) fixUp(h *memoryNode) *memoryNode {
	if h.right.isRed() && !h.left.isRed() {
		h = tree.rotateLeft(h)
	}
	if h.left.isRed() && h.left.left.isRed() {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *memoryLLRB) deleteMin(h *memoryNode) *memoryNode {
	// NOTE this is a check for the sentinel
	if h.left == nil {
		return nil
	}
	if !h.left.isRed() && !h.left.left.isRed() {
		h = tree.moveRedLeft(h)
	}
	h.left = tree.deleteMin(h.left)
	return tree.fixUp(h)
}
