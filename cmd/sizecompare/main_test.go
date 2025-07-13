package main

import (
	"testing"
)

func BenchmarkSerializeToJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		serializeToJSON(metadata)
	}
}

func BenchmarkSerializeToXML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		serializeToXML(metadata)
	}
}

func BenchmarkSerializeToProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		serializeToProto(genMetadata)
	}
}

/*
Output
goos: darwin
goarch: arm64
pkg: github.com/akkahshh24/movieapp/cmd/sizecompare
cpu: Apple M1
BenchmarkSerializeToJSON-8       4100863               283.8 ns/op
BenchmarkSerializeToXML-8         587469              1871 ns/op
BenchmarkSerializeToProto-8      7711086               155.1 ns/op
PASS
ok      github.com/akkahshh24/movieapp/cmd/sizecompare  4.242s
*/
