package pool

import (
	"runtime"
	"sync"
)

// modified from https://github.com/ivpusic/grpool/blob/master/grpool.go

type worker[A interface{}, F func(A)] struct {
	workerPool chan *worker[A, F]
	jobChannel chan A
	stop       chan struct{}
	fun        F
}

func (w *worker[A, F]) start() {
	go func(w *worker[A, F]) {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		var job A
		for {
			// worker free, add it to pool
			w.workerPool <- w

			select {
			case job = <-w.jobChannel:
				w.fun(job)
				runtime.Gosched()
			case <-w.stop:
				w.stop <- struct{}{}
				return
			}
		}
	}(w)
}

func newWorker[A interface{}, F func(A)](pool chan *worker[A, F], fun F) *worker[A, F] {
	return &worker[A, F]{
		workerPool: pool,
		jobChannel: make(chan A),
		stop:       make(chan struct{}),
		fun:        fun,
	}
}

// Accepts jobs from clients, and waits for first free worker to deliver job
type dispatcher[A interface{}, F func(A)] struct {
	workerPool chan *worker[A, F]
	jobQueue   chan A
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

func newDispatcher[A interface{}, F func(A)](workerPool chan *worker[A, F], jobQueue chan A, fun F) *dispatcher[A, F] {
	d := &dispatcher[A, F]{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		stop:       make(chan struct{}),
	}

	for i := 0; i < cap(d.workerPool); i++ {
		worker := newWorker(d.workerPool, fun)
		worker.start()
	}

	go d.dispatch()
	return d
}

type Pool[A interface{}, F func(A)] struct {
	JobQueue   chan A
	dispatcher *dispatcher[A, F]
	wg         sync.WaitGroup
}

func NewPool[A interface{}, F func(A)](numWorkers int, jobQueueLen int, fun F) *Pool[A, F] {
	jobQueue := make(chan A, jobQueueLen)
	workerPool := make(chan *worker[A, F], numWorkers)

	pool := &Pool[A, F]{
		JobQueue:   jobQueue,
		dispatcher: newDispatcher[A, F](workerPool, jobQueue, fun),
	}

	return pool
}

// In case you are using WaitAll fn, you should call this method
// every time your job is done.
//
// If you are not using WaitAll then we assume you have your own way of synchronizing.
func (p *Pool[A, F]) JobDone() {
	p.wg.Done()
}

// How many jobs we should wait when calling WaitAll.
// It is using WaitGroup Add/Done/Wait
func (p *Pool[A, F]) WaitCount(count int) {
	p.wg.Add(count)
}

// Will wait for all jobs to finish.
func (p *Pool[A, F]) WaitAll() {
	p.wg.Wait()
}

// Will release resources used by pool
func (p *Pool[A, F]) Release() {
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}
