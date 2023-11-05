package component

import "sync"

// Job represents the job to be run
type Job struct {
	Payload Payload
}

type Payload struct {
	Message string
}

// Worker represents the worker that executes the job
type Worker struct {
	JobQueue chan Job
	wg       *sync.WaitGroup
	quit     chan bool
}

func NewWorker(jobQueue chan Job, wg *sync.WaitGroup) Worker {
	return Worker{JobQueue: jobQueue, wg: wg, quit: make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.JobQueue:
				// do the job
				println(job.Payload.Message)
				w.wg.Done()
			case <-w.quit:
				return
			}
		}
	}()
}
