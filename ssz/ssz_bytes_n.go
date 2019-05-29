package ssz

import (
	"fmt"
	"reflect"
	"unsafe"
	"zrnt-ssz/ssz/unsafe_util"
)

type SSZBytesN struct {
	length uint32
}

func NewSSZBytesN(typ reflect.Type) (*SSZBytesN, error) {
	if typ.Kind() != reflect.Array {
		return nil, fmt.Errorf("typ is not a fixed-length bytes array")
	}
	if typ.Elem().Kind() != reflect.Uint8 {
		return nil, fmt.Errorf("typ is not a bytes array")
	}
	length := typ.Len()
	res := &SSZBytesN{length: uint32(length)}
	return res, nil
}

func (v *SSZBytesN) FixedLen() uint32 {
	return v.length
}

func (v *SSZBytesN) IsFixed() bool {
	return true
}

func (v *SSZBytesN) Encode(eb *sszEncBuf, p unsafe.Pointer) {
	sh := unsafe_util.GetSliceHeader(p, v.length)
	data := *(*[]byte)(unsafe.Pointer(sh))
	eb.Write(data)
}

func (v *SSZBytesN) Decode(p unsafe.Pointer) {
	// TODO
}
func (v *SSZBytesN) Ignore() {
	// TODO skip ahead Length bytes in input
}