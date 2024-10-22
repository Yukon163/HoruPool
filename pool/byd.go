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

type job struct {
	fun   jobFunc
	argus []interface{}
}

type jobFunc interface {
	Call(args ...interface{})
}

type funcAdapter struct {
	fun func(...interface{})
}

func (fa funcAdapter) Call(args ...interface{}) {
	fa.fun(args...)
}

type worker struct {
	workerPool chan *worker
	jobChannel chan job
	stop       chan struct{}
}

func (w *worker) start() {
	go func(w *worker) {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		var j job
		for {
			// worker free, add it to pool
			w.workerPool <- w

			select {
			case j = <-w.jobChannel:
				//j.fun.Call()
				j.fun.Call(j.argus...)
				runtime.Gosched()
			case <-w.stop:
				w.stop <- struct{}{}
				return
			}
		}
	}(w)
}

func newWorker(pool chan *worker) *worker {
	return &worker{
		workerPool: pool,
		jobChannel: make(chan job),
		stop:       make(chan struct{}),
	}
}

// Accepts jobs from clients, and waits for first free worker to deliver job
type dispatcher struct {
	workerPool chan *worker
	jobQueue   chan job
	stop       chan struct{}
}

func (d *dispatcher) dispatch() {
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

func newDispatcher(workerPool chan *worker, jobQueue chan job) *dispatcher {
	d := &dispatcher{
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

type Pool struct {
	JobQueue   chan job
	dispatcher *dispatcher
	wg         sync.WaitGroup
}

func NewPool(numWorkers int, jobQueueLen int) *Pool {
	jobQueue := make(chan job, jobQueueLen)
	workerPool := make(chan *worker, numWorkers)

	pool := &Pool{
		JobQueue:   jobQueue,
		dispatcher: newDispatcher(workerPool, jobQueue),
	}
	return pool
}

func (p *Pool) addJob(fun func(...interface{}), argus ...interface{}) {
	p.JobQueue <- job{fun: funcAdapter{fun: fun}, argus: argus}
}

// In case you are using WaitAll fn, you should call this method
// every time your job is done.

// JobDone If you are not using WaitAll then we assume you have your own way of synchronizing.
func (p *Pool) JobDone() {
	p.wg.Done()
}

// WaitCount How many jobs we should wait when calling WaitAll.
// It is using WaitGroup Add/Done/Wait
func (p *Pool) WaitCount(count int) {
	p.wg.Add(count)
}

// WaitAll Will wait for all jobs to finish.
func (p *Pool) WaitAll() {
	p.wg.Wait()
}

// Release Will release resources used by pool
func (p *Pool) Release() {
	p.dispatcher.stop <- struct{}{}
	<-p.dispatcher.stop
}
