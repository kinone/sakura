package console

import (
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type Wrapper struct {
	sigC       chan os.Signal
	sigHanlder map[os.Signal][]func()
	wg         *sync.WaitGroup
	close      int32
}

func NewWrapper() (w *Wrapper) {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	w = &Wrapper{
		sigC:       ch,
		wg:         new(sync.WaitGroup),
		sigHanlder: make(map[os.Signal][]func()),
	}

	return
}

func (w *Wrapper) Go(f func(args ...interface{}), args ...interface{}) *Wrapper {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		f(args...)
	}()

	return w
}

func (w *Wrapper) GoLoop(f func()) *Wrapper {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			if atomic.LoadInt32(&w.close) == 1 {
				return
			}

			f()
		}
	}()

	return w
}

func (w *Wrapper) HandleSignal(f func(), s ...os.Signal) *Wrapper {
	for _, v := range s {
		w.sigHanlder[v] = append(w.sigHanlder[v], f)
	}

	return w
}

func (w *Wrapper) HandleShutdown(f func()) *Wrapper {
	return w.HandleSignal(
		f,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
	)
}

func (w *Wrapper) Wait() {
	w.HandleShutdown(w.shutdown)

	for s := range w.sigC {
		for _, f := range w.sigHanlder[s] {
			f()
		}
	}
}

func (w *Wrapper) shutdown() {
	atomic.StoreInt32(&w.close, 1)
	w.wg.Wait()
	close(w.sigC)
}
