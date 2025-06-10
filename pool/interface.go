package pool

// Action For WorkerPool
type WorkerPoolAction interface {
	AddWorker()
	RemoveWorker()
	SubmitJob(job string)
	WorkerCount() int
	PrintStatus()
	Stop()
}
