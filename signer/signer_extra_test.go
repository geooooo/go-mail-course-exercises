package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestExtra(t *testing.T) {
	var recieved uint32
	freeFlowJobs := []Job{
		func(_, out chan any) {
			defer close(out)
			out <- uint32(1)
			out <- uint32(3)
			out <- uint32(4)
		},
		func(in, out chan any) {
			defer close(out)
			for val := range in {
				out <- val.(uint32) * 3
				time.Sleep(time.Millisecond * 100)
			}
		},
		func(in, _ chan any) {
			for val := range in {
				fmt.Println("collected", val)
				atomic.AddUint32(&recieved, val.(uint32))
			}
		},
	}

	start := time.Now()

	ExecutePipeline(freeFlowJobs...)

	end := time.Since(start)
	expectedTime := time.Millisecond * 350

	if end > expectedTime {
		t.Errorf("execition too long\nGot: %s\nExpected: <%s", end, expectedTime)
	}

	if recieved != (1+3+4)*3 {
		t.Errorf("f3 have not collected inputs, recieved = %d", recieved)
	}
}
