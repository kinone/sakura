package orm

import (
	"fmt"
	"testing"
)

type Foo struct {
	*Dao
	ID      uint32 `db:"id"`
	Name    string `db:"name"`
	Age     uint8  `db:"age"`
	Created string `db:"created"`
}

func NewFoo() (f *Foo) {
	f = &Foo{
		Dao: NewDao(),
	}

	f.SetCurrent(f)

	return
}

func TestDao_AllFields(t *testing.T) {
	fmt.Println(NewFoo().AllFields())
}

func TestDao_FieldPtr(t *testing.T) {
	f := NewFoo()
	ptrs := f.FieldPtr([]string{"id", "name", "age"})

	id := ptrs[0].(*uint32)
	*id = 1001
	name := ptrs[1].(*string)
	*name = "zhenhao"
	age := ptrs[2].(*uint8)
	*age = 25

	fmt.Println(f)
}
