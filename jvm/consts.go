package jvm

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type JConst struct {
	tag   u1
	value []byte
}

func (cst JConst) valueAsString() string {
	// note: this won't work in some case since Java uses a modified version of
	// UTF-8
	return string(cst.value)
}

// dumpValue takes a pointer on something and "cast" this const's value in it.
// If the receiver type is too small the result will be silently truncated.
func (cst JConst) dumpValue(ret interface{}) error {
	if ret == nil {
		return ErrNullPointer
	}
	buf := bytes.NewBuffer(cst.value)
	err := binary.Read(buf, binary.BigEndian, ret)
	return err
}

func (cst JConst) String() string {
	switch cst.tag {
	case TAG_STRING:
		return cst.valueAsString()

	case TAG_INT:
		var v int32
		cst.dumpValue(&v)
		return fmt.Sprintf("Integer(%d)", v)

	case TAG_FLOAT:
		var v float32
		cst.dumpValue(&v)
		return fmt.Sprintf("Integer(%d)", v)

	case TAG_LONG:
		var v float64
		cst.dumpValue(&v)
		return fmt.Sprintf("Long(%d)", v)

	case TAG_DOUBLE:
		var v int64
		cst.dumpValue(&v)
		return fmt.Sprintf("Double(%d)", v)

	case TAG_CLASS_REF:
		fallthrough
	case TAG_STRING_REF:
		fallthrough
	case TAG_METHOD_TYPE:
		return fmt.Sprintf("#%d", bytesToInt(cst.value))

	case TAG_FIELD_REF:
		fallthrough
	case TAG_METHOD_REF:
		fallthrough
	case TAG_INTERFACE_METHOD_REF:
		fallthrough
	case TAG_NAME_TYPE_DESC:
		return fmt.Sprintf("#%d:#%d", bytesToInt(cst.value[:2]),
			bytesToInt(cst.value[2:]))

	case TAG_METHOD_HANDLE:
		// TODO

	case TAG_INVOKE_DYN:
		// TODO
	}

	return "(unknown)"
}
