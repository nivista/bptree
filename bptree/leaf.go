package bptree

import (
	"errors"
	"fmt"

	"golang.org/x/sys/unix"
)

const leafHeaderLen = 7

type leaf struct {
	fd             int // file descriptor
	n              int // number of records
	leftmostOffset int // offset of leftmost record
	availableSpace int // available space
}

func initLeaf(fd int) leaf {
	data, err := unix.Mmap(fd, 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	// specify type as leaf
	data[0] = 1

	// specify number of records
	copy(data[1:], intToTwoBytes(0))

	// specify leftMostOffset
	copy(data[3:], intToTwoBytes(pageSize))

	// specify available space
	copy(data[5:], intToTwoBytes(pageSize-leafHeaderLen))

	// unmap data
	unix.Munmap(data)

	return leaf{fd, 0, pageSize, pageSize - leafHeaderLen}
}

func (l *leaf) insert(key []byte, val []byte) error {
	keyLen := len(key)
	valLen := len(val)

	if l.availableSpace < keyLen+valLen {
		return fmt.Errorf("not enough space")
	}

	data, err := unix.Mmap(l.fd, 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	insertionIdx := leafHeaderLen

	// scan for insertion point
	for offsetIdx := leafHeaderLen; offsetIdx < leafHeaderLen+2*l.n; offsetIdx += 2 {
		recordIdx := twoBytesToInt(data[offsetIdx : offsetIdx+2])
		recordKeyLen := twoBytesToInt(data[recordIdx : recordIdx+2])
		recordKey := data[recordIdx+4 : recordIdx+4+recordKeyLen]

		compare := compareBytes(key, recordKey)
		if compare == 1 {
			insertionIdx = offsetIdx
			break
		}
		if compare == 0 {
			return fmt.Errorf("key already exists")
		}
	}

	// move offsets out of way
	copy(data[insertionIdx+2:], data[insertionIdx:leafHeaderLen+l.n*2])

	//update metadata
	l.leftmostOffset -= keyLen + valLen + 4
	copy(data[3:], intToTwoBytes(l.leftmostOffset))

	l.availableSpace -= keyLen + valLen + 6
	copy(data[5:], intToTwoBytes(l.availableSpace))

	l.n++
	copy(data[1:], intToTwoBytes(l.n))

	// copy data
	copy(data[insertionIdx:], intToTwoBytes(l.leftmostOffset))
	copy(data[l.leftmostOffset:], intToTwoBytes(keyLen))
	copy(data[l.leftmostOffset+2:], intToTwoBytes(valLen))
	copy(data[l.leftmostOffset+4:], key)
	copy(data[l.leftmostOffset+4+keyLen:], val)

	return nil
}

func (l *leaf) delete(key []byte) error {
	data, err := unix.Mmap(l.fd, 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	var deletionIdx int
	var recordKeyLen int
	var recordValLen int
	for offsetIdx := leafHeaderLen; offsetIdx < leafHeaderLen+2*l.n; offsetIdx += 2 {
		recordIdx := twoBytesToInt(data[offsetIdx : offsetIdx+2])
		recordKeyLen = twoBytesToInt(data[recordIdx : recordIdx+2])
		recordKey := data[recordIdx+4 : recordIdx+4+recordKeyLen]
		recordValLen = twoBytesToInt(data[recordIdx+2 : recordIdx+4])
		if compareBytes(key, recordKey) == 0 {
			deletionIdx = offsetIdx
			break
		}
	}

	if deletionIdx == 0 {
		return fmt.Errorf("couldn't find key")
	}

	copy(data[deletionIdx:], data[deletionIdx+2:leafHeaderLen+l.n*2])

	//update metadata
	l.availableSpace += recordKeyLen + recordValLen + 6
	copy(data[5:], intToTwoBytes(l.availableSpace))

	l.n--
	copy(data[1:], intToTwoBytes(l.n))

	unix.Munmap(data)
	return nil
}

func (l *leaf) get(key []byte) (val []byte, err error) {
	data, err := unix.Mmap(l.fd, 0, pageSize, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		panic(err)
	}
	defer unix.Munmap(data)

	for offsetIdx := leafHeaderLen; offsetIdx < leafHeaderLen+2*l.n; offsetIdx += 2 {
		recordIdx := twoBytesToInt(data[offsetIdx : offsetIdx+2])
		recordKeyLen := twoBytesToInt(data[recordIdx : recordIdx+2])
		recordKey := data[recordIdx+4 : recordIdx+4+recordKeyLen]

		if compareBytes(key, recordKey) == 0 {
			recordValLen := twoBytesToInt(data[recordIdx+2 : recordIdx+4])

			recordVal := data[recordIdx+4+recordKeyLen : recordIdx+4+recordKeyLen+recordValLen]
			out := make([]byte, recordValLen)
			copy(out, recordVal)
			return out, nil
		}
	}
	return nil, errors.New("not found")
}
