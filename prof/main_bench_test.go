package main

import (
	"io/ioutil"
	"testing"
)

func init() {
	SlowSearch(ioutil.Discard)
	FastSearch(ioutil.Discard)
}

func BenchmarkSlow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SlowSearch(ioutil.Discard)
	}
}

func BenchmarkFast(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FastSearch(ioutil.Discard)
	}
}
