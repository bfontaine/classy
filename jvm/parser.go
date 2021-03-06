package jvm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

func readBytes(f *os.File, buff []byte) error {
	n, err := f.Read(buff)

	if err != nil {
		return err
	}

	if n < len(buff) {
		return ErrNotEnoughBytes
	}

	return nil
}

func readBinary(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}

func bytesToU2(buf []byte) u2 {
	return u2((buf[0] << 8) + buf[1])
}

func bytesToInt(buf []byte) int {
	intVal := 0

	// not sure this works with negative numbers
	for _, b := range buf {
		// big endian
		intVal <<= 8
		intVal += int(b)
	}

	return intVal
}

func readBinaryBuffer(r io.Reader, size int) ([]byte, error) {
	buf := make([]byte, size)
	if err := readBinary(r, buf); err != nil {
		return buf, err
	}
	return buf, nil
}

func readU4(r io.Reader, data *u4) error {
	buf, err := readBinaryBuffer(r, 4)

	if err != nil {
		return err
	}

	*data = u4(buf[0]<<24 + buf[1]<<16 + buf[2]<<8 + buf[3])
	return nil
}

func readU2(r io.Reader, data *u2) error {
	buf, err := readBinaryBuffer(r, 2)

	if err != nil {
		return err
	}

	*data = bytesToU2(buf)
	return nil
}

func readU1(r io.Reader, data *u1) error {
	buf, err := readBinaryBuffer(r, 1)

	if err != nil {
		return err
	}

	*data = u1(buf[0])
	return nil
}

func parseConstantPool(cls *JClass, constantPoolSize u2, r io.Reader) error {

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

		case TAG_CLASS_REF, TAG_STRING_REF, TAG_METHOD_TYPE:
			size = 2

		case TAG_METHOD_HANDLE:
			size = 3

		case TAG_INT, TAG_FLOAT, TAG_FIELD_REF, TAG_METHOD_REF,
			TAG_INTERFACE_METHOD_REF, TAG_NAME_TYPE_DESC, TAG_INVOKE_DYN:
			size = 4

		case TAG_LONG, TAG_DOUBLE:
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

func parseFieldAttributes(field *JField, attrsCount u2, r io.Reader) error {

	var index u2

	field.initAttributes(attrsCount)

	for index = 0; index < attrsCount; index++ {
		attr := JAttr{}

		if err := readU2(r, &attr.nameIndex); err != nil {
			return err
		}

		var attrLen u4
		if err := readU4(r, &attrLen); err != nil {
			return err
		}

		attr.initInfo(attrLen)
		if err := readBinary(r, &attr.info); err != nil {
			return err
		}

		field.addAttribute(index, attr)
	}

	return nil
}

func parseFields(cls *JClass, fieldsCount u2, r io.Reader) error {

	var index u2

	cls.initFields(fieldsCount)

	for index = 0; index < fieldsCount; index++ {
		field := JField{}

		if err := readU2(r, &field.accessFlags); err != nil {
			return err
		}

		if err := readU2(r, &field.nameIndex); err != nil {
			return err
		}

		if err := readU2(r, &field.descriptorIndex); err != nil {
			return err
		}

		var attributesCount u2
		if err := readU2(r, &attributesCount); err != nil {
			return err
		}

		if err := parseFieldAttributes(&field, attributesCount, r); err != nil {
			return err
		}

		cls.addField(index, field)
	}
	return nil
}

func parseMethods(cls *JClass, methods []byte) error {
	// TODO
	return nil
}

func parseAttrs(cls *JClass, attrs []byte) error {
	// TODO
	return nil
}

// ParseClassFromFiles takes an open file and parses a JClass from it.
func ParseClassFromFile(f *os.File) (JClass, error) {

	// magic number
	buf4 := make([]byte, 4)
	if err := readBytes(f, buf4); err != nil {
		return JClass{}, err
	}

	if !bytes.Equal(buf4, []byte{0xCA, 0xFE, 0xBA, 0xBE}) {
		return JClass{}, ErrWrongMagicNumber
	}

	cls := JClass{}

	// minor version number
	if err := readU2(f, &cls.minorVersion); err != nil {
		return cls, err
	}

	// major version number
	if err := readU2(f, &cls.majorVersion); err != nil {
		return cls, err
	}

	// constant pool size
	var constantPoolSize u2
	if err := readU2(f, &constantPoolSize); err != nil {
		return cls, err
	}

	// constant pool
	if err := parseConstantPool(&cls, constantPoolSize, f); err != nil {
		return cls, err
	}

	// access flags
	if err := readU2(f, &cls.accessFlags); err != nil {
		return cls, err
	}

	// this class
	if err := readU2(f, &cls.classIndex); err != nil {
		return cls, err
	}

	// super class
	if err := readU2(f, &cls.superClassIndex); err != nil {
		return cls, err
	}

	// interfaces
	var interfacesCount u2
	if err := readU2(f, &interfacesCount); err != nil {
		return cls, err
	}

	cls.initInterfaces(interfacesCount)
	if err := readBinary(f, &cls.interfaces); err != nil {
		return cls, err
	}

	// fields
	var fieldsCount u2
	if err := readU2(f, &fieldsCount); err != nil {
		return cls, err
	}

	if err := parseFields(&cls, fieldsCount, f); err != nil {
		return cls, err
	}

	/*
		// methods
		var methodsCount u2
		if err := readU2(f, &methodsCount); err != nil {
			return cls, err
		}

		methods := make([]byte, methodsCount)
		if err := parseMethods(&cls, methods); err != nil {
			return cls, err
		}

		// attributes
		var attrsCount u2
		if err := readU2(f, &attrsCount); err != nil {
			return cls, err
		}

		attrs := make([]byte, attrsCount)
		if err := parseAttrs(&cls, attrs); err != nil {
			return cls, err
		}
	*/

	return cls, nil
}
