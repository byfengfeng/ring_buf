package ring_buf

type RingBuf interface {
	Transmit()
	Read() []byte
	Write(b []byte)
	Destroy()
}
