package main

import (
	"fmt"
)

func main() {
	nums := [...]int{1, 2, 3}

	ExecutePipeline([]Job{
		func(_, out chan any) {
			defer close(out)
			for _, n := range nums {
				fmt.Printf("1) send: %d\n", n)
				out <- n
			}
		},
		func(in, out chan any) {
			defer close(out)
			for v := range in {
				vv := v.(int)
				r := vv * vv

				fmt.Printf("2) send: %d\n", r)
				out <- r
			}
		},
		func(in, _ chan any) {
			s := 0
			for v := range in {
				vv := v.(int)
				s += vv
			}
			fmt.Printf("3) result: %d\n", s)
		},
	}...)
}
