package pool

import (
	"fmt"
	"testing"
)

func TestNewPool(t *testing.T) {
	f := func(args ...any) {
		for _, arg := range args {
			fmt.Println("Function called with:", arg)
		}
	}
	pool := NewPool[bool, func(bool)](1, 1)
	fmt.Println(pool)
	pool.addJob(f, true, false)
	//pool.JobQueue <- (f, true)
	//pool.JobQueue <- false
	//pool.JobQueue <- true
	pool.WaitAll()
}
