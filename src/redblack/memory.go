package redblack

type memoryLLRB struct {
	root Node
}

type memoryNode struct {
	key         Key
	value       Value
	left, right Node
	color       Color
}

// LLRB implementation

func NewMemoryLLRB() LLRB {
	return NewRedBlackTree(&memoryLLRB{root: nil})
}

func (tree *memoryLLRB) NewNodeImpl(key Key, value Value) NodeImpl {
	return &memoryNode{key: key, value: value, color: RED}
}

func (tree *memoryLLRB) Root() Node {
	return tree.root
}

func (tree *memoryLLRB) SetRoot(root Node) {
	tree.root = root
}

// Node implementation

func (h *memoryNode) Key() Key {
	return h.key
}

func (h *memoryNode) SetKey(key Key) {
	h.key = key
}

func (h *memoryNode) Value() Value {
	return h.value
}

func (h *memoryNode) SetValue(value Value) {
	h.value = value
}

func (h *memoryNode) Left() Node {
	if h == nil {
		return nil
	}
	return h.left
}

func (n *memoryNode) SetLeft(h Node) {
	n.left = h
}

func (h *memoryNode) Right() Node {
	if h == nil {
		return nil
	}
	return h.right
}

func (n *memoryNode) SetRight(h Node) {
	n.right = h
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
