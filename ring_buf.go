package ring_buf

import (
	"encoding/binary"
)

const (
	zero                = 0
	dataPackSize        = 2
	shrinkageCountMax   = 10
	ringBufSize         = 1024
	bufferGrowThreshold = ringBufSize * 4
)

type ringBuf struct {
	shrinkageCount uint8         //缩容阈值计数 protected add or shrinkage buf size operate
	readWait       uint8         //读等待  is the read wait
	isDataEmpty    bool          //数据是否可读  is the data readable
	closeStatus    bool          //关闭状态 ring buf close status
	bufSize        uint32        //实体大小 array size
	rPos           uint32        //读取位置 read position
	wPos           uint32        //写入位置 write position
	buf            []byte        //实体数据 data array
	enCodeBuf      []byte        //包装头部数据 pack data head
	reqTransmit    chan []byte   //ringBuf data transmit
	resTransmit    chan []byte   //ringBuf data transmit
	readSignal     chan struct{} //ringBuf readSignal
	destroySignal  chan struct{} //ringBuf destroySignal signal
}

func NewRingBuff(resTransmit chan []byte) RingBuf {
	if resTransmit == nil {
		resTransmit = make(chan []byte)
	}
	r := &ringBuf{
		isDataEmpty:   true,
		bufSize:       ringBufSize,
		buf:           make([]byte, ringBufSize),
		enCodeBuf:     Get(dataPackSize),
		reqTransmit:   make(chan []byte),
		resTransmit:   resTransmit,
		readSignal:    make(chan struct{}),
		destroySignal: make(chan struct{}),
	}
	defer r.Init()
	return r
}

// Init 环形缓存读取初始化 ringBuf read Init
func (r *ringBuf) Init() {
	go r.ringEv()
}

// Read 读取数据 read ringBuf data
func (r *ringBuf) Read() []byte {
	r.readSignal <- struct{}{}
	return <-r.resTransmit
}

func (r *ringBuf) Write(b []byte) {
	r.reqTransmit <- b
}

// Write 数据写入,并发不安全，使其变并发安全即失去原本效果 data writing, concurrency is not safe, making it concurrency safe will lose the original effect
func (r *ringBuf) write(b []byte) (l uint32) {
	b = r.enCode(b)
	l = uint32(len(b))
	if l == zero {
		return
	}
	if (r.bufSize-r.rPos)/ringBufSize > 1 {
		if r.shrinkageCount < shrinkageCountMax {
			r.shrinkageCount++
		} else {
			r.shrinkage(l)
		}

	}

	free := r.available()
	if l > free {
		r.shrinkageCount = zero
		r.grow(r.bufSize + l - free)
	}
	//写入
	if r.wPos >= r.rPos {
		bufAseSize := r.bufSize - r.wPos
		//查看是否需要往前写
		if bufAseSize < l {
			needSize := l - bufAseSize
			copy(r.buf[r.wPos:], b[:bufAseSize])
			copy(r.buf, b[bufAseSize:])
			r.wPos = needSize
		} else {
			copy(r.buf[r.wPos:], b)
			r.wPos += l
		}
	} else {
		copy(r.buf[r.wPos:], b)
		r.wPos += l
	}
	r.isDataEmpty = false
	return
}

// Destroy 销毁 ring buf destroy
func (r *ringBuf) Destroy() {
	r.closeStatus = true
	close(r.destroySignal)
	close(r.reqTransmit)
	close(r.readSignal)
	close(r.resTransmit)
	Put(r.enCodeBuf)
	Put(r.buf)
}

// read 数据读出，并发安全，为了避免在扩容读取位置改变 Data reading, concurrency safety, in order to avoid changing the reading position during expansion
func (r *ringBuf) read(p []byte) (n uint32) {
	if len(p) == zero || r.isDataEmpty {
		return zero
	}
	if r.wPos > r.rPos {
		n = r.wPos - r.rPos
		if n > uint32(len(p)) {
			n = uint32(len(p))
		}
		copy(p, r.buf[r.rPos:r.rPos+n])
		r.rPos += n
		if r.rPos == r.wPos {
			r.Reset()
		}
		return
	}

	n = r.bufSize - r.rPos + r.wPos
	if n > uint32(len(p)) {
		n = uint32(len(p))
	}

	if r.rPos+n <= r.bufSize {
		copy(p, r.buf[r.rPos:r.rPos+n])
	} else {
		orthogonalPos := r.bufSize - r.rPos
		copy(p, r.buf[r.rPos:])
		negativePos := n - orthogonalPos
		copy(p[orthogonalPos:], r.buf[:negativePos])
	}

	r.rPos = (r.rPos + n) % r.bufSize
	if r.rPos == r.wPos {
		r.Reset()
	}
	return
}

// Reset 重置ringBuf读写位置  reset ringBuf read and write position
func (r *ringBuf) Reset() {
	r.isDataEmpty = true
	r.rPos, r.wPos = zero, zero
}

// shrinkage 判断是否缩容 judging whether to shrink
func (r *ringBuf) shrinkage(l uint32) {
	if r.bufSize/ringBufSize > 1 && r.shrinkageCount == shrinkageCountMax && r.wPos == r.rPos && r.rPos == zero {
		r.isDataEmpty = true
		newBuf := Get(ringBufSize)
		Put(r.buf)
		r.buf = newBuf
		r.bufSize = ringBufSize
		r.shrinkageCount = zero
	}
}

// available 判断bufSize可用长度 Judge the available length of bufSize
func (r *ringBuf) available() (n uint32) {
	if r.rPos == r.wPos {
		if r.isDataEmpty {
			return r.bufSize
		}
		return zero
	}

	if r.wPos < r.rPos {
		return r.rPos - r.wPos
	}

	return r.bufSize - r.wPos + r.rPos
}

// grow 扩容 expansion
func (r *ringBuf) grow(newCap uint32) {
	n := r.bufSize
	doubleCap := n + n
	if newCap <= doubleCap {
		if n < bufferGrowThreshold {
			newCap = doubleCap
		} else {
			for zero < n && n < newCap {
				n += n / 4
			}
			// The n calculation doesn't overflow, set n to newCap.
			if n > zero {
				newCap = n
			}
		}
	}
	newBuf := Get(int(newCap))
	oldLen := r.buffered()
	r.read(newBuf)
	Put(r.buf)
	r.buf = newBuf
	r.rPos = zero
	r.wPos = oldLen
	r.bufSize = newCap
	if r.wPos > zero {
		r.isDataEmpty = false
	}
}

// buffered 计算扩容之后写入位置 write the position after calculating the expansion
func (r *ringBuf) buffered() uint32 {
	if r.rPos == r.wPos {
		if r.isDataEmpty {
			return zero
		}
		return r.bufSize
	}

	if r.wPos > r.rPos {
		return r.wPos - r.rPos
	}

	return r.bufSize - r.rPos + r.wPos
}

// enCode 数据包装，确保接收和发送是同一数据 data packaging to ensure that the same data is received and sent
func (r *ringBuf) enCode(bytes []byte) []byte {
	binary.BigEndian.PutUint16(r.enCodeBuf, uint16(len(bytes)))
	return append(r.enCodeBuf, bytes...)
}

// ringEv
func (r *ringBuf) ringEv() {
	length := Get(dataPackSize)
	for {
		select {
		case <-r.readSignal:
			n := r.read(length)
			if n > zero {
				l := binary.BigEndian.Uint16(length)
				bytes := Get(int(l))
				n1 := r.read(bytes)
				if n1 > zero {
					r.resTransmit <- bytes
				}
			} else {
				r.readWait++
			}
		case data := <-r.reqTransmit:
			r.write(data)
			for r.readWait > 0 {
				n := r.read(length)
				if n > zero {
					l := binary.BigEndian.Uint16(length)
					bytes := Get(int(l))
					n1 := r.read(bytes)
					if n1 > zero {
						r.resTransmit <- bytes
						r.readWait--
					}
				}
			}
		case <-r.destroySignal:
			return
		}
	}
}
