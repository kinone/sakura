package workshop

import (
	"errors"
	"reflect"
)

type Job interface {
	Process()
}

type SimpleJob struct {
	callable reflect.Value
	payload  []reflect.Value
}

func NewSimpleJob(callable interface{}, payload ...interface{}) (s *SimpleJob, err error) {
	s = &SimpleJob{}

	for _, v := range payload {
		s.payload = append(s.payload, reflect.ValueOf(v))
	}

	s.callable = reflect.ValueOf(callable)
	if s.callable.Kind() != reflect.Func {
		s = nil
		err = errors.New("callable must be callable")
	}

	return
}

func (s *SimpleJob) Process() {
	s.callable.Call(s.payload)
}
