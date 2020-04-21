package utils

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	Unlocked int32 = 0
	Locked   int32 = 1 << Unlocked
)

type CASLocker struct {
	mutex  sync.Mutex
	status *int32
}

func NewCASLocker() *CASLocker {
	status := Unlocked
	return &CASLocker{
		mutex:  sync.Mutex{},
		status: &status,
	}
}

func (cl *CASLocker) Lock() {
	atomic.AddInt32(cl.status, Locked)
	cl.mutex.Lock()
}

func (cl *CASLocker) UnLock() {
	cl.mutex.Unlock()
	atomic.AddInt32(cl.status, Unlocked)
}

func (cl *CASLocker) TryLock() bool {
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&cl.mutex)), Unlocked, Locked) {
		atomic.AddInt32(cl.status, Locked)
		return true
	}
	return false
}

func (cl *CASLocker) IsLocked() bool {
	if atomic.LoadInt32(cl.status) == Locked {
		return true
	}
	return false
}
