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

func bytesToInt(buf []byte) int {
	intVal := 0

	// not sure this work with negative numbers
	for _, b := range buf {
		// big endian
		intVal <<= 8
		intVal += int(b)
	}

	return intVal
}

func readInt(r io.Reader, data *int, length int) error {
	buf := make([]byte, length)
	if err := readBinary(r, buf); err != nil {
		return err
	}

	*data = bytesToInt(buf)
	return nil
}

func setFlag(n int, magic int, flag *bool) {
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
	if err := readInt(f, &cls.minorVersion, 2); err != nil {
		return cls, err
	}

	// major version number
	if err := readInt(f, &cls.majorVersion, 2); err != nil {
		return cls, err
	}

	// constant pool size
	var constantPoolSize int
	if err := readInt(f, &constantPoolSize, 2); err != nil {
		return cls, err
	}

	// constant pool
	if err := cls.parseConstantPool(constantPoolSize, f); err != nil {
		return cls, err
	}

	// access flags
	var accessFlags int
	if err := readInt(f, &accessFlags, 2); err != nil {
		return cls, err
	}
	cls.setAccessFlags(accessFlags)

	// this class
	if err := readInt(f, &cls.classIndex, 2); err != nil {
		return cls, err
	}

	// super class
	if err := readInt(f, &cls.superClassIndex, 2); err != nil {
		return cls, err
	}

	/* TODO

	// interfaces
	var interfacesCount int
	if err := readInt(f, &interfacesCount, 2); err != nil {
		return cls, err
	}

	interfaces := make([]byte, interfacesCount)
	if err := parseInterfaces(&cls, interfaces); err != nil {
		return cls, err
	}

	// fields
	var fieldsCount int
	if err := readInt(f, &fieldsCount, 2); err != nil {
		return cls, err
	}

	fields := make([]byte, fieldsCount)
	if err := parseFields(&cls, fields); err != nil {
		return cls, err
	}

	// methods
	var methodsCount int
	if err := readInt(f, &methodsCount, 2); err != nil {
		return cls, err
	}

	methods := make([]byte, methodsCount)
	if err := parseMethods(&cls, methods); err != nil {
		return cls, err
	}

	// attributes
	var attrsCount int
	if err := readInt(f, &attrsCount, 2); err != nil {
		return cls, err
	}

	attrs := make([]byte, attrsCount)
	if err := parseAttrs(&cls, attrs); err != nil {
		return cls, err
	}
	*/

	return cls, nil
}
