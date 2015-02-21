package main

// most of the code here is based on informations found here:
//  - https://en.wikipedia.org/wiki/Java_class_file
//  - http://docs.oracle.com/javase/specs/jvms/se7/html/jvms-4.html#jvms-4.4

import (
	"flag"
	"fmt"
	"github.com/bfontaine/goclassy/jvm"
	"os"
)

func inspectFilename(source string) (jvm.JClass, error) {
	f, err := os.Open(source)
	if err != nil {
		return jvm.JClass{}, err
	}

	defer f.Close()

	return jvm.ParseClassFromFile(f)
}

func printClass(filename string, cls jvm.JClass) {
	fmt.Printf("%s:\n"+
		"  class: %s\n"+
		"  version: %s\n"+
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
