package pool

import (
	"fmt"
	"sync"

	workers "vk-go-test/worker"
)

// WorkerPool manages a collection of workers and jobs distribution
type WorkerPool struct {
	Workers    map[int]*workers.Worker
	JobChannel chan string
	Wg         sync.WaitGroup
	NextID     int
	Mu         sync.Mutex
}

// NewWorkerPool creates new pool with initital workers
func NewWorkerPool(initialWorkers int) *WorkerPool {
	pool := &WorkerPool{
		Workers:    make(map[int]*workers.Worker),
		JobChannel: make(chan string, 100), // buffer for 100 jobs
		NextID:     1,
	}

	for range initialWorkers {
		pool.AddWorker()
	}

	return pool
}

// AddWorker creates and starts a new worker
func (p *WorkerPool) AddWorker() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	id := p.NextID
	p.NextID++

	worker := workers.NewWorker(id, p.JobChannel, &p.Wg)
	p.Workers[id] = worker
	worker.Start()

	fmt.Printf("Worker %d added\n> ", id)
}

// RemoveWorker stops and removes a worker
func (p *WorkerPool) RemoveWorker() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if len(p.Workers) == 0 {
		fmt.Println("No workers to remove")
		return
	}

	for id, worker := range p.Workers {
		worker.Stop()
		delete(p.Workers, id)
		fmt.Printf("Worker %d removed\n> ", id)
		return
	}
}

// SubmitJob sends job to the worker pool
func (p *WorkerPool) SubmitJob(job string) {
	p.JobChannel <- job
}

// WorkerCount return current number of workers
func (p *WorkerPool) WorkerCount() int {
	p.Mu.Lock()
	defer p.Mu.Unlock()
	return len(p.Workers)
}

// PrintStatus display current tasks assignment
func (p *WorkerPool) PrintStatus() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if len(p.Workers) == 0 {
		fmt.Println("No active workers")
		return
	}

	fmt.Println("Current worker tasks:")
	for id, worker := range p.Workers {
		if worker.CurrentTask != "" {
			fmt.Printf("Worker %d: %s\n", id, worker.CurrentTask)
		} else {
			fmt.Printf("Worker %d: idle\n", id)
		}
	}
}

// Stop all workers for exiting
func (p *WorkerPool) Stop() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for _, worker := range p.Workers {
		worker.Stop()
	}

	p.Wg.Wait()
	close(p.JobChannel)
}
