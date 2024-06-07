package sandbox

import "bytes"

// Buffer is a named bytes.Buffer.
type Buffer struct {
	name string
	*bytes.Buffer
}

func NewBuffer(name string, data []byte) *Buffer {
	return &Buffer{name, bytes.NewBuffer(data)}
}

// Name returns the name of the buffer.
func (b Buffer) Name() string {
	return b.name
}

// Close is a no-op.
func (b Buffer) Close() error {
	return nil
}
