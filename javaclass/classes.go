package javaclass

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	ACC_PUBLIC     = 0x0001
	ACC_FINAL      = 0x0010
	ACC_SUPER      = 0x0020
	ACC_INTERFACE  = 0x0200
	ACC_ABSTRACT   = 0x0400
	ACC_SYNTHETIC  = 0x1000
	ACC_ANNOTATION = 0x2000
	ACC_ENUM       = 0x4000
)

type u1 uint8
type u2 uint16
type u4 uint32

type JClass struct {
	majorVersion u2
	minorVersion u2

	publicFlag    bool
	finalFlag     bool
	superFlag     bool
	interfaceFlag bool
	abstractFlag  bool

	constants []JConst

	classIndex      u2
	superClassIndex u2

	interfaces []u2
}

func (cls *JClass) initConstantPool(size u2) {
	cls.constants = make([]JConst, size)
}

func (cls *JClass) addConstant(index u2, tag u1, data []byte) {
	cls.constants[index] = JConst{tag, data}
}

func (cls *JClass) parseConstantPool(constantPoolSize u2, r io.Reader) error {

	var tag u1
	var size, index u2
	var data []byte

	cls.initConstantPool(constantPoolSize)

	for index = 1; index < constantPoolSize; index++ {
		if err := readU1(r, &tag); err != nil {
			return err
		}

		switch tag {
		case TAG_STRING:
			if err := readU2(r, &size); err != nil {
				return err
			}

		case TAG_CLASS_REF:
			fallthrough
		case TAG_STRING_REF:
			fallthrough
		case TAG_METHOD_TYPE:
			size = 2

		case TAG_METHOD_HANDLE:
			size = 3

		case TAG_INT:
			fallthrough
		case TAG_FLOAT:
			fallthrough
		case TAG_FIELD_REF:
			fallthrough
		case TAG_METHOD_REF:
			fallthrough
		case TAG_INTERFACE_METHOD_REF:
			fallthrough
		case TAG_NAME_TYPE_DESC:
			fallthrough
		case TAG_INVOKE_DYN:
			size = 4

		case TAG_LONG:
			fallthrough
		case TAG_DOUBLE:
			size = 4
			// 8-bytes constants take two slots in the constant pool table
			index += 1
			size = 8

		default:
			return errors.New(fmt.Sprintf("Unknown tag '%d'", tag))
		}

		data = make([]byte, size)
		if err := readBinary(r, &data); err != nil {
			return err
		}
		cls.addConstant(index, tag, data)
	}

	return nil
}

func (cls *JClass) setAccessFlags(accessFlags u2) {
	setFlag(accessFlags, ACC_PUBLIC, &cls.publicFlag)
	setFlag(accessFlags, ACC_FINAL, &cls.finalFlag)
	setFlag(accessFlags, ACC_SUPER, &cls.superFlag)
	setFlag(accessFlags, ACC_INTERFACE, &cls.interfaceFlag)
	setFlag(accessFlags, ACC_ABSTRACT, &cls.abstractFlag)
}

func (cls *JClass) Version() string {
	// we don't use the minor version here
	switch cls.majorVersion {
	case 45:
		return "JDK 1.1"
	case 46:
		return "JDK 1.2"
	case 47:
		return "JDK 1.3"
	case 48:
		return "JDK 1.4"
	case 49:
		return "J2SE 5.0"
	case 50:
		return "J2SE 6.0"
	case 51:
		return "J2SE 7"
	case 52:
		return "J2SE 8"
	}
	return "Unknown version"
}

func (cls *JClass) StringFlags() string {
	var buffer bytes.Buffer

	if cls.publicFlag {
		buffer.WriteString("public ")
	}
	if cls.finalFlag {
		buffer.WriteString("final ")
	}
	if cls.interfaceFlag {
		buffer.WriteString("interface ")
	}
	if cls.abstractFlag {
		buffer.WriteString("abstract ")
	}

	return buffer.String()
}

func (cls *JClass) resolveConstantIndex(index u2) JConst {
	cst := cls.constants[index]

	switch cst.tag {
	case TAG_CLASS_REF:
		fallthrough
	case TAG_STRING_REF:
		fallthrough
	case TAG_METHOD_TYPE:
		return cls.resolveConstantIndex(bytesToU2(cst.value))
	default:
		return cst
	}
}

func (cls *JClass) StringConstantsIndent(indent int) string {
	var buffer bytes.Buffer

	lineStart := strings.Repeat(" ", indent)

	for idx, cst := range cls.constants {
		if idx == 0 {
			continue
		}

		buffer.WriteString(fmt.Sprintf("%s%3d = %v\n", lineStart, idx, cst))
	}

	return buffer.String()
}

func (cls *JClass) ClassName() string {
	return cls.resolveConstantIndex(cls.classIndex).valueAsString()
}
