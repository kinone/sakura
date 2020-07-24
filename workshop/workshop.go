package workshop

import (
	"context"
	"sync"
)

type Workshop struct {
	pipe  chan Job
	wg    sync.WaitGroup
	close context.CancelFunc
}

func Open(hc int) (w *Workshop) {
	w = &Workshop{
		pipe: make(chan Job),
	}

	ctx := context.Background()
	ctx, w.close = context.WithCancel(ctx)

	for i := 0; i < hc; i++ {
		w.wg.Add(1)
		go func() {
			defer w.wg.Done()
			NewWorker(ctx, w.pipe).Start()
		}()
	}

	return
}

func (w *Workshop) Do(job Job) {
	w.pipe <- job
}

func (w *Workshop) Close() {
	defer w.wg.Wait()
	w.close()
}
