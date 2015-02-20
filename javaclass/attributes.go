package javaclass

import ()

type JAttr struct {
	nameIndex u2
	info      []u1
}

func (attr *JAttr) initInfo(size u4) {
	attr.info = make([]u1, size)
}
