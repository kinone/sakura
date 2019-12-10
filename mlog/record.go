package mlog

import (
	"encoding/json"
	"strings"
)

type Args []interface{}

type Record struct {
	level  Level
	format string
	args   Args
}

func NewRecord(l Level, f string, args Args) *Record {
	return &Record{l, f, args}
}

func (r *Record) Args() Args {
	return r.args
}

func (r *Record) Level() string {
	l := r.level.String()

	return strings.ToLower(l[1 : len(l)-1])
}

func (r *Record) Format() string {
	return r.format
}

func (r *Record) SetFormat(f string) {
	r.format = f
}

func (r *Record) SetArgs(args Args) {
	r.args = args
}

func (r *Record) SetLevel(l string) {
	r.level = NewLevel(l)
}

func (r *Record) Json() *Record {
	l := len(r.args)
	if l == 0 {
		return r
	}

	var (
		data []byte
		err  error
	)

	if l > 1 {
		data, err = json.Marshal(r.args)
	} else {
		data, err = json.Marshal(r.args[0])
	}

	if nil != err {
		data = []byte("<Json format error>")
	}

	r.format = "%s"
	r.args = Args{string(data)}

	return r
}
