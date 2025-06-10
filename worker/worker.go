package worker

import (
	"fmt"
	"sync"
	"time"
)

// The structure represents a Worker and contains information about it
type Worker struct {
	ID          int
	JobChannel  chan string
	QuitChannel chan struct{}
	Wg          *sync.WaitGroup
}

// NewWorker create new Worker (Fabric)
func NewWorker(id int, jobChannel chan string, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:          id,
		JobChannel:  jobChannel,
		QuitChannel: make(chan struct{}),
		Wg:          wg,
	}
}

// Start start worker and use imitation for work
func (w *Worker) Start() {
	w.Wg.Add(1)

	go func() {
		defer w.Wg.Done()

		for {
			select {
			case job := <-w.JobChannel:
				fmt.Printf("Worker %d work with %s\n", w.ID, job)
				time.Sleep(100 * time.Millisecond) // Just imitation for worker
			case <-w.QuitChannel:
				fmt.Printf("Worker %d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop stop worker
func (w *Worker) Stop() {
	close(w.QuitChannel)
}
