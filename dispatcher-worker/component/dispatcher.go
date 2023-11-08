package component

import "sync"

type Dispatcher struct {
	JobQueue   chan Job
	wg         sync.WaitGroup
	maxWorkers int
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	jobQueue := make(chan Job)
	return &Dispatcher{JobQueue: jobQueue, maxWorkers: maxWorkers}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.JobQueue, &d.wg)
		worker.Start()
	}
}

func (d *Dispatcher) Add(payload Payload) {
	job := Job{Payload: payload}
	d.wg.Add(1)
	d.JobQueue <- job
}

func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

func (d *Dispatcher) Stop() {
	close(d.JobQueue)
}
