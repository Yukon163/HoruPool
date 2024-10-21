package pool

import (
	"fmt"
	"testing"
)

func TestNewPool(t *testing.T) {
	f := func(boo bool) {
		fmt.Println(boo)
	}
	pool := NewPool[bool, func(bool)](1, 1, f)
	fmt.Println(pool)
	pool.JobQueue <- true
	pool.JobQueue <- false
	pool.JobQueue <- true
	pool.WaitAll()
}
