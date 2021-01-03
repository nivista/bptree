package bptree

const BASE_DIR = "/home/yaniv/go/src/github.com/nivista/bptree/data"

// Tree is a B+-tree
type Tree struct {
	root *node
}

// New returns a new Tree
func New() Tree {
	leafFile, fd := newFile()
	initLeaf(fd)
	_, fd = newFile()
	node := initNode(fd, []byte(leafFile))
	return Tree{&node}
}

// Contains checks if the tree contains an element
func (t Tree) Contains(key []byte) (bool, error) {
	curr := t.root
	for {
		filename := string(curr.getChild(key))
		child := getNodeR(filename)
		switch child := child.(type) {
		case *node:
			curr = child
		case *leaf:
			_, err := child.get(key)
			return err == nil, nil
		case error:
			panic(child)
		default:
			panic("unknown type")
		}
	}
}

func (t Tree) Get(key []byte) ([]byte, error) {
	curr := t.root
	for {
		filename := string(curr.getChild(key))
		child := getNodeR(filename)
		switch child := child.(type) {
		case *node:
			curr = child
		case *leaf:
			return child.get(key)
		case error:
			panic(child)
		default:
			panic("unknown type")
		}
	}
}

// Insert inserts an element
func (t Tree) Insert(key, val []byte) error {
	curr := t.root
	for {
		filename := string(curr.getChild(key))
		child := getNodeRW(filename)
		switch child := child.(type) {
		case *node:
			curr = child
		case *leaf:
			return child.insert(key, val)
		case error:
			panic(child)
		default:
			panic("unknown type")
		}
	}
}

// Delete deletes an element
func (t *Tree) Delete(key []byte) error {
	curr := t.root
	for {
		filename := string(curr.getChild(key))
		child := getNodeRW(filename)
		switch child := child.(type) {
		case *node:
			curr = child
		case *leaf:
			return child.delete(key)
		case error:
			panic(child)
		default:
			panic("unknown type")
		}
	}
}

// b1 is a 0 if its a node, a 1 if its a leaf

// file structure for node
// b2-3 will indicate number of kv pairs (N)
// b4-5 will indicate leftmostOffset
// b6-7 will indicate available space
// b8-12 will be the left pointer
// b10-(10+N*2) will be record offsets
// records will go from left to right
// the first two bytes will be key size, then the key, then the value will be FIVE bytes
// NOTE: filenames will be 5 digit numbers

// file structure for leaf
// b2-3 will indicate number of kv pairs (N)
// b4-5 will indicate leftmostOffset
// b6-7 will indicate available space
// b8-(8+N*2) will be record offsets
// records will go from left to right
// the first two bytes will be key size , the next two will be value size, then the key, then the value

// I can tell whether something will be a node or a leaf based on where we are in the tree
// this might have to change when we do concurrency, what if the height changes while i'm doing something?
