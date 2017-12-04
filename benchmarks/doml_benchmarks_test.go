/*
 * Copyright (c) 2017. Braedon Wooding
 * Created under LICENSE, see the file LICENSE for information
 */

package benchmarks

import (
	"doml/parser"
	"doml/runtime"
	"testing"
)

var smallTestString = `
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
`

var reasonableTestString = `
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
`

var largeTestString = `
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
@ example =System.Cool...
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
.x = 2, 4, 9, "bob" ; example.x = 4949, 0x03939 ;
.b = 2
`

func BenchmarkPush(b *testing.B) {
	runtimeInstance := runtime.NewRuntime(1, 0, nil)

	for i := 0; i < b.N; i++ {
		runtimeInstance.Push(2, false)
		runtimeInstance.Unsafe.SetStackPtr(-1)
	}
}

// Effectively going to be equivalent to speed of push, just without the set stack ptr
func BenchmarkPeek(b *testing.B) {
	runtimeInstance := runtime.NewRuntime(1, 0, nil)
	runtimeInstance.Push(2, false)

	for i := 0; i < b.N; i++ {
		runtimeInstance.Peek()
	}
}

func BenchmarkParserWithSmallString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseDOMLFromString(smallTestString)
	}
}

func BenchmarkParserWithReasonableString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseDOMLFromString(reasonableTestString)
	}
}

func BenchmarkParserWithLargeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parser.ParseDOMLFromString(largeTestString)
	}
}
