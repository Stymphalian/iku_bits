package bitreader

import (
	"fmt"
	"io"
)

// Type used for reading a single bits out of a Byte stream
type BitReader struct {
	r io.Reader

	storedByte byte
	numBits    uint
}

// Return a new bit reader, reading bits off the io.Reader stream one byte at
// a time.
func NewBitReader(r io.Reader) (*BitReader, error) {
	return &BitReader{r, 0, 0}, nil
}

// See Reader.ReadBit() doc
func (this *BitReader) ReadBit() (int, error) {
	err := this.getByteFromStream()
	if err != nil {
		return 0, err
	}

	this.numBits -= 1
	if (this.storedByte & (1 << (this.numBits))) > 0 {
		return 1, nil
	} else {
		return 0, nil
	}
}

// See Peeker.Peek() doc
func (this *BitReader) Peek() (int, error) {
	err := this.getByteFromStream()
	if err != nil {
		return 0, err
	}

	if (this.storedByte & (1 << (this.numBits - 1))) > 0 {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (this *BitReader) getByteFromStream() error {
	if this.numBits == 0 {
		p := make([]byte, 1)
		n, err := this.r.Read(p)
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("Requested 1 byte from stream but got %d", n)
		}
		this.storedByte = p[0]
		this.numBits = 8
	}
	return nil
}
