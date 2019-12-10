package mlog

import "encoding/json"

type Args []interface{}

type Record struct {
	level  Level
	format string
	args   Args
}

func (r *Record) Args() Args {
	return r.args
}

func (r *Record) SetFormat(f string) {
	r.format = f
}

func (r *Record) SetArgs(args Args) {
	r.args = args
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
