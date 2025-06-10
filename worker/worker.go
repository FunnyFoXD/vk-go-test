package worker

import (
	"sync"
	"time"
)

// Worker process
type Worker struct {
	ID          int
	JobChannel  chan string
	QuitChannel chan struct{}
	Wg          *sync.WaitGroup
	CurrentTask string
	mu          sync.Mutex
}

// NewWorker creates a new worker instance
func NewWorker(id int, jobChannel chan string, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:          id,
		JobChannel:  jobChannel,
		QuitChannel: make(chan struct{}),
		Wg:          wg,
	}
}

// Start begin processing jobs in goroutine
func (w *Worker) Start() {
	w.Wg.Add(1)

	go func() {
		defer w.Wg.Done()

		for {
			select {
			case job := <-w.JobChannel:
				// Update current task
				w.mu.Lock()
				w.CurrentTask = job
				w.mu.Unlock()

				// Work imitation
				time.Sleep(4 * time.Second)

				// Clear task
				w.mu.Lock()
				w.CurrentTask = ""
				w.mu.Unlock()

			case <-w.QuitChannel:
				return // Exit goroutine
			}
		}
	}()
}

// Stop signals the worker to terminate
func (w *Worker) Stop() {
	close(w.QuitChannel)
}
