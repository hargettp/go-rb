This package is an implementation of [red-black trees](http://en.wikipedia.org/wiki/Red-black_tree "red-black trees")
in [Go](http://golang.org "Go") using the algorithms described in 
[Left-Leaning Red-Black Trees, R. Sedgwick 2008](http://www.cs.princeton.edu/~rs/talks/LLRB/LLRB.pdf "Left-Leaning Red-Black Trees, R. Sedgwick 2008")

The implementation is intended to be sufficiently generic that any suitable data types can be used
for keys or values.  Although it is not currently implemented, the intent is that eventually persistent
trees will be implemented in an efficient, append-only file format.  Again, this implementation will be
generic so that the algorithms for organizing data for storage can be separated from the actual underlying
storage format.

See the LICENSE file for applicable licensing of this implementation.