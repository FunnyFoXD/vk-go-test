package pool

import (
	"fmt"
	"sync"

	workers "vk-go-test/worker"
)

// The structure represents a WorkerPool and contains information about it
type WorkerPool struct {
	Workers    map[int]*workers.Worker
	JobChannel chan string
	Wg         sync.WaitGroup
	NextID     int
	Mu         sync.Mutex
}

// NewWorkerPool create new PoolWorker with initialWorkers (Fabric)
func NewWorkerPool(initialWorkers int) *WorkerPool {
	pool := &WorkerPool{
		Workers:    make(map[int]*workers.Worker),
		JobChannel: make(chan string, 100),
		NextID:     1,
	}

	// Initial Pool
	for range initialWorkers {
		pool.AddWorker()
	}

	return pool
}

// AddWorker add new Worker in WorkerPool
func (p *WorkerPool) AddWorker() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	id := p.NextID
	p.NextID++

	worker := workers.NewWorker(id, p.JobChannel, &p.Wg)
	p.Workers[id] = worker
	worker.Start()

	fmt.Printf("Worker %d is added\n", id)
}

// RemoveWorker remove Worker from WorkerPool
func (p *WorkerPool) RemoveWorker() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if len(p.Workers) == 0 {
		fmt.Println("No workers to remove!")
		return
	}

	for id, worker := range p.Workers {
		worker.Stop()
		delete(p.Workers, id)

		fmt.Printf("Worker %d is deleted\n", id)
		return
	}
}

// SubmitJob submit job into channel
func (p *WorkerPool) SubmitJob(job string) {
	p.JobChannel <- job
}

// WorkerCount return number of workers in WorkerPool
func (p *WorkerPool) WorkerCount() int {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	return len(p.Workers)
}

// Stop stop WorkerPool if it has no jobs
func (p *WorkerPool) Stop() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for _, worker := range p.Workers {
		worker.Stop()
	}

	// WAIT!
	p.Wg.Wait()

	close(p.JobChannel)
}
