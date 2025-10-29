package main

import (
	"fmt"
	"os"
)

/*
cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
BenchmarkSlow-8               50          23393609 ns/op        20349353 B/op     182848 allocs/op
BenchmarkFast-8              165           7139242 ns/op         2714035 B/op      47481 allocs/op
*/
func main() {
	if len(os.Args) > 1 && os.Args[1] == "-f" {
		fmt.Println("Fast")
		FastSearch(os.Stdout)
	} else {
		fmt.Println("Slow")
		SlowSearch(os.Stdout)
	}
}
