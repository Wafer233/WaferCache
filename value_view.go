package WaferCache

// 个人理解，这一步是跟并发没什么关系。 主要是创建一个KV的深拷贝，使得我们得到value的时候不会被直接篡改
type ValueView struct {
	b []byte
}

func (view ValueView) Len() int {
	return len(view.b)
}

func (view ValueView) ByteSlice() []byte {
	return cloneBytes(view.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (view ValueView) String() string {
	return string(view.b)
}
