package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type jclass struct {
	majorVersion int
	minorVersion int
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

func readInt(r io.Reader, data *int, length int) error {
	buf := make([]byte, length)
	intVal := 0
	if err := binary.Read(r, binary.BigEndian, buf); err != nil {
		return err
	}

	for _, b := range buf {
		intVal <<= 8
		intVal += int(b)
	}

	*data = intVal
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

	return jc, nil
}

func (jc *jclass) Version() string {
	// magic numbers source:
	// https://en.wikipedia.org/wiki/Java_class_file

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

func printClass(filename string, jc jclass) {
	fmt.Printf("%s:\n  %s\n", filename, jc.Version())
}

func main() {
	flag.Parse()

	for _, source := range flag.Args() {
		jc, err := inspectFilename(source)
		if err != nil {
			fmt.Printf("Can't inspect '%s': %s\n", source, err)
		}

		printClass(source, jc)
	}
}
