package workshop

import "context"

type Worker struct {
	id   int
	pipe <-chan Job
}

func NewWorker(id int, pipe <-chan Job) (w *Worker) {
	w = &Worker{
		id:   id,
		pipe: pipe,
	}

	return
}

func (w *Worker) Start(ctx context.Context) {
	defer Logger.Debugf("Worker %d stopped", w.id)
	Logger.Debugf("Workder %d started", w.id)

	for {
		select {
		case job := <-w.pipe:
			job.Process()
		case <-ctx.Done():
			return
		}
	}
}
