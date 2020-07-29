package orm

import (
	"database/sql"
	"reflect"
	"unsafe"
)

type DaoInterface interface {
	AllFields() map[string]string
	FieldPtr([]string) []interface{}
}

type RowInterface interface {
	Scan(...interface{}) error
}

type RowsInterface interface {
	RowInterface
	Columns() ([]string, error)
}

type Dao struct {
	current DaoInterface
}

func NewDao() (d *Dao) {
	d = &Dao{}

	return
}

func (d *Dao) SetCurrent(c DaoInterface) {
	d.current = c
}

func (d *Dao) Load(r RowsInterface) (err error) {
	var cols []string

	if cols, err = r.Columns(); nil != err {
		return
	}

	err = r.Scan(d.FieldPtr(cols)...)

	return
}

func (d *Dao) LoadRow(r RowInterface, cols []string) error {
	return r.Scan(d.FieldPtr(cols)...)
}

func (d *Dao) AllFields() (m map[string]string) {
	rtt := reflect.TypeOf(d.current).Elem()
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

func (d *Dao) FieldPtr(cols []string) (r []interface{}) {
	m := d.current.AllFields()
	rvt := reflect.ValueOf(d.current).Elem()
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
