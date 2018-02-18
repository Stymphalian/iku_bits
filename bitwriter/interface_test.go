package bitwriter

import (
	"bytes"
	"io"
	"log"
	"testing"
)

func createInterfaceForTest(w io.Writer) Interface {
	h, err := NewBitWriter(w)
	if err != nil {
		log.Fatal("Failed to create the bit reader", err)
	}
	return h
}
func TestInterface_TestWriteAndFlush(t *testing.T) {
	for index, tc := range []struct {
		bits     []int
		preFlush []byte
		want     []byte
	}{
		// Write bits within one byte
		{[]int{0}, []byte{}, []byte{0x00}},
		{[]int{1}, []byte{}, []byte{0x80}},
		{[]int{1, 0, 1, 0}, []byte{}, []byte{0xa0}},
		{[]int{1, 0, 1, 0, 1, 1, 1, 1}, []byte{0xaf}, []byte{0xaf}},
		{[]int{
			1, 0, 1, 0, 1, 1, 1, 1,
			0, 0, 1, 1}, []byte{0xaf}, []byte{0xaf, 0x30}},
		{[]int{
			1, 0, 1, 0, 1, 1, 1, 1,
			0, 0, 1, 1, 1, 1, 1, 1,
			0, 0}, []byte{0xaf, 0x3f}, []byte{0xaf, 0x3f, 0x00}},

		// No writes, nothing should show up in flush
		{[]int{}, []byte{}, []byte{}},
	} {

		buf := bytes.NewBufferString("")
		w := createInterfaceForTest(buf)

		for i := 0; i < len(tc.bits); i++ {
			err := w.WriteBit(tc.bits[i])
			if err != nil {
				t.Errorf("Failed[%d]: Failed to write bit %dth bit in sequence %#v",
					index, i, tc.bits)
			}
		}

		if len(tc.bits)%8 != w.Remain() {
			t.Errorf("Failed[%d]: Was expecting %d bits to be remaining but got %d",
				index, len(tc.bits)%8, w.Remain())
		}
		if bytes.Compare(buf.Bytes(), tc.preFlush) != 0 {
			t.Errorf("Failed[%d]: Didnt get expected preflush bytes, got %#v, want %#v",
				index, buf.Bytes(), tc.preFlush)
		}
		w.Flush()
		if bytes.Compare(buf.Bytes(), tc.want) != 0 {
			t.Errorf("Failed[%d]: Didnt get expected final bytes, got %#v, want %#v",
				index, buf.Bytes(), tc.want)
		}
	}
}
