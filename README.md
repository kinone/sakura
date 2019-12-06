# sakura
Just for fun

* mlog

```go
package main

import "github.com/kinone/sakura/mlog"

var logger mlog.LevelLogger

func init() {
    logger = mlog.NewLogger(&mlog.Option{
        File:   "",
        Levels: []string{"info+"},
        //Type: mlog.TMultiHandler,
        //Handlers: []*mlog.HandlerOption{
        //    {
        //        File:   "",
        //        Levels: []string{"debug"},
        //    },
        //
        //    {
        //        File:   "",
        //        Levels: []string{"info"},
        //    },
        //},
    })
}

func main() {
    defer logger.Close()

    for i := 0; i < 10; i++ {
        logger.Debug("Hello world")
        logger.Infof("Hello %d", i)
    }
}
```

* console

```go
package main

import (
    "github.com/kinone/sakura/console"
    "github.com/kinone/sakura/mlog"
    "log"
    "syscall"
    "time"
)

type Foo struct {
    console.CommandTraits
    subject string
    logger  mlog.LevelLogger
}

func (f *Foo) Configure() {
    f.AddArgument(&console.Argument{
        Name:  "s",
        Arg:   &f.subject,
        Value: "info+",
        Usage: "just for test",
    })
}

func (f *Foo) Name() string {
    return "foo:bar"
}

func (f *Foo) Execute() error {
    f.logger = mlog.NewLogger(&mlog.Option{
        Levels: nil,
    })

    f.App().HandleSignal(func() {
        f.logger.Debug("================ reload ==================")
        _ = f.logger.Reload()
    }, syscall.SIGUSR1)

    f.App().HandleShutdown(func() {
        f.logger.Debug("cmd finished")
        f.logger.Close()
    }).GoLoop(func() {
        f.logger.Alertf("Hello xxx %s", f.subject)
        time.Sleep(time.Second)
    }).GoLoop(func() {
        f.logger.Info("hahahahahah")
        time.Sleep(time.Second)
    }).Wait()

    return nil
}

func main() {
    a := console.NewApp()
    a.AddCommand(&Foo{})
    if err := a.Run(); nil != err {
        log.Fatal(err)
    }
}
```