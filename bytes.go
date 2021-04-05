package goutil

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

const (
	BYTES_IN_INT64 = 8

	_kib = 1024
	_mib = 1048576
	_gib = 1073741824

	_kb = 1000
	_mb = 1000000
	_gb = 1000000000
)

var (
	Endianess    = 0
	LittleEndian = 1
	BigEndian    = 2
)

func PrettyMemString(numBytes uint64) string {
	if numBytes < _kb {
		return fmt.Sprintf("%d", numBytes)
	}
	if numBytes < _mb {
		return fmt.Sprintf("%f kB", float64(numBytes)/float64(_kb))
	}
	if numBytes < _gb {
		return fmt.Sprintf("%f MB", float64(numBytes)/float64(_mb))
	}
	return fmt.Sprintf("%f GB", float64(numBytes)/float64(_gb))
}

func DetectEndianess() {
	var x int = 0x012345678
	var p unsafe.Pointer = unsafe.Pointer(&x)
	if 0x01 == *(*byte)(p) {
		Endianess = BigEndian
	} else if (0x78 & 0xff) == ((*(*byte)(p)) & 0xff) {
		Endianess = LittleEndian
	} else {
		panic("could not determine endianness")
	}
}

func UnsafeCaseUInt64ToBytes(val uint64) []byte {
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&val)), Len: BYTES_IN_INT64, Cap: BYTES_IN_INT64}
	return *(*[]byte)(unsafe.Pointer(&hdr))
}

func UnsafeCaseBytesToUInt64(b []byte) uint64 {
	return *(*uint64)(unsafe.Pointer(&b[0]))
}

func UnsafeCaseBytesToUInt32(b []byte) uint32 {
	return *(*uint32)(unsafe.Pointer(&b[0]))
}

func UnsafePtrCaseInt64ToUint64(i *int64) *uint64 {
	return (*uint64)(unsafe.Pointer(i))
}

func UnsafePtrCaseUint64ToInt64(u *uint64) *int64 {
	return (*int64)(unsafe.Pointer(u))
}

func UnsafePtrCaseFloat64ToUint64(f *float64) *uint64 {
	return (*uint64)(unsafe.Pointer(f))
}

func UnsafePtrCaseUint64ToFloat64(u *uint64) *float64 {
	return (*float64)(unsafe.Pointer(u))
}

func Int64ToBytes(val int64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, uint64(val))
	return bs
}

func UInt64ToBytes(val uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, val)
	return bs
}

func UInt32ToBytes(val uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, val)
	return bs
}

func Uint16ToBytes(val uint16) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, val)
	return bs
}

func Uint8ToBytes(val uint8) []byte {
	bs := make([]byte, 1)
	bs[0] = byte(val)
	return bs
}

func BoolToBytes(val bool) []byte {
	if val {
		return []byte{1}
	}
	return []byte{0}
}

func BytesToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

func BytesToUInt64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func BytesToUInt32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func BytesToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func BytesToUint8(b []byte) uint8 {
	if len(b) < 1 {
		return 0
	}
	return uint8(b[0])
}

func BytesToBool(b []byte) bool {
	if len(b) < 1 {
		return false
	}
	return b[0] == 1
}

func Uint64sToBytes(vals ...uint64) []byte {
	bs := make([]byte, 8*len(vals))
	for i, v := range vals {
		binary.LittleEndian.PutUint64(bs[i*8:(i+1)*8], v)
	}
	return bs
}

func Float64ToBytes(v float64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, math.Float64bits(v))
	return bs
}

func Float64sToBytes(vals ...float64) []byte {
	bs := make([]byte, 8*len(vals))
	for i, v := range vals {
		binary.LittleEndian.PutUint64(bs[i*8:(i+1)*8], math.Float64bits(v))
	}
	return bs
}

// You must ensure that the byte slice has the right length
func BytesToUint64s(b []byte, vals ...*uint64) {
	for i, v := range vals {
		*v = binary.LittleEndian.Uint64(b[i*8 : (i+1)*8])
	}
}

func GenUint32FromString(s string) uint32 {
	bs := []byte(s)
	if len(bs) < 4 {
		// pad with empty bits
		for i := 0; i < 4-len(bs); i++ {
			bs = append(bs, byte(0))
		}
	}
	return UnsafeCaseBytesToUInt32(bs)
}

func GetInterfaceSize(value interface{}, isParent bool, log func(...interface{})) (s int64) {
	var isvalid bool
	var v reflect.Value
	switch _v := value.(type) {
	case map[string]interface{}:
		for k, mv := range _v {
			ks := GetInterfaceSize(k, true, log)  // treat as parent because it is not included in unsafe calculation
			vs := GetInterfaceSize(mv, true, log) // treat as parent because it is not included in unsafe calculation
			s += ks + vs
			if log != nil {
				log("map[string]interface{} key", k, "size:", ks, "value:", mv, "size:", vs)
			}
		}
	case string:
		s = int64(len(_v))
	case []byte:
		s = int64(len(_v))
	case int, int32, int64, uint, uint32, uint64, float64:
		if isParent { // else already included by unsafe
			return int64(binary.Size(value))
		}
	default:
		isvalid, v = getReflectValue(value)
		if !isvalid {
			break
		}
		kind := v.Kind()
		switch kind {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				fv := v.Field(i).Interface()
				fvsize := GetInterfaceSize(fv, false, log)
				if log != nil {
					log("field", v.Type().Field(i).Name, "size:", fvsize)
				}
				s += fvsize
			}
			if isParent {
				if log != nil {
					log("reflect struct size:", v.Type().Size())
				}
				s += int64(v.Type().Size())
			}
			return
		case reflect.Slice:
			s = int64(binary.Size(value))
			if s < 0 { // non fixed-size value
				s = 0
				for i := 0; i < v.Len(); i++ {
					s += GetInterfaceSize(v.Index(i).Interface(), true, log)
				}
			}
		case reflect.Map:
			if v.Len() < 1 {
				s = 0
				break
			}
			for _, k := range v.MapKeys() {
				s += GetInterfaceSize(k.Interface(), true, log)
				s += GetInterfaceSize(v.MapIndex(k).Interface(), true, log)
			}
		default:
			if value != nil {
				s = int64(binary.Size(value))
			}
			if s < 0 { // non fixed-size value
				s = 0
			}
		}
		if isParent {
			if v.IsValid() {
				s += int64(v.Type().Size())
			}
			// s += int64(v.Type().Size())
			return
		}
	}
	if isParent {
		if !v.IsValid() {
			isvalid, v = getReflectValue(value)
		}
		if isvalid {
			s += int64(v.Type().Size())
		}
	}
	return
}

func GetValueSize(value interface{}) (s int64) {
	return GetInterfaceSize(value, true, nil)
}

func GetBytes(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeBytes(b []byte, v interface{}) (err error) {
	dec := gob.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(v)
	return
}
