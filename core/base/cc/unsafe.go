package cc

import (
	"sync/atomic"
	"unsafe"
)

func StoreInt8(addr *int8, val int8) {
	unsafeptr := (*int32)(unsafe.Pointer(addr))
	atomic.StoreInt32(unsafeptr, int32(val))
}

func LoadInt8(addr *int8) (val int8) {
	unsafeptr := (*int32)(unsafe.Pointer(addr))
	val = int8(atomic.LoadInt32(unsafeptr))
	return
}

func StoreUint8(addr *uint8, val uint8) {
	unsafeptr := (*uint32)(unsafe.Pointer(addr))
	atomic.StoreUint32(unsafeptr, uint32(val))
}

func LoadUint8(addr *uint8) (val uint8) {
	unsafeptr := (*uint32)(unsafe.Pointer(addr))
	val = uint8(atomic.LoadUint32(unsafeptr))
	return
}

func StoreInt16(addr *int16, val int16) {
	unsafeptr := (*int32)(unsafe.Pointer(addr))
	atomic.StoreInt32(unsafeptr, int32(val))
}

func LoadInt16(addr *int16) (val int16) {
	unsafeptr := (*int32)(unsafe.Pointer(addr))
	val = int16(atomic.LoadInt32(unsafeptr))
	return
}

func StoreUint16(addr *uint16, val uint16) {
	unsafeptr := (*uint32)(unsafe.Pointer(addr))
	atomic.StoreUint32(unsafeptr, uint32(val))
}

func LoadUint16(addr *uint16) (val uint16) {
	unsafeptr := (*uint32)(unsafe.Pointer(addr))
	val = uint16(atomic.LoadUint32(unsafeptr))
	return
}
