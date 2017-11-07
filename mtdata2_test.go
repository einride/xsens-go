package xsens

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test1632(t *testing.T) {
	buf := bytes.NewReader([]byte{0x24, 0x39, 0x58, 0x10, 0x00, 0x03})
	var fp xsens1632
	err := binary.Read(buf, binary.BigEndian, &fp)
	if nil != err {
		t.Error("Could not read data: ", err)
	}
	if 0x0003 != fp.Integer {
		t.Error("Integer part wrong, expected 0x0003 was ", fp.Integer)
	}
	if 0x24395810 != fp.Fraction {
		t.Error("Fraction part wrong, expected 0x24395810 was ", fp.Fraction)
	}

	f := fp.ToFloat()
	diff := math.Abs(f - 3.1415)
	if 0.00001 < diff {
		t.Error("Error when converting to float64, expected 3.1415 was ", f)
	}
}

func TestRun(t *testing.T) {
	test := assert.New(t)
	err := Open()
	test.Nil(err)
	defer Close()

	err = readmsgs()
	test.Nil(err)
}
