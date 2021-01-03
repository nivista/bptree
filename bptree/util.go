package bptree

import (
	"fmt"
	"math/rand"

	"golang.org/x/sys/unix"
)

const fileCreateRetries = 10

func intToTwoBytes(n int) []byte {
	return []byte{
		byte((n & 0xff00) >> 8),
		byte((n & 0xff)),
	}
}

func twoBytesToInt(data []byte) int {
	if len(data) != 2 {
		panic("expected 2 bytes")
	}
	return (int(data[0]) << 8) | int(data[1])
}

// returns 1 if a is greater than b, 0, if they're equal, -1 if a is less
func compareBytes(a, b []byte) int {
	// skip zeroes
	for len(a) != 0 && a[0] == 0 {
		a = a[1:]
	}
	for len(b) != 0 && b[0] == 0 {
		b = b[1:]
	}

	if len(a) > len(b) {
		return 1
	}
	if len(a) < len(b) {
		return -1
	}
	for idx := range a {
		if a[idx] > b[idx] {
			return 1
		}
		if a[idx] < b[idx] {
			return -1
		}
	}
	return 0
}

func newFile() (filename string, fd int) {
	var err error
	for i := 0; i < fileCreateRetries; i++ {
		filename = fmt.Sprintf("%05d", rand.Uint32()%10000)
		fd, err = createFile(filename)
		if err == nil {
			break
		}
	}
	if err != nil {
		panic(err)
	}
	return filename, fd
}

// returns leaf, node or error
func getNodeR(filename string) interface{} {
	fd, err := getFileR(filename)
	if err != nil {
		return err
	}
	return getNode(fd)
}

func getNodeRW(filename string) interface{} {
	fd, err := getFileRW(filename)
	if err != nil {
		return err
	}
	return getNode(fd)
}

func getNode(fd int) interface{} {
	data, err := unix.Mmap(fd, 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer unix.Munmap(data)

	switch data[0] {
	case 0: // node
		return &node{
			fd:             fd,
			n:              twoBytesToInt(data[1:3]),
			leftmostOffset: twoBytesToInt(data[3:5]),
			availableSpace: twoBytesToInt(data[5:7]),
			leftPointer:    data[7:12],
		}
	case 1: // leaf
		return &leaf{
			fd:             fd,
			n:              twoBytesToInt(data[1:3]),
			leftmostOffset: twoBytesToInt(data[3:5]),
			availableSpace: twoBytesToInt(data[5:7]),
		}
	default:
		panic("unknown node type")
	}

}
