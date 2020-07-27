package workshop

import (
	"context"
	"github.com/kinone/sakura/mlog"
	"sync"
)

var Logger mlog.LevelLogger = mlog.NewNullLogger()

type Workshop struct {
	pipe   chan Job
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func Open(hc int) (w *Workshop) {
	ctx, cancel := context.WithCancel(context.Background())

	w = &Workshop{
		pipe:   make(chan Job),
		ctx:    ctx,
		cancel: cancel,
	}

	w.hire(hc)

	return
}

func (w *Workshop) Do(job Job) {
	w.pipe <- job
}

func (w *Workshop) Close() {
	defer w.wg.Wait()
	w.cancel()
}

func (w *Workshop) hire(hc int) {
	for i := 0; i < hc; i++ {
		w.wg.Add(1)
		go func(id int) {
			defer w.wg.Done()
			NewWorker(id, w.pipe).Start(w.ctx)
		}(i)
	}
}
