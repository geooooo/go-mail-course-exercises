package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func init() {
	SlowSearch(ioutil.Discard)
	FastSearch(ioutil.Discard)
}

func TestSearch(t *testing.T) {
	slowOut := new(bytes.Buffer)
	SlowSearch(slowOut)
	slowResult := slowOut.String()

	fastOut := new(bytes.Buffer)
	FastSearch(fastOut)
	fastResult := fastOut.String()

	if slowResult != fastResult {
		t.Errorf("results not match\nGot:\n%v\nExpected:\n%v", fastResult, slowResult)
	}
}