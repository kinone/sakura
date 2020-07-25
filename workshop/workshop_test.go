package workshop

import (
	"fmt"
	"github.com/kinone/sakura/mlog"
	"math/rand"
	"testing"
	"time"
)

func TestWorkshop_Do(t *testing.T) {
	Logger = mlog.NewLogger(nil)
	defer Logger.Close()

	ws := Open(10)
	defer ws.Close()

	hello := func(i int) {
		time.Sleep(time.Second * time.Duration(rand.Intn(3)))
		fmt.Println("Hello ", i)
	}

	for i := 0; i <= 20; i++ {
		job, _ := NewSimpleJob(hello, i)
		ws.Do(job)
	}
}
