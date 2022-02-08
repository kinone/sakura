package mlog

import (
	"log"
	"os"
	"sync"
)

type HandlerOption struct {
	File   string
	Levels []string
}

type Handler interface {
	Log(r *Record)
	Reload() error
	Close()
}

type FileHandler struct {
	filename  string
	filter    []Filter
	driver    *log.Logger
	flag      int
	fh        *os.File
	showLevel bool
}

func NewFileHandler(file string) (h *FileHandler) {
	h = &FileHandler{
		filename: file,
		flag:     log.LstdFlags,
	}

	if err := h.init(); nil != err {
		panic(err)
	}

	return
}

func NewBareHandler(file string) (h *FileHandler) {
	h = NewFileHandler(file)
	h.flag = 0
	h.driver.SetFlags(h.flag)

	return
}

func NewLevelHandler(file string, level ...string) (h *FileHandler) {
	h = NewFileHandler(file)
	h.AddFilter(LevelFilter(level...))
	h.showLevel = true

	return
}

func (h *FileHandler) AddFilter(f ...Filter) {
	h.filter = append(h.filter, f...)
}

func (h *FileHandler) Log(r *Record) {
	for _, f := range h.filter {
		if !f(r) {
			return
		}
	}

	var (
		v      = r.args
		format = r.format
	)

	if h.showLevel {
		v = append([]interface{}{r.level.String()}, v...)
		if len(format) > 0 {
			format = "%s " + r.format
		}
	}

	if len(format) > 0 {
		h.driver.Printf(format, v...)
	} else {
		h.driver.Println(v...)
	}
}

func (h *FileHandler) Reload() (err error) {
	if nil == h.fh {
		return
	}

	old := h.fh
	defer func() { // file changed
		if nil != old && nil == err {
			err = old.Close()
		}
	}()

	return h.init()
}

func (h *FileHandler) Close() {
	if nil != h.fh {
		_ = h.fh.Close()
	}
}

func (h *FileHandler) init() (err error) {
	if len(h.filename) > 0 {
		if h.fh, err = os.OpenFile(h.filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664); nil != err {
			return
		}
	}

	if nil != h.fh {
		h.driver = log.New(h.fh, "", h.flag)
	} else {
		h.driver = log.New(os.Stdout, "", h.flag)
	}

	return
}

type JsonHandler struct {
	Handler
}

func NewJsonHandler(handler Handler) (h *JsonHandler) {
	h = &JsonHandler{handler}

	return
}

func (h *JsonHandler) Log(r *Record) {
	h.Handler.Log(r.Json())
}

type SmartHandler struct {
	Handler
	ch chan *Record
	wg *sync.WaitGroup
}

func NewSmartHandler(handler Handler) (h *SmartHandler) {
	h = &SmartHandler{
		Handler: handler,
		ch:      make(chan *Record),
		wg:      new(sync.WaitGroup),
	}

	h.wg.Add(1)
	go h.consumer()

	return
}

func (h *SmartHandler) consumer() {
	defer h.wg.Done()
	//defer h.Handler.Log(&Record{level: Debug, args: Args{"log consumer stoped"}})
	//h.Handler.Log(&Record{level: Debug, args: Args{"log consumer started"}})
	for r := range h.ch {
		h.Handler.Log(r)
	}
}

func (h *SmartHandler) Log(r *Record) {
	h.ch <- r
}

func (h *SmartHandler) Close() {
	defer h.Handler.Close()

	close(h.ch)
	h.wg.Wait()
}
