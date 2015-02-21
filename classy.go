package main

// most of the code is based on informations found here:
//  - https://en.wikipedia.org/wiki/Java_class_file
//  - http://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/bfontaine/classy/jvm"
	"os"
	"strings"
)

func inspectFilename(source string) (jvm.JClass, error) {
	f, err := os.Open(source)
	if err != nil {
		return jvm.JClass{}, err
	}

	defer f.Close()

	return jvm.ParseClassFromFile(f)
}

func stringConstantsIndent(cls jvm.JClass, indent int) string {
	var buffer bytes.Buffer

	lineStart := strings.Repeat(" ", indent)

	for idx, cst := range cls.Constants() {
		if idx == 0 {
			continue
		}

		buffer.WriteString(fmt.Sprintf("%s%3d = %v\n", lineStart, idx, cst))
	}

	return buffer.String()
}

func stringFlags(cls *jvm.JClass) string {
	var buffer bytes.Buffer

	if cls.HasAccessFlag(jvm.ACC_PUBLIC) {
		buffer.WriteString("public ")
	}
	if cls.HasAccessFlag(jvm.ACC_FINAL) {
		buffer.WriteString("final ")
	}
	if cls.HasAccessFlag(jvm.ACC_INTERFACE) {
		buffer.WriteString("interface ")
	}
	if cls.HasAccessFlag(jvm.ACC_ABSTRACT) {
		buffer.WriteString("abstract ")
	}

	return buffer.String()
}

func printClass(filename string, cls jvm.JClass) {
	fmt.Printf("%s:\n"+
		"  class: %s\n"+
		"  version: %s\n"+
		"  access: %s\n"+
		"  constants:\n%s\n",
		filename,
		cls.ClassName(),
		cls.JavaVersion(),
		stringFlags(&cls),
		stringConstantsIndent(cls, 4))
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
