package bptree

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestCreateFile(t *testing.T) {
	createFile("hello")
}

func TestNewFile(t *testing.T) {
	name, _ := newFile()
	t.Log(name)
	dumpPage(name, t)
}

func TestInitLeaf(t *testing.T) {
	name, fd := newFile()
	t.Log(name)
	initLeaf(fd)
	dumpPage(name, t)
}

func TestInsert(t *testing.T) {
	name, fd := newFile()
	t.Log(name)
	l := initLeaf(fd)
	err := l.insert([]byte("aaaaa"), []byte("bbbbbbbbbb"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.insert([]byte("cc"), []byte("d"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.insert([]byte("aa"), []byte("bb"))
	if err != nil {
		t.Fatal(err)
	}
	dumpPage(name, t)
}

func TestGet(t *testing.T) {
	name, fd := newFile()
	t.Log(name)
	l := initLeaf(fd)
	err := l.insert([]byte("aaaaa"), []byte("bbbbbbbbbb"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.insert([]byte("cc"), []byte("d"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.insert([]byte("aa"), []byte("bb"))
	if err != nil {
		t.Fatal(err)
	}
	val, err := l.get([]byte("cc"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
	if !reflect.DeepEqual(val, []byte("d")) {
		t.Fatalf("expected %v, got %v\n", val, []byte("d"))
	}
	dumpPage(name, t)
}

func TestDelete(t *testing.T) {
	name, fd := newFile()
	t.Log(name)
	l := initLeaf(fd)
	err := l.insert([]byte("aaaaa"), []byte("bbbbbbbbbb"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.insert([]byte("cc"), []byte("d"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.insert([]byte("aa"), []byte("bb"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.delete([]byte("cc"))
	if err != nil {
		t.Fatal(err)
	}
	err = l.delete([]byte("cc"))
	if err == nil {
		t.Fatal("expected second delete to fail")
	}

	dumpPage(name, t)
}

func dumpPage(filename string, t *testing.T) {
	file, err := os.Open(path.Join(BASE_DIR, filename))
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(data)
}
