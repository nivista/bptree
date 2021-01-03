package bptree

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	tree := New()
	err := tree.Insert([]byte("hello"), []byte("world"))
	if err != nil {
		t.Fatal(err)
	}

	err = tree.Insert([]byte("hey"), []byte("universe"))
	if err != nil {
		t.Fatal(err)
	}

	contains, err := tree.Contains([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

	if !contains {
		t.Fatal("expected to contain hello")
	}

	val, err := tree.Get([]byte("hey"))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, []byte("universe")) {
		t.Fatalf("expected universe, got %v", string(val))
	}

	err = tree.Delete([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}

}
