package orm

import (
	"fmt"
)

func ExampleQuery_Select() {
	query := NewQuery(nil).
		Select("id", "name").
		From("users").
		Where("name LIKE ?", "nickname LIKE ?", LogicOR).
		AndWhere([]interface{}{"is_delete=?", []string{"flag=?", "haha=?"}, LogicOR}).
		OrWhere("id=?").
		OrderBy("id DESC").
		Offset(10).
		MaxResults(20).
		Bind("z", "z", 20, 1, 0)

	fmt.Println(query.String())
	// output:
	// SELECT id,name FROM users WHERE (name LIKE ? OR nickname LIKE ?) AND ((is_delete=? OR (flag=? AND haha=?))) OR (id=?) ORDER BY id DESC LIMIT 20 OFFSET 10
}

func ExampleQuery_Update() {
	query := NewQuery(nil).Update("foo").Set("name", "?", "age", "age+1").Where("id=1001", "flag=0")

	fmt.Println(query.String())
	// output:
	// UPDATE foo SET name=?, age=age+1 WHERE (id=1001 AND flag=0)
}

func ExampleQuery_Insert() {
	query := NewQuery(nil).Insert("foo").Values(
		"name", "'zhenhao'",
		"email", "'hit_zhenhao@163.com'",
		"created", "?",
	)
	fmt.Println(query.String())
	// output:
	// INSERT INTO foo (name, email, created) VALUES ('zhenhao', 'hit_zhenhao@163.com', ?)
}

func ExampleQuery_Delete() {
	query := NewQuery(nil).Delete("foo").Where("id=101").AndWhere("flag=0")
	fmt.Println(query.String())
	// output:
	// DELETE FROM foo WHERE (id=101) AND (flag=0)
}
