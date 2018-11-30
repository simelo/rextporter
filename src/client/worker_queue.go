package client

import (
	"log"
	"sync"
	"time"
)

type RequestInfo struct {
	Client Client
	Res    chan []byte
	Err    chan error
}

type workRequest struct {
	client Client
	res    chan []byte
	err    chan error
}

type workQueue chan workRequest
type workerQueue chan workQueue

type Pool struct {
	workers  workerQueue
	works    workQueue
	nWorkers uint
	wg       sync.WaitGroup
}

func NewPool(workersNum uint) *Pool {
	return &Pool{
		workers:  make(chan workQueue, workersNum),
		works:    make(chan workRequest, workersNum*2),
		nWorkers: workersNum,
		wg:       sync.WaitGroup{},
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Apply(ri RequestInfo) {
	work := workRequest{client: ri.Client, res: ri.Res, err: ri.Err}
	p.works <- work
}

type workerT struct {
	works    chan workRequest
	workers  chan workQueue
	quitChan chan bool
	wg       *sync.WaitGroup
}

func (p *Pool) newWorker() workerT {
	log.Println("creating worker")
	return workerT{
		works:    make(chan workRequest),
		workers:  p.workers,
		quitChan: make(chan bool),
		wg:       &p.wg,
	}
}

func (w *workerT) start() {
	log.Println("starting worker")
	w.wg.Add(1)
	go func() {
		for {
			// Add ourselves as available into the worker queue.
			w.workers <- w.works
			select {
			// Wait for a work request.
			case work := <-w.works:
				log.Println("start processing work")
				time.Sleep(time.Second * 1)
				log.Println("work.client", work.client)
				data, err := work.client.GetData()
				if err == nil {
					work.res <- data
				} else {
					work.err <- err
				}
				log.Println("finish processing work")
				// wait for a quit msg
			case <-w.quitChan:
				// Receive a close worker request.
				w.wg.Done()
				return
			}
		}
	}()
}

func (w *workerT) stop() {
	go func() {
		w.quitChan <- true
		// close(w.quitChan)
	}()
}

func (p *Pool) Stop() {
	// for {
	// 	worker := <-p.workers
	// 	worker.stop()
	// }
}

func (p *Pool) StartDispatcher() {
	for i := uint(0); i < p.nWorkers; i++ {
		worker := p.newWorker()
		worker.start()
	}
	go func() {
		for {
			select {
			// wait for an incoming work
			case work := <-p.works:
				log.Println("work recived")
				// wait for an available worker
				worker := <-p.workers
				log.Println("worker available")
				// dispatch work into worker
				worker <- work
			}
		}
	}()
}
