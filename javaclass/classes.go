package javaclass

import (
	"bytes"
	"fmt"
	"strings"
)

type JClass struct {
	majorVersion u2
	minorVersion u2
	accessFlags  u2

	constants  []JConst
	fields     []JField
	interfaces []u2

	classIndex      u2
	superClassIndex u2
}

func (cls *JClass) initConstantPool(size u2) {
	cls.constants = make([]JConst, size)
}

func (cls *JClass) initFields(size u2) {
	cls.fields = make([]JField, size)
}

func (cls *JClass) initInterfaces(size u2) {
	cls.interfaces = make([]u2, size)
}

func (cls *JClass) addConstant(index u2, tag u1, data []byte) {
	cls.constants[index] = JConst{tag, data}
}

func (cls *JClass) Version() string {
	// we could also use the minor version here
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

func (cls *JClass) hasAccessFlag(flag u2) bool {
	return cls.accessFlags&flag == flag
}

func (cls *JClass) StringFlags() string {
	var buffer bytes.Buffer

	if cls.hasAccessFlag(ACC_PUBLIC) {
		buffer.WriteString("public ")
	}
	if cls.hasAccessFlag(ACC_FINAL) {
		buffer.WriteString("final ")
	}
	if cls.hasAccessFlag(ACC_INTERFACE) {
		buffer.WriteString("interface ")
	}
	if cls.hasAccessFlag(ACC_ABSTRACT) {
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
