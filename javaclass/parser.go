package javaclass

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

var (
	ErrNotEnoughBytes   = errors.New("Can't read enough bytes")
	ErrWrongMagicNumber = errors.New("Wrong magic number")
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

func readU2(r io.Reader, data *u2) error {
	buf, err := readBinaryBuffer(r, 2)

	if err != nil {
		return err
	}

	*data = bytesToU2(buf)
	return nil
}

func readU1(r io.Reader, data *u1) error {
	buf, err := readBinaryBuffer(r, 2)

	if err != nil {
		return err
	}

	*data = u1(buf[0])
	return nil
}

func setFlag(n u2, magic u2, flag *bool) {
	*flag = n&magic == magic
}

func parseInterfaces(cls *JClass, interfaces []byte) error {
	// TODO
	return nil
}

func parseFields(cls *JClass, fields []byte) error {
	// TODO
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
	if err := cls.parseConstantPool(constantPoolSize, f); err != nil {
		return cls, err
	}

	// access flags
	var accessFlags u2
	if err := readU2(f, &accessFlags); err != nil {
		return cls, err
	}
	cls.setAccessFlags(accessFlags)

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

	cls.interfaces = make([]u2, interfacesCount)
	if err := readBinary(f, &cls.interfaces); err != nil {
		return cls, err
	}

	// fields
	var fieldsCount u2
	if err := readU2(f, &fieldsCount); err != nil {
		return cls, err
	}

	/*

		fields := make([]byte, fieldsCount)
		if err := parseFields(&cls, fields); err != nil {
			return cls, err
		}

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
