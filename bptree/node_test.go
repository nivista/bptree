package bptree

import (
	"reflect"
	"testing"
)

func TestInitNode(t *testing.T) {
	leafFile, fd := newFile()
	initLeaf(fd)
	nodeFile, fd := newFile()
	initNode(fd, []byte(leafFile))
	dumpPage(nodeFile, t)
}

func TestNodeInsert(t *testing.T) {
	leafFile1 := []byte("00000")
	leafFile2 := []byte("11111")
	leafFile3 := []byte("22222")

	nodeFile, fd := newFile()
	node := initNode(fd, leafFile1)
	node.insert([]byte("a"), leafFile2)
	node.insert([]byte("b"), leafFile3)
	dumpPage(nodeFile, t)
}

func TestGetChild(t *testing.T) {
	leafFile1 := []byte("00000")
	leafFile2 := []byte("11111")
	leafFile3 := []byte("22222")

	nodeFile, fd := newFile()
	node := initNode(fd, leafFile1)
	node.insert([]byte("a"), leafFile2)
	node.insert([]byte("c"), leafFile3)

	// equals goes left
	res := node.getChild([]byte("a"))
	if !reflect.DeepEqual(res, []byte("00000")) {
		t.Errorf("expected %v, got %v", []byte("00000"), res)
	}

	res = node.getChild([]byte("b"))
	if !reflect.DeepEqual(res, []byte("11111")) {
		t.Errorf("expected %v, got %v", []byte("11111"), res)
	}

	res = node.getChild([]byte("d"))
	if !reflect.DeepEqual(res, []byte("22222")) {
		t.Errorf("expected %v, got %v", []byte("22222"), res)
	}
	dumpPage(nodeFile, t)

}
