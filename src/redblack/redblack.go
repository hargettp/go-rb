package redblack

const (
	RED   = true
	BLACK = false
)

type Color bool

type LLRB interface {
	Search(key Key) Value
	Insert(key Key, value Value)
	Delete(key Key)
	DeleteMin()
}

type node interface {
	getKey() Key
	getValue() Value
	isRed() bool
	leftChild() node
	rightChild() node
	flipColors()
}

type Key interface {
	CompareTo(other Key) int
}

type Value interface {
}
