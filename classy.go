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
)

type jclass struct {
	majorVersion int
	minorVersion int

	publicFlag    bool
	finalFlag     bool
	superFlag     bool
	interfaceFlag bool
	abstractFlag  bool
}

func readBytes(f *os.File, buff []byte) error {
	n, err := f.Read(buff)

	if err != nil {
		return err
	}

	if n < len(buff) {
		return errors.New("Can't read enough bytes")
	}

	return nil
}

func readBinary(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.BigEndian, data)
}

func readInt(r io.Reader, data *int, length int) error {
	buf := make([]byte, length)
	intVal := 0
	if err := readBinary(r, buf); err != nil {
		return err
	}

	for _, b := range buf {
		intVal <<= 8
		intVal += int(b)
	}

	*data = intVal
	return nil
}

func (jc *jclass) addConstant(tag int, size int, data []byte) {
	// TODO
}

func (jc *jclass) parseConstantPool(constantPoolSize int, r io.Reader) error {

	var tag, size int
	var data []byte

	// skip the first index
	if err := readInt(r, &tag, 2); err != nil {
		return err
	}

	for index := 1; index < constantPoolSize; index++ {
		if err := readInt(r, &tag, 1); err != nil {
			return err
		}

		switch tag {
		case TAG_STRING:
			if err := readInt(r, &size, 2); err != nil {
				return err
			}
			break

		case TAG_CLASS_REF:
		case TAG_STRING_REF:
		case TAG_METHOD_TYPE:
			size = 2
			break

		case TAG_METHOD_HANDLE:
			size = 3
			break

		case TAG_INT:
		case TAG_FLOAT:
		case TAG_FIELD_REF:
		case TAG_METHOD_REF:
		case TAG_INTERFACE_METHOD_REF:
		case TAG_NAME_TYPE_DESC:
		case TAG_INVOKE_DYN:
			size = 4
			break

		case TAG_LONG:
		case TAG_DOUBLE:
			size = 8
			break

		default:
			return errors.New(fmt.Sprintf("Unknown tag '%d'", tag))
		}

		data = make([]byte, size)
		if err := readBinary(r, &data); err != nil {
			return err
		}
		jc.addConstant(tag, size, data)
	}

	return nil
}

func parseInterfaces(jc *jclass, interfaces []byte) error {
	// TODO
	return nil
}

func parseFields(jc *jclass, fields []byte) error {
	// TODO
	return nil
}

func parseMethods(jc *jclass, methods []byte) error {
	// TODO
	return nil
}

func parseAttrs(jc *jclass, attrs []byte) error {
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
		return jclass{}, errors.New("Bad magic number")
	}

	jc := jclass{}

	// minor version number
	if err := readInt(f, &jc.minorVersion, 2); err != nil {
		return jc, err
	}

	// major version number
	if err := readInt(f, &jc.majorVersion, 2); err != nil {
		return jc, err
	}

	// constant pool size
	var constantPoolSize int
	if err := readInt(f, &constantPoolSize, 2); err != nil {
		return jc, err
	}

	// constant pool
	if err := jc.parseConstantPool(constantPoolSize, f); err != nil {
		return jc, err
	}

	// access flags
	var accessFlags int
	if err := readInt(f, &accessFlags, 2); err != nil {
		return jc, err
	}
	jc.SetAccessFlags(accessFlags)

	/* TODO

	// this class
	var classIndex int
	if err := readInt(f, &classIndex, 2); err != nil {
		return jc, err
	}
	// TODO

	// super class
	var superClassIndex int
	if err := readInt(f, &superClassIndex, 2); err != nil {
		return jc, err
	}
	// TODO

	// interfaces
	var interfacesCount int
	if err := readInt(f, &interfacesCount, 2); err != nil {
		return jc, err
	}

	interfaces := make([]byte, interfacesCount)
	if err := parseInterfaces(&jc, interfaces); err != nil {
		return jc, err
	}

	// fields
	var fieldsCount int
	if err := readInt(f, &fieldsCount, 2); err != nil {
		return jc, err
	}

	fields := make([]byte, fieldsCount)
	if err := parseFields(&jc, fields); err != nil {
		return jc, err
	}

	// methods
	var methodsCount int
	if err := readInt(f, &methodsCount, 2); err != nil {
		return jc, err
	}

	methods := make([]byte, methodsCount)
	if err := parseMethods(&jc, methods); err != nil {
		return jc, err
	}

	// attributes
	var attrsCount int
	if err := readInt(f, &attrsCount, 2); err != nil {
		return jc, err
	}

	attrs := make([]byte, attrsCount)
	if err := parseAttrs(&jc, attrs); err != nil {
		return jc, err
	}
	*/

	return jc, nil
}

func setFlag(n int, magic int, flag *bool) {
	*flag = n&magic == magic
}

func (jc *jclass) SetAccessFlags(accessFlags int) {
	setFlag(accessFlags, 0x0001, &jc.publicFlag)
	setFlag(accessFlags, 0x0010, &jc.finalFlag)
	setFlag(accessFlags, 0x0020, &jc.superFlag)
	setFlag(accessFlags, 0x0200, &jc.interfaceFlag)
	setFlag(accessFlags, 0x0400, &jc.abstractFlag)
}

func (jc *jclass) Version() string {
	// we don't use the minor version here
	switch jc.majorVersion {
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

func (jc *jclass) StringFlags() string {
	var buffer bytes.Buffer

	if jc.publicFlag {
		buffer.WriteString("public ")
	}
	if jc.finalFlag {
		buffer.WriteString("final ")
	}
	if jc.interfaceFlag {
		buffer.WriteString("interface ")
	}
	if jc.abstractFlag {
		buffer.WriteString("abstract ")
	}

	return buffer.String()
}

func printClass(filename string, jc jclass) {
	fmt.Printf("%s:\n"+
		"  version:  %s\n"+
		"  access: %s\n",
		filename, jc.Version(), jc.StringFlags())
}

func main() {
	flag.Parse()

	for _, source := range flag.Args() {
		jc, err := inspectFilename(source)
		if err != nil {
			fmt.Printf("Can't inspect '%s': %s\n", source, err)
			os.Exit(1)
		}

		printClass(source, jc)
	}
}
