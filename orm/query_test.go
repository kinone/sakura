package orm

import "fmt"

func ExampleNewQuery() {
	query := NewQuery().
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
