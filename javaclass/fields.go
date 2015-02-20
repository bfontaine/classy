package javaclass

import ()

type JField struct {
	accessFlags     u2
	nameIndex       u2
	descriptorIndex u2
	attributes      []JAttr
}

func (f *JField) initAttributes(attrsSize u2) {
	f.attributes = make([]JAttr, attrsSize)
}

func (f *JField) addAttribute(index u2, attr JAttr) {
	f.attributes[index] = attr
}
