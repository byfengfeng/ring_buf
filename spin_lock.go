package ring_buf

import (
	"runtime"
	"sync/atomic"
)

type SpinLock uint32

const maxRepCount = 16

// Lock 加锁
func (s *SpinLock) Lock() {
	//尝试次数
	repCount := 1
	//尝试将锁的值从0改成1
	for !atomic.CompareAndSwapUint32((*uint32)(s), 0, 1) {
		for i := 0; i < repCount; i++ {
			//协程谦让
			runtime.Gosched()
		}
		//判断锁获取的次数是否达到最大值
		if repCount < maxRepCount {
			repCount <<= 1
		}
	}
}

// Unlock 解锁
func (s *SpinLock) Unlock() {
	//将锁的值改成0
	atomic.StoreUint32((*uint32)(s), 0)
}
