package bitreader

import (
	"bytes"
	"io"
	"log"
	"reflect"
	"testing"
)

func createInterfaceForTest(r io.Reader) Interface {
	b, err := NewBitReader(r)
	if err != nil {
		log.Fatal("Failed to create the bit reader", err)
	}
	return b
}
func TestInterface_TestReader(t *testing.T) {
	for index, tc := range []struct {
		buf        []byte
		bitsToRead int
		want       []int
	}{
		// Read bits only from one byte
		{[]byte{0x55}, 1, []int{0}},
		{[]byte{0x80}, 1, []int{1}},
		{[]byte{0x55}, 4, []int{0, 1, 0, 1}},
		{[]byte{0x5c}, 8, []int{0, 1, 0, 1, 1, 1, 0, 0}},

		// Read bits across one byte boundary
		{[]byte{0x5c, 0x55}, 10, []int{0, 1, 0, 1, 1, 1, 0, 0, 0, 1}},
		{[]byte{0x5c, 0x55}, 12, []int{0, 1, 0, 1, 1, 1, 0, 0, 0, 1, 0, 1}},
	} {
		buf := bytes.NewBuffer(tc.buf)
		r := createInterfaceForTest(buf)
		got := make([]int, 0)
		for i := 0; i < tc.bitsToRead; i++ {
			v, err := r.ReadBit()
			if err != nil {
				t.Errorf("Failed[%d]: %v", index, err)
			}
			got = append(got, v)
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Failed[%d]: got %v, but want %v", index, got, tc.want)
		}
	}
}

func TestInterface_ReadPastBuf(t *testing.T) {
	for index, tc := range []struct {
		buf        []byte
		bitsToRead int
	}{
		// Read bits only from one byte
		{[]byte{0x80}, 9},
		// Read bits across one byte boundary
		{[]byte{0x5c, 0x55}, 17},
	} {
		buf := bytes.NewBuffer(tc.buf)
		r := createInterfaceForTest(buf)
		for i := 0; i < tc.bitsToRead-1; i++ {
			_, err := r.ReadBit()
			if err != nil {
				t.Errorf("Failed[%d]: %v", index, err)
			}
		}
		_, err := r.ReadBit()
		if err == nil {
			t.Errorf("Failed[%d]: Trying to read out of byte boundary should fail.",
				index)
		}
	}
}

func TestInterface_TestPeeker(t *testing.T) {
	for index, tc := range []struct {
		buf          []byte
		bitsToRead   int
		expectedPeek int
	}{
		// Read bits only from one byte
		{[]byte{0x55}, 0, 0},
		{[]byte{0x55}, 1, 1},
		{[]byte{0x80}, 1, 0},
		{[]byte{0x55}, 4, 0},
		{[]byte{0x57}, 7, 1},
		// {[]byte{0x5c}, 8, 1},

		// Read bits across one byte boundary
		{[]byte{0x5c, 0x85}, 8, 1},
		{[]byte{0x5c, 0x55}, 10, 0},
		{[]byte{0x5c, 0x55}, 13, 1},
		{[]byte{0x5c, 0x57}, 15, 1},
	} {
		buf := bytes.NewBuffer(tc.buf)
		r := createInterfaceForTest(buf)
		for i := 0; i < tc.bitsToRead; i++ {
			_, err := r.ReadBit()
			if err != nil {
				t.Errorf("Failed[%d]: %v", index, err)
			}
		}
		v, err := r.Peek()
		if err != nil {
			t.Errorf("Failed[%d]: %v", index, err)
		}

		if v != tc.expectedPeek {
			t.Errorf("Failed[%d]: expected to get %d when peeking but got %d",
				index, tc.expectedPeek, v)
		}
	}
}

func TestInterface_TestPeekPastBuffer(t *testing.T) {
	for index, tc := range []struct {
		buf        []byte
		bitsToRead int
	}{
		// Read bits only from one byte
		{[]byte{0x55}, 8},

		// Read bits across one byte boundary
		{[]byte{0x5c, 0x85}, 16},
	} {
		buf := bytes.NewBuffer(tc.buf)
		r := createInterfaceForTest(buf)
		for i := 0; i < tc.bitsToRead; i++ {
			_, err := r.ReadBit()
			if err != nil {
				t.Errorf("Failed[%d]: %v", index, err)
			}
		}
		_, err := r.Peek()
		if err == nil {
			t.Errorf("Failed[%d]: Peek should fail if there is not more data in"+
				" buffer", index)
		}
	}
}
