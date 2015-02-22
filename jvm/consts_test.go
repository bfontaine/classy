package jvm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyValueAsString(t *testing.T) {
	cst := JConst{}

	assert.Equal(t, "", cst.valueAsString())
}

func TestValueAsString(t *testing.T) {
	cst := JConst{value: []byte{'h', 'e', 'l', 'l', 'o'}}

	assert.Equal(t, "hello", cst.valueAsString())
}

func TestDumpValueNilArg(t *testing.T) {
	cst := JConst{value: []byte{'a'}}
	err := cst.dumpValue(nil)

	assert.NotNil(t, err)
}

func TestDumpValueTooSmallArgType(t *testing.T) {
	cst := JConst{value: []byte{1, 1, 1, 1}}
	var res int16
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	// result is truncated
	assert.Equal(t, int16(0x0101), res)
}

func TestDumpValueTooLargeArgType(t *testing.T) {
	cst := JConst{value: []byte{1, 1, 1, 1}}
	var res int64
	err := cst.dumpValue(&res)

	assert.NotNil(t, err)
}

func TestDumpValue(t *testing.T) {
	cst := JConst{value: []byte{1, 1, 1, 1}}
	var res uint32

	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, uint32(0x1010101), res)
}

func TestStringString(t *testing.T) {
	cst := JConst{tag: TAG_STRING, value: []byte{'h', 'e', 'l', 'l', 'o'}}
	assert.Equal(t, "hello", cst.String())
}

func TestStringPositiveInteger(t *testing.T) {
	cst := JConst{tag: TAG_INT, value: []byte{0, 0, 0, 1}}

	var res int32
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int32(1), res)
	assert.Equal(t, "Integer(1)", cst.String())
}

func TestStringNegativeInteger(t *testing.T) {
	cst := JConst{tag: TAG_INT, value: []byte{0xFF, 0xFF, 0xFF, 0xD6}}

	var res int32
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int32(-42), res)
	assert.Equal(t, "Integer(-42)", cst.String())
}

func TestStringNullInteger(t *testing.T) {
	cst := JConst{tag: TAG_INT, value: []byte{0, 0, 0, 0}}

	var res int32
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int32(0), res)
	assert.Equal(t, "Integer(0)", cst.String())
}

func TestStringPositiveFloat(t *testing.T) {
	cst := JConst{tag: TAG_FLOAT, value: []byte{0x42, 0x28, 0, 0}}

	var res float32
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, float32(42.0), res)
	assert.Equal(t, "Float(42.000000)", cst.String())
}

func TestStringNegativeFloat(t *testing.T) {
	cst := JConst{tag: TAG_FLOAT, value: []byte{0xc2, 0x28, 0xb8, 0x52}}

	var res float32
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, float32(-42.18), res)
	assert.Equal(t, "Float(-42.180000)", cst.String())
}

func TestStringNullFloat(t *testing.T) {
	cst := JConst{tag: TAG_FLOAT, value: []byte{0, 0, 0, 0}}

	var res float32
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, float32(0.0), res)
	assert.Equal(t, "Float(0.000000)", cst.String())
}

func TestStringPositiveLong(t *testing.T) {
	cst := JConst{tag: TAG_LONG, value: []byte{0, 0, 0, 0, 0, 0, 0, 1}}

	var res int64
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int64(1), res)
	assert.Equal(t, "Long(1)", cst.String())
}

func TestStringNegativeLong(t *testing.T) {
	cst := JConst{tag: TAG_LONG, value: []byte{0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xD6}}

	var res int64
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int64(-42), res)
	assert.Equal(t, "Long(-42)", cst.String())
}

func TestStringNullLong(t *testing.T) {
	cst := JConst{tag: TAG_LONG, value: []byte{0, 0, 0, 0, 0, 0, 0, 0}}

	var res int64
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), res)
	assert.Equal(t, "Long(0)", cst.String())
}

func TestStringPositiveDouble(t *testing.T) {
	cst := JConst{tag: TAG_DOUBLE, value: []byte{0x3f, 0xf5, 0x1e, 0xb8,
		0x51, 0xeb, 0x85, 0x1f}}

	var res float64
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, float64(1.32), res)
	assert.Equal(t, "Double(1.320000)", cst.String())
}

func TestStringNegativeDouble(t *testing.T) {
	cst := JConst{tag: TAG_DOUBLE, value: []byte{0xc0, 0x09, 0x21, 0xca, 0xc0,
		0x83, 0x12, 0x6f}}

	var res float64
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, float64(-3.1415), res)
	assert.Equal(t, "Double(-3.141500)", cst.String())
}

func TestStringNullDouble(t *testing.T) {
	cst := JConst{tag: TAG_DOUBLE, value: []byte{0, 0, 0, 0, 0, 0, 0, 0}}

	var res float64
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, float64(0), res)
	assert.Equal(t, "Double(0.000000)", cst.String())
}

func TestStringClassRef(t *testing.T) {
	cst := JConst{tag: TAG_CLASS_REF, value: []byte{0, 42}}

	var res int16
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int16(42), res)
	assert.Equal(t, "#42", cst.String())
}

func TestStringStringRef(t *testing.T) {
	cst := JConst{tag: TAG_STRING_REF, value: []byte{0, 42}}

	var res int16
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int16(42), res)
	assert.Equal(t, "#42", cst.String())
}

func TestStringMethodType(t *testing.T) {
	cst := JConst{tag: TAG_METHOD_TYPE, value: []byte{0, 42}}

	var res int16
	err := cst.dumpValue(&res)

	assert.Nil(t, err)
	assert.Equal(t, int16(42), res)
	assert.Equal(t, "#42", cst.String())
}

func TestStringFieldRef(t *testing.T) {
	cst := JConst{tag: TAG_FIELD_REF, value: []byte{0, 42, 0, 17}}
	assert.Equal(t, "#42:#17", cst.String())
}

func TestStringMethodRef(t *testing.T) {
	cst := JConst{tag: TAG_METHOD_REF, value: []byte{0, 42, 0, 17}}
	assert.Equal(t, "#42:#17", cst.String())
}

func TestStringInterfaceMethodRef(t *testing.T) {
	cst := JConst{tag: TAG_INTERFACE_METHOD_REF, value: []byte{0, 42, 0, 17}}
	assert.Equal(t, "#42:#17", cst.String())
}

func TestStringNameType(t *testing.T) {
	cst := JConst{tag: TAG_NAME_TYPE_DESC, value: []byte{0, 42, 0, 17}}
	assert.Equal(t, "#42:#17", cst.String())
}

func TestStringMethodHandle(t *testing.T) {
	// TODO
}

func TestStringInvokeDynamic(t *testing.T) {
	// TODO
}

func TestStringUnknownTag(t *testing.T) {
	cst := JConst{tag: 123, value: []byte{0, 42}}
	assert.Equal(t, "(unknown)", cst.String())
	// TODO
}
