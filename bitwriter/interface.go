package bitwriter

// A bit level writer which outputs bits one at a time. Data will typically
// exported to the output stream at a byte boundary.
// Make sure to Flush() the stream to get all the bits.
type Interface interface {
	Writer
	Flusher
}

// Writer writes a single bit into the output stream. There is no guarante that
// the bit will appear in the output stream until after a Flush() or it reaches
// the required byte boundary
// Returns error - nil if okay, error otherwise.
type Writer interface {
	WriteBit(int) error
}

// Flusher writes out any remainig bits into the output stream.
// You should always call Flush() when using a bit writer as it will typically
// export a 8 bit alignment.
// Returns int - the number of bits flushed
// Returns error- nil if okay, error otherwise
type Flusher interface {
	Flush() (int, error)
}
