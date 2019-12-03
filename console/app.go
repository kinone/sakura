package console

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

type Application struct {
	cmds map[string]CommandInterface
}

func NewApp() *Application {
	cs := make(map[string]CommandInterface)

	return &Application{
		cmds: cs,
	}
}

func (a *Application) AddCommand(c CommandInterface) {
	a.cmds[c.Name()] = c
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
