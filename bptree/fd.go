package bptree

import (
	"path"

	"golang.org/x/sys/unix"
)

// page size
const pageSize = 4096

var zeroes = make([]byte, pageSize)

func getFileR(filename string) (int, error) {
	filepath := path.Join(BASE_DIR, filename)
	return unix.Open(filepath, unix.O_RDWR, 0) // O_RDONLY is failiing!!
}

func getFileRW(filename string) (int, error) {
	filepath := path.Join(BASE_DIR, filename)
	return unix.Open(filepath, unix.O_RDWR, 0)
}

func createFile(filename string) (int, error) {
	filepath := path.Join(BASE_DIR, filename)
	fd, err := unix.Open(filepath, unix.O_RDWR|unix.O_EXCL|unix.O_CREAT, unix.S_IRWXU)
	if err != nil {
		return 0, err
	}
	unix.Write(fd, zeroes)
	return fd, nil
}

func releaseFD(fd int) error {
	return unix.Close(fd)
}
