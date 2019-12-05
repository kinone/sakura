package mlog

import (
	"log"
	"os"
	"sync"
)

type Args []interface{}

type Record struct {
	level  int
	format string
	args   Args
}

type Handler interface {
	Log(r *Record)
	Reload() error
	Close()
}

type NullHandler struct{}

func (h *NullHandler) Log(*Record)         {}
func (h *NullHandler) Reload() (err error) { return }
func (h *NullHandler) Close()              {}

type Filter func(*Record) bool

func LevelFilter(level int) Filter {
	return func(r *Record) bool {
		return r.level >= level
	}
}

type FileHandler struct {
	filename string
	filter   []Filter
	driver   *log.Logger
	fh       *os.File
}

func NewFileHandler(file string) (h *FileHandler) {
	h = &FileHandler{
		filename: file,
	}

	if err := h.init(); nil != err {
		panic(err)
	}

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

	v := append([]interface{}{Prefix(r.level)}, r.args...)
	if len(r.format) > 0 {
		format := "%s " + r.format
		h.driver.Printf(format, v...)
	} else {
		h.driver.Println(v...)
	}
}

func (h *FileHandler) Reload() (err error) {
	if nil == h.fh {
		return
	}

	if err = h.fh.Close(); nil != err {
		return
	}

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
		h.driver = log.New(h.fh, "", log.LstdFlags)
	} else {
		h.driver = log.New(os.Stdout, "", log.LstdFlags)
	}

	return
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
	defer func() {
		h.Handler.Log(&Record{level: Debug, args: Args{"log consumer stoped"}})
		h.wg.Done()
	}()

	h.Handler.Log(&Record{level: Debug, args: Args{"log consumer started"}})
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
