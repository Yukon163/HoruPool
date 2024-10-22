package pool

import (
	"runtime"
	"sync"
)

// modified from https://github.com/ivpusic/grpool/blob/master/grpool.go

//type job[A any, F any] struct {
//	fun   F
//	argus []A
//}

type JobFunc interface {
	Call(args ...interface{})
}

type funcAdapter struct {
	f func(...interface{})
}

type job struct {
	fun    JobFunc
	params []interface{} // 使用切片存储变长参数
}

func (fa funcAdapter) Call(args ...interface{}) {
	fa.f(args...)
}

type worker[A any, F any] struct {
	workerPool chan *worker[A, F]
	jobChannel chan job[A, F]
	stop       chan struct{}
}

func (w *worker[A, F]) start() {
	go func(w *worker[A, F]) {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		var j job[A, F]
		for {
			// worker free, add it to pool
			w.workerPool <- w

			select {
			case j = <-w.jobChannel:
				//j.fun.Call()
				runtime.Gosched()
			case <-w.stop:
				w.stop <- struct{}{}
				return
			}
		}
	}(w)
}

func newWorker[A any, F any](pool chan *worker[A, F]) *worker[A, F] {
	return &worker[A, F]{
		workerPool: pool,
		jobChannel: make(chan job[A, F]),
		stop:       make(chan struct{}),
	}
}

// Accepts jobs from clients, and waits for first free worker to deliver job
type dispatcher[A any, F any] struct {
	workerPool chan *worker[A, F]
	jobQueue   chan job[A, F]
	stop       chan struct{}
}

func (d *dispatcher[A, F]) dispatch() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for {
		select {
		case job := <-d.jobQueue:
			worker := <-d.workerPool
			runtime.Gosched()
			worker.jobChannel <- job
			runtime.Gosched()
		case <-d.stop:
			for i := 0; i < cap(d.workerPool); i++ {
				worker := <-d.workerPool
				runtime.Gosched()
				worker.stop <- struct{}{}
				runtime.Gosched()
				<-worker.stop
				runtime.Gosched()
			}

			d.stop <- struct{}{}
			return
		}
	}
}

func newDispatcher[A any, F any](workerPool chan *worker[A, F], jobQueue chan job[A, F]) *dispatcher[A, F] {
	d := &dispatcher[A, F]{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		stop:       make(chan struct{}),
	}

	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool)
		worker.start()
	}

	go d.dispatch()
	return d
}

type Pool[A any, F any] struct {
	JobQueue   chan job[A, F]
	dispatcher *dispatcher[A, F]
	wg         sync.WaitGroup
}

func NewPool[A any, F any](numWorkers int, jobQueueLen int) *Pool[A, F] {
	jobQueue := make(chan job[A, F], jobQueueLen)
	workerPool := make(chan *worker[A, F], numWorkers)

	pool := &Pool[A, F]{
		JobQueue:   jobQueue,
		dispatcher: newDispatcher[A, F](workerPool, jobQueue),
	}
	return pool
}

func (p *Pool[A, F]) addJob(fun F, argus ...A) {
	p.JobQueue <- job[A, F]{fun, argus}
}

// In case you are using WaitAll fn, you should call this method
// every time your job is done.

// JobDone If you are not using WaitAll then we assume you have your own way of synchronizing.
func (p *Pool[A, F]) JobDone() {
	p.wg.Done()
}

// WaitCount How many jobs we should wait when calling WaitAll.
// It is using WaitGroup Add/Done/Wait
func (p *Pool[A, F]) WaitCount(count int) {
	p.wg.Add(count)
}

// WaitAll Will wait for all jobs to finish.
func (p *Pool[A, F]) WaitAll() {
	p.wg.Wait()
}

// Release Will release resources used by pool
func (p *Pool[A, F]) Release() {
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}
