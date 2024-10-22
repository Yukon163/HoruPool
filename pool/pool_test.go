package pool

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPool(t *testing.T) {
	f := func(args ...any) {
		for _, arg := range args {
			fmt.Println("Function called with:", arg)
		}
	}

	//f1 := func(arg1 int, arg2 string) {
	//	fmt.Printf("Function horued with: %d, %s", arg1, arg2)
	//}
	pool := NewPool(1, 1)
	fmt.Println(pool)
	pool.addJob(f, true, false)
	//pool.addJob(f1, 1, false)
	//pool.JobQueue <- (fun, true)
	//pool.JobQueue <- false
	//pool.JobQueue <- true
	//pool.WaitAll()
	time.Sleep(50 * time.Millisecond)
}
