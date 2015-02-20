package main

// most of the code here is based on informations found here:
//  - https://en.wikipedia.org/wiki/Java_class_file
//  - http://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	TAG_STRING               = 1
	TAG_INT                  = 3
	TAG_FLOAT                = 4
	TAG_LONG                 = 5
	TAG_DOUBLE               = 6
	TAG_CLASS_REF            = 7
	TAG_STRING_REF           = 8
	TAG_FIELD_REF            = 9
	TAG_METHOD_REF           = 10
	TAG_INTERFACE_METHOD_REF = 11
	TAG_NAME_TYPE_DESC       = 12
	TAG_METHOD_HANDLE        = 15
	TAG_METHOD_TYPE          = 16
	TAG_INVOKE_DYN           = 18

	ACC_PUBLIC     = 0x0001
	ACC_FINAL      = 0x001
	ACC_SUPER      = 0x0020
	ACC_INTERFACE  = 0x0200
	ACC_ABSTRACT   = 0x0400
	ACC_SYNTHETIC  = 0x1000
	ACC_ANNOTATION = 0x2000
	ACC_ENUM       = 0x4000
)

var (
	ErrNotEnoughBytes   = errors.New("Can't read enough bytes")
	ErrWrongMagicNumber = errors.New("Wrong magic number")
)

type jconst struct {
	tag   int
	value []byte
}

type jclass struct {
	majorVersion int
	minorVersion int

	publicFlag    bool
	finalFlag     bool
	superFlag     bool
	interfaceFlag bool
	abstractFlag  bool

	constants []jconst

	classIndex      int
	superClassIndex int
}

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

func (cls *jclass) initConstantPool(size int) {
	cls.constants = make([]jconst, size)
}

func (cls *jclass) addConstant(index, tag int, data []byte) {
	cls.constants[index] = jconst{tag, data}
}

func (cls *jclass) parseConstantPool(constantPoolSize int, r io.Reader) error {

	var tag, size int
	var data []byte

	cls.initConstantPool(constantPoolSize)

	for index := 1; index < constantPoolSize; index++ {
		if err := readInt(r, &tag, 1); err != nil {
			return err
		}

		switch tag {
		case TAG_STRING:
			if err := readInt(r, &size, 2); err != nil {
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

func parseInterfaces(cls *jclass, interfaces []byte) error {
	// TODO
	return nil
}

func parseFields(cls *jclass, fields []byte) error {
	// TODO
	return nil
}

func parseMethods(cls *jclass, methods []byte) error {
	// TODO
	return nil
}

func parseAttrs(cls *jclass, attrs []byte) error {
	// TODO
	return nil
}

func inspectFilename(source string) (jclass, error) {
	f, err := os.Open(source)
	if err != nil {
		return jclass{}, err
	}

	defer f.Close()

	// magic number
	buf4 := make([]byte, 4)
	if err := readBytes(f, buf4); err != nil {
		return jclass{}, err
	}

	if !bytes.Equal(buf4, []byte{0xCA, 0xFE, 0xBA, 0xBE}) {
		return jclass{}, ErrWrongMagicNumber
	}

	cls := jclass{}

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
	cls.SetAccessFlags(accessFlags)

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

func setFlag(n int, magic int, flag *bool) {
	*flag = n&magic == magic
}

func (cls *jclass) SetAccessFlags(accessFlags int) {
	setFlag(accessFlags, ACC_PUBLIC, &cls.publicFlag)
	setFlag(accessFlags, ACC_FINAL, &cls.finalFlag)
	setFlag(accessFlags, ACC_SUPER, &cls.superFlag)
	setFlag(accessFlags, ACC_INTERFACE, &cls.interfaceFlag)
	setFlag(accessFlags, ACC_ABSTRACT, &cls.abstractFlag)
}

func (cls *jclass) Version() string {
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

func (cls *jclass) StringFlags() string {
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

func (cls *jclass) resolveConstantIndex(index int) jconst {
	cst := cls.constants[index]

	switch cst.tag {
	case TAG_CLASS_REF:
		fallthrough
	case TAG_STRING_REF:
		fallthrough
	case TAG_METHOD_TYPE:
		return cls.resolveConstantIndex(bytesToInt(cst.value))
	default:
		return cst
	}
}

func (cst jconst) valueAsString() string {
	return string(cst.value)
}

func (cst jconst) dumpValue(ret interface{}) error {
	buf := bytes.NewBuffer(cst.value)
	err := binary.Read(buf, binary.BigEndian, &ret)
	return err
}

func (cst jconst) valueAsInt64() uint64 {
	return binary.BigEndian.Uint64(cst.value)
}

func (cst jconst) String() string {
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

func (cls *jclass) StringConstantsIndent(indent int) string {
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

func (cls *jclass) ClassName() string {
	return cls.resolveConstantIndex(cls.classIndex).valueAsString()
}

func printClass(filename string, cls jclass) {
	fmt.Printf("%s:\n"+
		"  class:  %s\n"+
		"  version:  %s\n"+
		"  access: %s\n"+
		"  constants:\n%s\n",
		filename,
		cls.ClassName(),
		cls.Version(),
		cls.StringFlags(),
		cls.StringConstantsIndent(4))
}

func main() {
	flag.Parse()

	for _, source := range flag.Args() {
		cls, err := inspectFilename(source)
		if err != nil {
			fmt.Printf("Can't inspect '%s': %s\n", source, err)
			os.Exit(1)
		}

		printClass(source, cls)
	}
}
