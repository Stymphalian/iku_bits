package bitreader

// The default interface for a bit reader
type Interface interface {
	Reader
	Peeker
}

type Reader interface {
	// Read one bit out of the stream. Advances the head by one bit so that
	// next ReadBit() will return the next bit.
	// Returns int - 1, or 0
	// Returns error - 	nil if okay, error otherwise.
	ReadBit() (int, error)
}

type Peeker interface {
	// Only peek at the next bit, does not move the head forward.
	// Returns int - 1, or 0
	// Returns error - 	nil if okay, error otherwise. You can use this to check
	// for the end of the stream. If error is not null this typically
	// coresponds to the end of the stream.
	Peek() (int, error)
}
