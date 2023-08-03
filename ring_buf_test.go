package ring_buf

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

const (
	testSize = 50000
)

func testMap() {
	m := sync.Map{}
	for i := 0; i < 100000; i++ {
		m.Store(i, "2")
	}
	m.Range(func(key, value any) bool {
		m.Delete(key)
		return true
	})
	m.Range(func(key, value any) bool {
		fmt.Println(key, "----", value)
		return true
	})
}

func read(accept net.Conn, ringBuff RingBuf) {
	go func() {
		for {
			by := ringBuff.Read()
			if by == nil {
				return
			}
			fmt.Println(string(by))
		}
	}()
	head := make([]byte, 4)
	for {
		n, err := accept.Read(head)
		if err != nil {
			panic(err)
		}
		if n == 4 {
			l := binary.BigEndian.Uint32(head)
			data := make([]byte, l-4)
			n1, err1 := accept.Read(data)
			if err1 != nil {
				panic(err)
			}
			if uint32(n1) == l-4 {
				ringBuff.Write(data)
			}
		}
	}

}

func excode() {

}

func decode() {

}

func TestNet(t *testing.T) {
	ringBuff := NewRingBuff(nil)
	go func() {
		listen, err := net.Listen("tcp", ":9998")
		if err != nil {
			panic(err)
		}
		for {
			accept, err2 := listen.Accept()
			if err2 != nil {
				panic(err2)
			}
			go read(accept, ringBuff)
		}
	}()
	time.Sleep(1 * time.Second)
	go func() {
		conn, err := net.Dial("tcp", ":9998")
		if err != nil {
			panic(err)
		}
		t1 := []byte("123456789")
		for i := 0; i < 10000; i++ {

			//for j := 0; i > j; j++ {
			//	t1 = append(t1, []byte(fmt.Sprintf("%d", j))...)
			//}
			t1 = []byte(fmt.Sprintf("%d", i+1))
			data := Encode(t1)
			conn.Write(data)
		}
		time.Sleep(5 * time.Minute)
	}()
	time.Sleep(5 * time.Minute)
}

func TestRingBuff_Write(t *testing.T) {
	testMap()
	for j := 0; j < 1; j++ {
		go func(index int) {
			ringBuff := NewRingBuff(nil)
			checkList := make([]string, testSize)
			go func() {
				//time.Sleep(1 * time.Second)
				data := make([]byte, 2)
				//data1 := make([]byte, 2)
				c := 0
				//lock := sync.Mutex{}

				for {
					data = ringBuff.Read()
					if data == nil {
						return
					}
					c += 1
					if c > testSize {
						c = 1
					}
					fmt.Println(string(data), "--c", index, "---", checkList[c-1] == string(data), "---", c)
					if checkList[c-1] != string(data) {
						fmt.Println(111)
					}
					Put(data)
				}
			}()
			fmt.Println(time.Now().Unix())
			data := []byte("123456789")
			data1 := []byte("987654321")
			for i := 0; i < testSize; i++ {

				if i%2 == 0 {
					checkList[i] = "123456789"
					ringBuff.Write(data)
					//	//wd = ringBuff.Write(data1[4:])
				} else {
					checkList[i] = "987654321"
					ringBuff.Write(data1)
					//wd = ringBuff.Write(data[7:])
				}
				if i == 30000 {
					//time.Sleep(3 * time.Second)
					fmt.Println("-------------------------------------------------------")
				}
			}
			fmt.Println("------------------end-----------------------")
			time.Sleep(30 * time.Second)
			data = []byte("11111111")
			data1 = []byte("22222222")
			for i := 0; i < 1000; i++ {
				if i%2 == 0 {
					checkList[i] = "22222222"
					ringBuff.Write(data1)
					//	//wd = ringBuff.Write(data1[4:])
				} else {
					checkList[i] = "11111111"
					ringBuff.Write(data)
					//wd = ringBuff.Write(data[7:])
				}
			}
			time.Sleep(10 * time.Second)
			ringBuff.Destroy()
		}(j)
	}

	time.Sleep(5 * time.Minute)
}

func Encode(data []byte) []byte {
	length := uint32(len(data)) + uint32(4)
	buf := make([]byte, length)
	binary.BigEndian.PutUint32(buf, length)
	buf = append(buf[:4], data...)
	return buf
}

func TestRingBuff_Read(t *testing.T) {
	fmt.Println(time.Now().Unix())
	c := uint8(0)
	c = 1 << 1
	fmt.Println(1&(c>>1) == 1)
	c = 1 << 0
	fmt.Println(1&(c>>1) == 1)

	c = 1 << 1
	fmt.Println(1&(c>>1) == 1)
	c = 1 << 0
	fmt.Println(1&(c>>1) == 1)
	//ringBuff := NewRingBuff()
	//date := ringBuff.Read(9)
	//fmt.Println(date)
}
