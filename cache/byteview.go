package cache

type ByteView struct {
	b []byte
}

func (this ByteView) Len() int {
	return len(this.b)
}

func (this ByteView) String() string {
	return string(this.b)
}

func (this ByteView) ByteSlice() []byte {
	return cloneByte(this.b)
}

func cloneByte(b []byte) []byte {
	clone := make([]byte, len(b))
	copy(clone, b)
	return clone
}
