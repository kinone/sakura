package orm

import (
	"database/sql"
	"reflect"
	"unsafe"
)

type RowInterface interface {
	Scan(...interface{}) error
}

type RowsInterface interface {
	RowInterface
	Columns() ([]string, error)
}

type EntityInterface interface {
	AllFields() map[string]string
	FieldPtr([]string) []interface{}
	Load(r RowsInterface) error
	LoadRow(RowInterface, []string) error
	SetCurrent(EntityInterface)
}

type Entity struct {
	current EntityInterface
}

func NewEntity() (e *Entity) {
	e = &Entity{}

	return
}

func (e *Entity) SetCurrent(c EntityInterface) {
	e.current = c
}

func (e *Entity) Load(r RowsInterface) (err error) {
	var cols []string

	if cols, err = r.Columns(); nil != err {
		return
	}

	err = r.Scan(e.FieldPtr(cols)...)

	return
}

func (e *Entity) LoadRow(r RowInterface, cols []string) error {
	return r.Scan(e.FieldPtr(cols)...)
}

func (e *Entity) AllFields() (m map[string]string) {
	rtt := reflect.TypeOf(e.current).Elem()
	m = make(map[string]string)
	for i := 0; i < rtt.NumField(); i++ {
		f := rtt.Field(i)
		n := f.Tag.Get("db")

		if len(n) == 0 || n == "-" {
			continue
		}

		m[n] = f.Name
	}

	return
}

func (e *Entity) FieldPtr(cols []string) (r []interface{}) {
	m := e.current.AllFields()
	rvt := reflect.ValueOf(e.current).Elem()
	for _, col := range cols {
		var (
			name string
			e    bool
		)
		// Tag中未声明的字段
		if name, e = m[col]; !e {
			r = append(r, new(sql.RawBytes))
			continue
		}

		f := rvt.FieldByName(name)
		p := unsafe.Pointer(f.UnsafeAddr())
		v := reflect.NewAt(f.Type(), p)
		r = append(r, v.Interface())
	}

	return
}
