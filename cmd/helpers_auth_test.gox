package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1_randomString(t *testing.T) {
	l := 12
	r1 := randomString(l)
	log.Println(r1)
	r2 := randomString(l)
	log.Println(r2)
	assert.EqualValues(t, len(r1), l, "Random string length issue.")
}
