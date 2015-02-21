package main

import (
	"github.com/bfontaine/goclassy/jvm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringFlagsEmptyClass(t *testing.T) {
	jc := jvm.JClass{}

	assert.Equal(t, stringFlags(&jc), "")
}
