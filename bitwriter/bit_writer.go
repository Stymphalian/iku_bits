package bitwriter

import (
	"fmt"
	"io"
)

type BitWriter struct {
	w       io.Writer
	buf     byte
	numBits uint
}

func NewBitWriter(w io.Writer) (*BitWriter, error) {
	return &BitWriter{w, 0, 0}, nil
}

func (this *BitWriter) WriteBit(v int) error {
	if v == 1 {
		this.buf |= (1 << (uint(7) - this.numBits))
	}
	this.numBits += 1

	if this.numBits == 8 {
		err := this.writeByteToStream()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *BitWriter) Flush() (int, error) {
	n := this.numBits
	if n > 0 {
		err := this.writeByteToStream()
		if err != nil {
			return 0, err
		}
	}
	return int(n), nil
}

func (this *BitWriter) Remain() int {
	return int(this.numBits)
}

func (this *BitWriter) writeByteToStream() error {
	n, err := this.w.Write([]byte{this.buf})
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("Tried to write byte to stream but failed")
	}
	this.buf = 0
	this.numBits = 0
	return nil
}
