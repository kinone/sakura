package workshop

import "context"

type Worker struct {
	pipe <-chan Job
	ctx  context.Context
}

func NewWorker(ctx context.Context, pipe <-chan Job) (w *Worker) {
	w = &Worker{
		pipe: pipe,
		ctx:  ctx,
	}

	return
}

func (w *Worker) Start() {
	for {
		select {
		case job := <-w.pipe:
			job.Process()
		case <-w.ctx.Done():
			return
		}
	}
}
