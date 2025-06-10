package worker

// Actions for Worker
type WorkerAction interface {
	Start()
	Stop()
}
