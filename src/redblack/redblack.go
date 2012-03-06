package redblack

import "fmt"
import "strings"

import l4g "code.google.com/p/log4go"

var trace l4g.Logger

func init() {
	trace = l4g.NewDefaultLogger(l4g.INFO)
}

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
		Return the number of keys in the tree
	*/
	Size() int
	String() string

	// Internal methods

	rotateLeft(h Node) Node
	rotateRight(h Node) Node
	moveRedLeft(h Node) Node
	moveRedRight(h Node) Node
	fixUp(h Node) Node
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
	return tree.search(tree.Root(), key)
}

func (tree *llrb) search(h Node, key Key) Value {
	// NOTE this is a check for the sentinel
	for h != nil {
		cmp := key.Compare(h.Key())
		if cmp == 0 {
			return h.Value()
		} else if cmp < 0 {
			h = h.Left()
		} else if cmp > 0 {
			h = h.Right()
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
	if tree.Root() != nil {
		tree.Root().SetColor(BLACK)
	}
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
	if isRed(h.Left()) && isRed(h.Right()) {
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
	if isRed(h.Right()) {
		h = tree.rotateLeft(h)
	}
	if isRed(h.Left()) && isRed(h.Left().Left()) {
		h = tree.rotateRight(h)
	}
	return h
}

func (tree *llrb) delete(h Node, key Key) Node {
	if h == nil {
		return h
	}
	trace.Trace("Deleting %v from \n%v", key, h)
	if key.Compare(h.Key()) < 0 {
		if !isRed(h.Left()) && !isRed(h.Left().Left()) {
			h = tree.moveRedLeft(h)
		}
		h.SetLeft(tree.delete(h.Left(), key))
	} else {
		if isRed(h.Left()) {
			h = tree.rotateRight(h)
		}
		if key.Compare(h.Key()) == 0 && h.Right() == nil {
			return nil
		}
		if !isRed(h.Right()) && h.Right() != nil && !isRed(h.Right().Left()) {
			h = tree.moveRedRight(h)
		}
		if key.Compare(h.Key()) == 0 {
			minRight := h.Right().min()
			h.SetValue(tree.search(h.Right(), minRight))
			h.SetKey(minRight)
			h.SetRight(tree.deleteMin(h.Right()))
		} else {
			h.SetRight(tree.delete(h.Right(), key))
		}
	}
	return tree.fixUp(h)
}

func (tree *llrb) deleteMin(h Node) Node {
	trace.Trace("Before deleting min from %v\n%v", h, tree)
	if h.Left() == nil {
		return nil
	}
	if !isRed(h.Left()) && !isRed(h.Left().Left()) {
		h = tree.moveRedLeft(h)
	}
	h.SetLeft(tree.deleteMin(h.Left()))
	trace.Trace("After deleting min from %v\n%v", h, tree)
	return tree.fixUp(h)
}

func (tree *llrb) rotateLeft(h Node) Node {
	trace.Trace("Before rotate left of %v\n%v", h, tree)
	x := h.Right()
	h.SetRight(x.Left())
	x.SetLeft(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	trace.Trace("After rotate left of %v\n%v", x, tree)
	return x
}

func (tree *llrb) rotateRight(h Node) Node {
	trace.Trace("Before rotate right of %v\n%v", h, tree)
	x := h.Left()
	h.SetLeft(x.Right())
	x.SetRight(h)
	x.SetColor(h.Color())
	h.SetColor(RED)
	trace.Trace("After rotate right of %v\n%v", x, tree)
	return x
}

func (tree *llrb) moveRedLeft(h Node) Node {
	trace.Trace("Before move red left of %v\n%v", h, tree)
	h.flipColors()
	if isRed(h.Right().Left()) {
		h.SetRight(tree.rotateRight(h.Right()))
		h = tree.rotateLeft(h)
		h.flipColors()
	}
	trace.Trace("After move red left of %v\n%v", h, tree)
	return h
}

func (tree *llrb) moveRedRight(h Node) Node {
	trace.Trace("Before move red right of %v\n%v", h, tree)
	h.flipColors()
	if isRed(h.Left().Left()) {
		h = tree.rotateRight(h)
		h.flipColors()
	}
	trace.Trace("After move red right of %v\n%v", h, tree)
	return h
}

func (tree *llrb) fixUp(h Node) Node {
	trace.Trace("Before fix up of %v\n%v", h, tree)
	if isRed(h.Right()) {
		h = tree.rotateLeft(h)
	}
	if isRed(h.Left()) && isRed(h.Left().Left()) {
		h = tree.rotateRight(h)
	}
	if isRed(h.Left()) && isRed(h.Right()) {
		h.flipColors()
	}
	trace.Trace("After fix up of %v\n%v", h, tree)
	return h
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
	min() Key
	max() Key
	String() string
}

type NodeImpl interface {
	Key() Key
	SetKey(key Key)
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

func isRed(h Node) bool {
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

func (h *node) min() Key {
	if h.Left() != nil {
		return h.Left().min()
	}
	return h.Key()
}

func (h *node) max() Key {
	if h.Right() != nil {
		return h.Right().max()
	}
	return h.Key()
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
