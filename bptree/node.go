package bptree

import (
	"fmt"

	"golang.org/x/sys/unix"
)

const nodeHeaderLen = 12

type node struct {
	fd             int
	n              int
	leftmostOffset int
	availableSpace int
	leftPointer    []byte
}

func initNode(fd int, leftPointer []byte) node {
	if len(leftPointer) != 5 {
		panic("len of leftpointer should be 5")
	}

	data, err := unix.Mmap(fd, 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	// specify type as node
	data[0] = 0

	// specify number of records
	copy(data[1:], intToTwoBytes(0))

	// specify leftmostOffset
	copy(data[3:], intToTwoBytes(pageSize))

	// specify available space
	copy(data[5:], intToTwoBytes(pageSize-nodeHeaderLen))

	// specify left pointer
	copy(data[7:], leftPointer)

	// unmap data
	unix.Munmap(data)

	return node{fd, 0, pageSize, pageSize - nodeHeaderLen, leftPointer}
}

// returns node, leaf, or error
func (n *node) getChild(key []byte) []byte {
	data, err := unix.Mmap(n.fd, 0, pageSize, unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer unix.Munmap(data)

	currPointer := n.leftPointer
	for offsetIdx := nodeHeaderLen; offsetIdx < nodeHeaderLen+2*n.n; offsetIdx += 2 {
		recordIdx := twoBytesToInt(data[offsetIdx : offsetIdx+2])
		recordKeyLen := twoBytesToInt(data[recordIdx : recordIdx+2])
		recordKey := data[recordIdx+2 : recordIdx+2+recordKeyLen]

		if compareBytes(recordKey, key) >= 0 {
			break
		} else {
			copy(currPointer, data[recordIdx+2+recordKeyLen:recordIdx+2+recordKeyLen+5])
		}
	}

	return currPointer
}

// NOTE: when there's a split, the left page gets rewritten, the right page gets created
// therefore this takes the fd of the fh that it should point to on the right
func (n *node) insert(key []byte, filename []byte) error {
	if len(filename) != 5 {
		panic("bad filename length insert")
	}
	if 2+2+len(key)+5 > n.availableSpace {
		return fmt.Errorf("not enough space")
	}

	data, err := unix.Mmap(n.fd, 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer unix.Munmap(data)

	insertionIdx := nodeHeaderLen

	// scan for insertion point
	for offsetIdx := nodeHeaderLen; offsetIdx < nodeHeaderLen+2*n.n; offsetIdx += 2 {
		recordIdx := twoBytesToInt(data[offsetIdx : offsetIdx+2])
		recordKeyLen := twoBytesToInt(data[recordIdx : recordIdx+2])
		recordKey := data[recordIdx+2 : recordIdx+2+recordKeyLen]

		compare := compareBytes(key, recordKey)
		if compare == 1 {
			insertionIdx = offsetIdx + 2
			break
		}
		if compare == 0 {
			return fmt.Errorf("key already exists")
		}
	}

	// move offsets out of way
	copy(data[insertionIdx+2:], data[insertionIdx:nodeHeaderLen+n.n*2])

	//update metadata
	n.leftmostOffset -= len(key) + 7
	copy(data[3:], intToTwoBytes(n.leftmostOffset))

	n.availableSpace -= len(key) + 9
	copy(data[5:], intToTwoBytes(n.availableSpace))

	n.n++
	copy(data[1:], intToTwoBytes(n.n))

	// copy data
	copy(data[insertionIdx:], intToTwoBytes(n.leftmostOffset))
	copy(data[n.leftmostOffset:], intToTwoBytes(len(key)))
	copy(data[n.leftmostOffset+2:], key)
	copy(data[n.leftmostOffset+2+len(key):], filename)

	return nil
}
