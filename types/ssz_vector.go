package types

import (
	"fmt"
	"reflect"
	"unsafe"
	. "zssz/dec"
	. "zssz/enc"
	. "zssz/htr"
)

type VectorLength interface {
	VectorLength() uint32
}

type SSZVector struct {
	length uint32
	elemMemSize uintptr
	elemSSZ SSZ
	isFixedLen bool
	fixedLen uint32
}

func NewSSZVector(factory SSZFactoryFn, typ reflect.Type) (*SSZVector, error) {
	if typ.Kind() != reflect.Array {
		return nil, fmt.Errorf("typ is not a fixed-length array")
	}
	length := uint32(typ.Len())
	elemTyp := typ.Elem()

	elemSSZ, err := factory(elemTyp)
	if err != nil {
		return nil, err
	}
	res := &SSZVector{
		length: length,
		elemMemSize: elemTyp.Size(),
		elemSSZ: elemSSZ,
		isFixedLen: elemSSZ.IsFixed(),
		fixedLen: elemSSZ.FixedLen() * length,
	}
	return res, nil
}

func (v *SSZVector) VectorLength() uint32 {
	return v.length
}

func (v *SSZVector) FixedLen() uint32 {
	return v.fixedLen
}

func (v *SSZVector) IsFixed() bool {
	return v.isFixedLen
}

func (v *SSZVector) Encode(eb *EncodingBuffer, p unsafe.Pointer) {
	if v.elemSSZ.IsFixed() {
		EncodeFixedSeries(v.elemSSZ.Encode, v.length, v.elemMemSize, eb, p)
	} else {
		EncodeVarSeries(v.elemSSZ.Encode, v.length, v.elemMemSize, eb, p)
	}
}

func (v *SSZVector)  Decode(dr *DecodingReader, p unsafe.Pointer) error {
	if v.elemSSZ.IsFixed() {
		return DecodeFixedSeries(v.elemSSZ.Decode, v.length, v.elemMemSize, dr, p)
	} else {
		return DecodeVarSeries(v.elemSSZ.Decode, v.length, v.elemMemSize, dr, p)
	}
}

func (v *SSZVector) HashTreeRoot(h *Hasher, p unsafe.Pointer) [32]byte {
	elemHtr := v.elemSSZ.HashTreeRoot
	elemSize := v.elemMemSize
	leaf := func(i uint32) []byte {
		v := elemHtr(h, unsafe.Pointer(uintptr(p)+(elemSize * uintptr(i))))
		return v[:]
	}
	return Merkleize(h, v.length, leaf)
}