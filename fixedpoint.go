package xsens

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// fixed to floating point conversion factors
const (
	factorFP1220 = 1 << 20
	factorFP1632 = 1 << 32
)

// FP1220 is a fixed point 12.20 value.
//
// The 12.20 fixed point output is calculated with:
//
//  int32_t fixedPointValue12p20 = round(floatingPointValue * 2^20)
//
// The resulting 32bit integer value is transmitted in big-endian order (MSB first).
//
// The range of a 12.20 fixed point value is [-2048.0, 2047.9999990].
type FP1220 [4]byte

// String returns a string representation of the value.
func (fp FP1220) String() string {
	return fmt.Sprintf("FP1220(%s)", hex.EncodeToString(fp[:]))
}

// Float64 returns the value as a 64-bit floating point value.
func (fp FP1220) Float64() float64 {
	d := fp[:]
	u := binary.BigEndian.Uint32(d)
	i := int32(u) // reinterpret as signed
	f := float64(i) / factorFP1220
	return f
}

// FP1220 is a fixed point 16.32 value.
//
// The 16.32 fixed point output is calculated with:
//
//  int64_t fixedPointValue16p32 = round(floatPointValue * 2^32)
//
// Of the resulting 64 bit integer only the 6 least significant bytes are transmitted.
//
// If these are the bytes b0 to b5 (with b0 the LSB) they are transmitted in this order:
//
//  [b3, b2, b1, b0, b5, b4]
//  [0,  1,  2,  3,  4,  5 ]
//
// This can be interpreted as first transmitting the 32bit fractional part and then the
// 16 bit integer part, both parts are in big-endian order (MSB first).
//
// The range of a 16.32 fixed point value is [-32768.0, 32767.9999999998].
type FP1632 [6]byte

// String returns a string representation of the value.
func (fp FP1632) String() string {
	return fmt.Sprintf("FP1632(%s)", hex.EncodeToString(fp[:]))
}

// Float64 returns the value as a 64-bit floating point value.
func (fp FP1632) Float64() float64 {
	d := []byte{0, 0, fp[4], fp[5], fp[0], fp[1], fp[2], fp[3]}
	if d[2]&0x80 > 0 {
		// sign-extend to 64 bits
		d[0] = 0xff
		d[1] = 0xff
	}
	u := binary.BigEndian.Uint64(d)
	i := int64(u) // reinterpret as signed
	f := float64(i) / factorFP1632
	return f
}
