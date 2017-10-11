// ************************************************************
// DO NOT EDIT.
// THIS FILE IS AUTO-GENERATED BY codecgen.
// ************************************************************

package v1alpha1

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"

	codec1978 "github.com/ugorji/go/codec"
)

const (
	// ----- content types ----
	codecSelferC_UTF88024 = 1
	codecSelferC_RAW8024  = 0
	// ----- value types used ----
	codecSelferValueTypeArray8024 = 10
	codecSelferValueTypeMap8024   = 9
	// ----- containerStateValues ----
	codecSelfer_containerMapKey8024    = 2
	codecSelfer_containerMapValue8024  = 3
	codecSelfer_containerMapEnd8024    = 4
	codecSelfer_containerArrayElem8024 = 6
	codecSelfer_containerArrayEnd8024  = 7
)

var (
	codecSelferBitsize8024                         = uint8(reflect.TypeOf(uint(0)).Bits())
	codecSelferOnlyMapOrArrayEncodeToStructErr8024 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer8024 struct{}

func init() {
	if codec1978.GenVersion != 5 {
		_, file, _, _ := runtime.Caller(0)
		err := fmt.Errorf("codecgen version mismatch: current: %v, need %v. Re-generate file: %v",
			5, codec1978.GenVersion, file)
		panic(err)
	}
	if false { // reference the types, but skip this branch at build/run time
	}
}

func (x PrefixType) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer8024
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	yym1 := z.EncBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.EncExt(x) {
	} else if !yym1 && z.IsJSONHandle() {
		z.EncJSONMarshal(x)
	} else {
		r.EncodeInt(int64(x))
	}
}

func (x *PrefixType) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer8024
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym1 := z.DecBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else if !yym1 && z.IsJSONHandle() {
		z.DecJSONUnmarshal(x)
	} else {
		*((*int)(x)) = int(r.DecodeInt(codecSelferBitsize8024))
	}
}
