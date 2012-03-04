package redblack

import "fmt"
import "strings"

const (
	RED   = Color(true)
	BLACK = Color(false)
)

/*
Representation of node colors
*/
type Color bool

func (c Color) String() string {
	if c == RED {
		return "RED"
	}
	return "BLACK"
}

//=============================================================================
//
// Left-leaning red-black trees
//
//=============================================================================

/*
General red-black tree interface, and a left-leaning red-black tree specifically.  
While this interface can be extended most custom implementations will want to 
implement LLRBImpl instead
*/
type LLRB interface {
	LLRBImpl
	/*
		Create a new node for use with this tree
	*/
	NewNode(key Key, value Value) Node
	/*
		Search for the value associated with the provided key
	*/
	Search(key Key) Value
	/*
		Insert a new key into the tree with the associated value
	*/
	Insert(key Key, value Value)
	/*
		Delete the indicated key and its corresponding value from the tree
	*/
	Delete(key Key)
	/*
		Delete from the tree the value with the minimum key
	*/
	DeleteMin()
	/*
		Return the number of keys in the tree
	*/
	Size() int
	String() string
}

/*
Implement this interface to customize red-black tree behavior
*/
type LLRBImpl interface {
	/*
		Create an instance of the node implementation for the provided key and value;
		will be used to create a new node internally within the tree
	*/
	NewNodeImpl(key Key, value Value) NodeImpl
	/*
		Return the root of the tree
	*/
	Root() Node
	/*
		Change the root of the tree.
	*/
	SetRoot(h Node)
}

type llrb struct {
	LLRBImpl
}

/*
Create a new red-black tree using the provided implementation
*/
func NewRedBlackTree(impl LLRBImpl) LLRB {
	return &llrb{impl}
}

//
// Default implementations of LLRB methods
//

func (tree *llrb) NewNode(key Key, value Value) Node {
	// return &memoryNode{key: key, value: value, color: RED}
	return &node{tree.NewNodeImpl(key, value)}
}

func (tree *llrb) Search(key Key) Value {
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

func (tree *llrb) Insert(key Key, value Value) {
	tree.SetRoot(tree.insert(tree.Root(), key, value))
	tree.Root().SetColor(BLACK)
}

func (tree *llrb) Delete(key Key) {
	tree.SetRoot(tree.delete(tree.Root(), key))
	tree.Root().SetColor(BLACK)
}

func (tree *llrb) DeleteMin() {
	tree.SetRoot(tree.deleteMin(tree.Root()))
	tree.Root().SetColor(BLACK)
}

func (tree *llrb) Size() int {
	var count func(h Node) int
	count = func(h Node) int {
		if h != nil {
			return 1 + count(h.Left()) + count(h.Right())
		}
		return 0
	}
	return count(tree.Root())
}

func (tree *llrb) String() string {
	var visit func(h Node, depth int) string
	visit = func(h Node, depth int) string {
		if h == nil {
			return ""
		}
		s := fmt.Sprintf("%v%v\n", strings.Repeat("\t", depth), h)
		s += visit(h.Left(), depth+1)
		s += visit(h.Right(), depth+1)
		return s
	}
	return visit(tree.Root(), 0)
}

// LLRB implementation

func (tree *llrb) insert(h Node, key Key, value Value) Node {
	// NOTE this is a check for the sentinel
	if h == nil {
		return tree.NewNode(key, value)
	}
	if h.Left() != nil && h.Left().isRed() && h.Right() != nil && h.Right().isRed() {
		h.flipColors()
	}
	cmp := key.Compare(h.Key())
	if cmp == 0 {
		h.SetValue(value)
	} else if cmp < 0 {
		h.SetLeft(tree.insert(h.Left(), key, value))
	} else if cmp > 0 {
		h.SetRight(tree.insert(h.Right(), key, value))
	}
	if h.Right() != nil && h.Right().isRed() && (h.Left() == nil || !h.Left().isRed()) {
		h = tree.rotateLeft(h)
	}
	if h.Left() != nil && h.Left().isRed() && h.Left().Left() != nil && h.Left().Left().isRed() {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *llrb) rotateLeft(h Node) Node {
	x := h.Right()
	h.SetRight(x.Left())
	x.SetLeft(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	return x
}

func (tree *llrb) rotateRight(h Node) Node {
	x := h.Left()
	h.SetLeft(x.Right())
	x.SetRight(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	return x
}

func (tree *llrb) delete(h Node, key Key) Node {
	return nil
}

func (tree *llrb) moveRedLeft(h Node) Node {
	h.flipColors()
	if h.Right().Left().isRed() {
		h.SetRight(tree.rotateRight(h.Right()))
		h = tree.rotateLeft(h)
		h.flipColors()
	}
	return h
}

func (tree *llrb) moveRedRight(h Node) Node {
	h.flipColors()
	if h.Left().Left().isRed() {
		h = tree.rotateRight(h)
		h.flipColors()
	}
	return h
}

func (tree *llrb) fixUp(h Node) Node {
	if h.Right().isRed() && !h.Left().isRed() {
		h = tree.rotateLeft(h)
	}
	if h.Left().isRed() && h.Left().Left().isRed() {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *llrb) deleteMin(h Node) Node {
	// NOTE this is a check for the sentinel
	if h.Left() == nil {
		return nil
	}
	if !h.Left().isRed() && !h.Left().Left().isRed() {
		h = tree.moveRedLeft(h)
	}
	h.SetLeft(tree.deleteMin(h.Left()))
	return tree.fixUp(h)
}

//=============================================================================
//
// Keys
//
//=============================================================================

/*
Generalized interface for keys used in red-black trees
*/
type Key interface {
	/*
		Compare 2 keys, and return 0 if they are equivalent, -1 if key1
		is less than key 2, and 1 if key1 is greater than key2
	*/
	Compare(other Key) int
	String() string
}

type IntKey int

func (key1 IntKey) Compare(key2 Key) int {
	int1 := int(key1)
	int2 := int(key2.(IntKey))
	switch {
	case int1 == int2:
		return 0
	case int1 < int2:
		return -1
	case int1 > int2:
		return 1
	}
	return 1
}

func (key IntKey) String() string {
	return fmt.Sprintf("%v", int(key))
}

//=============================================================================
//
// Values
//
//=============================================================================

/*
Generalized interface for values stored in a red-black tree
*/
type Value interface {
	String() string
}

type StringValue string

func (value StringValue) String() string {
	return string(value)
}

type BytesValue []byte

func (value BytesValue) String() string {
	return string(([]byte)(value))
}

//=============================================================================
//
// Nodes
//
//=============================================================================

/*
Generalized interface for red-black tree nodes
*/
type Node interface {
	NodeImpl
	flipColors()
	isRed() bool
	String() string
}

type NodeImpl interface {
	Key() Key
	Value() Value
	SetValue(value Value)
	Left() Node
	SetLeft(h Node)
	Right() Node
	SetRight(h Node)
	Color() Color
	SetColor(c Color)
}

type node struct {
	NodeImpl
}

func (h *node) isRed() bool {
	if h == nil {
		return false
	}
	return h.Color() == RED
}

func (h *node) flipColors() {
	h.SetColor(!h.Color())
	h.Left().SetColor(!h.Left().Color())
	h.Right().SetColor(!h.Right().Color())
}

func (h *node) String() string {
	var leftKey, rightKey Key
	if h.Left() != nil {
		leftKey = h.Left().Key()
	} else {
		leftKey = nil
	}
	if h.Right() != nil {
		rightKey = h.Right().Key()
	} else {
		rightKey = nil
	}
	return fmt.Sprintf("key=%v,left={%v},right={%v},color=%v,value=%v",
		h.Key(), leftKey, rightKey, h.Color(), h.Value())
}
