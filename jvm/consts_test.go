package jvm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmptyValueAsString(t *testing.T) {
	cst := JConst{}

	assert.Equal(t, cst.valueAsString(), "")
}

func TestValueAsString(t *testing.T) {
	cst := JConst{value: []byte{'h', 'e', 'l', 'l', 'o'}}

	assert.Equal(t, cst.valueAsString(), "hello")
}

func TestDumpValueNilArg(t *testing.T) {
	cst := JConst{value: []byte{'a'}}
	err := cst.dumpValue(nil)

	assert.NotNil(t, err)
}

func TestDumpValueTooSmallArgType(t *testing.T) {
	cst := JConst{value: []byte{1, 1, 1, 1}}
	var res int16
	err := cst.dumpValue(&res)

	assert.NotNil(t, err)
}
