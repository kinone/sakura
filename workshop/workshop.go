package workshop

import (
	"context"
	"github.com/kinone/sakura/mlog"
	"sync"
)

var Logger mlog.LevelLogger = mlog.NewNullLogger()

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
		go func(id int) {
			defer w.wg.Done()
			NewWorker(id, w.pipe).Start(ctx)
		}(i)
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
