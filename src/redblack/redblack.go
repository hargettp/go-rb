package redblack

import "fmt"

const (
	RED   = Color(true)
	BLACK = Color(false)
)

type Color bool

func (c Color) String() string {
	if c == RED {
		return "RED"
	}
	return "BLACK"
}

/*
General red-black tree interface
*/
type LLRB interface {
	Search(key Key) Value
	Insert(key Key, value Value)
	Delete(key Key)
	DeleteMin()
	Size() int
	String() string
}

/*
Generalized interface for keys used in red-black trees
*/
type Key interface {
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

/*
Internal node interface
*/
type Node interface {
	Key() Key
	Value() Value
	IsRed() bool
	Left() Node
	Right() Node
	FlipColors()
	Color() Color
	SetColor(c Color)
	String() string
}
