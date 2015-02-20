package javaclass

import ()

type JField struct {
	accessFlags     u2
	nameIndex       u2
	descriptorIndex u2
	attributes      []JAttr
}
