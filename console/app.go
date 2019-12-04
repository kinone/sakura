package console

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Application struct {
	close      int32
	wg         sync.WaitGroup
	cmds       map[string]CommandInterface
	sigC       chan os.Signal
	sigHanlder map[os.Signal]func()
}

func NewApp() *Application {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGKILL,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	return &Application{
		sigC:       ch,
		cmds:       make(map[string]CommandInterface),
		sigHanlder: make(map[os.Signal]func()),
	}
}

func (a *Application) SigCh() <-chan os.Signal {
	return a.sigC
}

func (a *Application) HandleSignal(f func(), s ...os.Signal) {
	for _, v := range s {
		a.sigHanlder[v] = f
	}
}

func (a *Application) HandleShutdown(f func()) {
	a.HandleSignal(
		f,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGKILL,
	)
}

func (a *Application) SignalSinff() {
	for s := range a.sigC {
		if f, e := a.sigHanlder[s]; e {
			f()
		}
	}
}

func (a *Application) AddCommand(c CommandInterface) {
	c.SetApp(a)
	a.cmds[c.Name()] = c
}

func (a *Application) GoLoop(f func()) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		for {
			if atomic.LoadInt32(&a.close) == 1 {
				return
			}

			f()
		}
	}()
}

func (a *Application) Shutdown() {
	atomic.StoreInt32(&a.close, 1)
	a.wg.Wait()
	close(a.sigC)
}

func (a *Application) Run() (err error) {
	if len(os.Args) < 2 {
		err = errors.New("no command name input")
		return
	}

	name := os.Args[1]
	var (
		c CommandInterface
		e bool
	)

	if c, e = a.cmds[name]; !e {
		err = errors.New(fmt.Sprintf("no command %s found", name))
	}

	c.Configure()

	if err = a.parseArgs(c); nil != err {
		return
	}

	if err = flag.CommandLine.Parse(os.Args[2:]); nil != err {
		return
	}

	return c.Execute()
}

func (a *Application) parseArgs(c CommandInterface) (err error) {
	for _, v := range c.Args() {
		switch v.Arg.(type) {
		case *string:
			flag.StringVar(v.Arg.(*string), v.Name, v.Value.(string), v.Usage)
		case *int:
			flag.IntVar(v.Arg.(*int), v.Name, v.Value.(int), v.Usage)
		case *uint:
			flag.UintVar(v.Arg.(*uint), v.Name, v.Value.(uint), v.Usage)
		case *int64:
			flag.Int64Var(v.Arg.(*int64), v.Name, v.Value.(int64), v.Usage)
		case *uint64:
			flag.Uint64Var(v.Arg.(*uint64), v.Name, v.Value.(uint64), v.Usage)
		case *bool:
			flag.BoolVar(v.Arg.(*bool), v.Name, v.Value.(bool), v.Usage)
		case *time.Duration:
			flag.DurationVar(v.Arg.(*time.Duration), v.Name, v.Value.(time.Duration), v.Usage)
		default:
			err = errors.New("cannot processd arg type of " + v.Name)
		}
	}

	return
}
