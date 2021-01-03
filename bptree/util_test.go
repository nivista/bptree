package bptree

import (
	"reflect"
	"testing"
)

func TestIntToTwoBytes(t *testing.T) {
	cases := []struct {
		n int
		b []byte
	}{
		{1, []byte{0, 1}},
		{256, []byte{1, 0}},
		{257, []byte{1, 1}},
	}

	for _, testCase := range cases {
		res := intToTwoBytes(testCase.n)
		if !reflect.DeepEqual(testCase.b, res) {
			t.Errorf("expected %v, got %v\n", testCase.b, res)
		}
	}
}

func TestTwoBytesToInt(t *testing.T) {
	cases := []struct {
		n int
		b []byte
	}{
		{1, []byte{0, 1}},
		{256, []byte{1, 0}},
		{257, []byte{1, 1}},
	}

	for _, testCase := range cases {
		res := twoBytesToInt(testCase.b)
		if testCase.n != res {
			t.Errorf("expected %v, got %v\n", testCase.n, res)
		}
	}
}

func TestCompareBytes(t *testing.T) {
	cases := []struct {
		a, b []byte
		out  int
	}{
		{[]byte{'a'}, []byte{'a', 'b'}, -1},
		{[]byte{'a', 'b'}, []byte{'b', 'b'}, -1},
		{[]byte{'a'}, []byte{'a'}, 0},
	}

	for _, testCase := range cases {
		res := compareBytes(testCase.a, testCase.b)
		if res != testCase.out {
			t.Errorf("expected %v, got %v", testCase.out, res)
		}
	}
}
