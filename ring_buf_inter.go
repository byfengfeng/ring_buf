package ring_buf

type RingBuf interface {
	Read() []byte
	Write(b []byte)
	Destroy()
}
