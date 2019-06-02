package unsafe_util

import (
	"unsafe"
)

// Allocates a new slice of the given length, with the given space for each element,
// and write a slice-header for it to the pointer.
// Note: p is assumed to be a pointer to a slice header,
// and the pointer is assumed to keep the referenced data alive, away from the GC.
func AllocateSliceSpaceAndBind(p unsafe.Pointer, length uint32, elemMemSize uintptr) unsafe.Pointer {
	dataLen := uint32(elemMemSize) * length
	data := make([]byte, 0, dataLen)
	contentsPtr := unsafe.Pointer(&data)
	sh := GetSliceHeader(contentsPtr, length)
	*(*SliceHeader)(p) = *sh
	return contentsPtr
}
