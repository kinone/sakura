package workshop

import (
	"fmt"
	"github.com/kinone/sakura/mlog"
	"math/rand"
	"testing"
	"time"
)

type Point struct {
	X int
	Y int
}

func (p *Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func TestWorkshop_Do(t *testing.T) {
	count := 20

	Logger = mlog.NewLogger(&mlog.Option{
		Levels: []string{"info+"},
	})
	defer Logger.Close()

	res := make([]*Point, count)
	ws := Open(10)
	defer func() {
		ws.Close()
		for _, v := range res {
			Logger.Infof("%s\n", v)
		}
	}()

	processor := func(f *Point) {
		r := rand.Intn(3)
		f.Y = f.X + r
		time.Sleep(time.Second * time.Duration(r))
		fmt.Println("processed ", f.X)
	}

	for i := 0; i < count; i++ {
		res[i] = &Point{X: i}
		job, _ := NewSimpleJob(processor, res[i])
		ws.Do(job)
	}
}
