package orm

import (
	"fmt"
	"testing"
)

type Row struct {
}

func (r *Row) Scan(ptrs ...interface{}) error {
	for _, v := range ptrs {
		switch v.(type) {
		case *uint32:
			*v.(*uint32) = 1001
		case *string:
			*v.(*string) = "zhenhao"
		case *uint8:
			*v.(*uint8) = 25

		}
	}

	return nil
}

type Rows struct {
	*Row
}

func (r *Rows) Columns() (cols []string, err error) {
	cols = []string{"id", "name", "age"}

	return
}

type Foo struct {
	*Entity
	ID   uint32 `db:"id"`
	Name string `db:"name"`
	Age  uint8  `db:"age"`
}

func NewFoo() (f *Foo) {
	f = &Foo{
		Entity: NewEntity(),
	}

	f.SetCurrent(f)

	return
}

func TestDao_AllFields(t *testing.T) {
	fmt.Println(NewFoo().AllFields())
}

func TestDao_Load(t *testing.T) {
	f := NewFoo()

	_ = f.Load(&Rows{})

	fmt.Println(f)
}

func TestDao_LoadRow(t *testing.T) {
	f := NewFoo()

	_ = f.LoadRow(&Row{}, []string{"name", "age"})

	fmt.Println(f)
}
