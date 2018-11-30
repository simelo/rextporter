package scrapper

import (
	"log"
	"sync"
	"time"
)

type ScrapResult struct {
	Val               interface{}
	ConstMetricIdxOut int
}

type ScrapErrResult struct {
	Err               error
	ConstMetricIdxOut int
}

type ScrapRequest struct {
	Scrap            Scrapper
	Res              chan ScrapResult
	ConstMetricIdxIn int
	Err              chan ScrapErrResult
}

type scrapWork struct {
	scrapper         Scrapper
	res              chan ScrapResult
	constMetricIdxIn int
	err              chan ScrapErrResult
}

type workQueue chan scrapWork
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
		works:    make(chan scrapWork, workersNum*2),
		nWorkers: workersNum,
		wg:       sync.WaitGroup{},
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Apply(ri ScrapRequest) {
	work := scrapWork{
		scrapper:         ri.Scrap,
		res:              ri.Res,
		constMetricIdxIn: ri.ConstMetricIdxIn,
		err:              ri.Err,
	}
	p.works <- work
}

type workerT struct {
	works    chan scrapWork
	workers  chan workQueue
	quitChan chan bool
	wg       *sync.WaitGroup
}

func (p *Pool) newWorker() workerT {
	log.Println("creating worker")
	return workerT{
		works:    make(chan scrapWork),
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
				log.Println("work.client", work.scrapper)
				val, err := work.scrapper.GetMetric()
				if err == nil {
					work.res <- ScrapResult{Val: val, ConstMetricIdxOut: work.constMetricIdxIn}
				} else {
					work.err <- ScrapErrResult{Err: err, ConstMetricIdxOut: work.constMetricIdxIn}
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
