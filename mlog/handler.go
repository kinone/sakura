package mlog

import (
	"log"
	"os"
	"sync"
)

type Handler interface {
	Log(int, ...interface{})
	Logf(int, string, ...interface{})
	Reload() error
	Close()
}

type NullHandler struct {
}

func (h *NullHandler) Log(int, ...interface{})          {}
func (h *NullHandler) Logf(int, string, ...interface{}) {}
func (h *NullHandler) Reload() (err error)              { return }
func (h *NullHandler) Close()                           {}

type FileHandler struct {
	filename string
	filter   func(level int) bool
	driver   *log.Logger
	fh       *os.File
}

func NewFileHandler(file string) (h *FileHandler) {
	h = &FileHandler{
		filename: file,
		filter: func(int) bool {
			return true
		},
	}

	if err := h.init(); nil != err {
		panic(err)
	}

	return
}

func (h *FileHandler) Log(level int, v ...interface{}) {
	if !h.filter(level) {
		return
	}

	v = append([]interface{}{Prefix(level)}, v...)
	h.driver.Println(v...)
}

func (h *FileHandler) Logf(level int, format string, v ...interface{}) {
	if !h.filter(level) {
		return
	}

	format = "%s " + format
	v = append([]interface{}{Prefix(level)}, v...)
	h.driver.Printf(format, v...)
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
	ch chan func()
	wg *sync.WaitGroup
}

func NewSmartHandler(handler Handler) (h *SmartHandler) {
	h = &SmartHandler{
		Handler: handler,
		ch:      make(chan func()),
		wg:      new(sync.WaitGroup),
	}

	h.wg.Add(1)
	go h.consumer()

	return
}

func (h *SmartHandler) consumer() {
	defer func() {
		h.Handler.Log(Debug, "log consumer stoped")
		h.wg.Done()
	}()

	h.Handler.Log(Debug, "log consumer started")
	for f := range h.ch {
		f()
	}
}

func (h *SmartHandler) Log(level int, v ...interface{}) {
	h.ch <- func() {
		h.Handler.Log(level, v...)
	}
}

func (h *SmartHandler) Logf(level int, format string, v ...interface{}) {
	h.ch <- func() {
		h.Handler.Logf(level, format, v...)
	}
}

func (h *SmartHandler) Close() {
	defer h.Handler.Close()

	close(h.ch)
	h.wg.Wait()
}
