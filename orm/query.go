package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

type QueryType uint8
type LogicType uint8

const (
	SELECT QueryType = iota
	DELETE
	UPDATE
	INSERT
)

const (
	_ LogicType = iota
	LogicAND
	LogicOR
)

func (l LogicType) String() string {
	switch l {
	case LogicAND:
		return "AND"
	case LogicOR:
		return "OR"
	default:
		return ""
	}
}

type QueryInterface interface {
	Query() (*sql.Rows, error)
	Execute() (sql.Result, error)

	Insert(string) QueryInterface
	Value(string, string) QueryInterface
	Values(...string) QueryInterface
	Update(string) QueryInterface
	Set(...string) QueryInterface
	Select(...string) QueryInterface
	From(string) QueryInterface
	Delete(string) QueryInterface
	Where(...interface{}) QueryInterface
	AndWhere(...interface{}) QueryInterface
	OrWhere(...interface{}) QueryInterface
	OrderBy(string) QueryInterface
	Offset(int) QueryInterface
	MaxResults(int) QueryInterface
	Bind(...interface{}) QueryInterface

	String() string
}

type Query struct {
	db         *sql.DB
	qtype      QueryType
	cols       []string
	tableName  string
	set        []string
	where      []string
	orderBy    []string
	values     []string
	binds      []interface{}
	maxResults int
	offset     int
}

func NewQuery(db *sql.DB) *Query {
	return &Query{
		db: db,
	}
}

func (q *Query) Query() (*sql.Rows, error) {
	return q.db.Query(q.getSQL(), q.binds...)
}

func (q *Query) Execute() (sql.Result, error) {
	return q.db.Exec(q.getSQL(), q.binds...)
}

func (q *Query) Delete(table string) QueryInterface {
	q.qtype = DELETE
	q.tableName = table

	return q
}

func (q *Query) Insert(table string) QueryInterface {
	q.qtype = INSERT
	q.tableName = table

	return q
}

func (q *Query) Values(values ...string) QueryInterface {
	l := len(values)
	if l%2 == 1 {
		panic("query.Set: odd argument count")
	}

	q.values = append(q.values, values...)

	return q
}

func (q *Query) Value(name, value string) QueryInterface {
	q.values = append(q.values, name, value)

	return q
}

func (q *Query) Update(table string) QueryInterface {
	q.qtype = UPDATE
	q.tableName = table

	return q
}

func (q *Query) Set(args ...string) QueryInterface {
	l := len(args)
	if l%2 == 1 {
		panic("query.Set: odd argument count")
	}

	for i := 0; i < l-1; i += 2 {
		q.set = append(q.set, args[i]+"="+args[i+1])
	}

	return q
}

func (q *Query) Select(cols ...string) QueryInterface {
	q.qtype = SELECT
	q.cols = append(q.cols, cols...)

	return q
}

func (q *Query) From(table string) QueryInterface {
	q.tableName = table

	return q
}

func (q *Query) Where(parts ...interface{}) QueryInterface {
	q.where = append(q.where, NewWhere(parts, LogicType(0)).String())

	return q
}

func (q *Query) AndWhere(parts ...interface{}) QueryInterface {
	q.where = append(q.where, NewWhere(parts, LogicAND).String())

	return q
}
func (q *Query) OrWhere(parts ...interface{}) QueryInterface {
	q.where = append(q.where, NewWhere(parts, LogicOR).String())

	return q
}

func (q *Query) OrderBy(order string) QueryInterface {
	q.orderBy = append(q.orderBy, order)

	return q
}

func (q *Query) Offset(offset int) QueryInterface {
	q.offset = offset

	return q
}

func (q *Query) MaxResults(maxResults int) QueryInterface {
	q.maxResults = maxResults

	return q
}

func (q *Query) Bind(bind ...interface{}) QueryInterface {
	q.binds = append(q.binds, bind...)

	return q
}

func (q *Query) String() string {
	return q.getSQL()
}

func (q *Query) getSQL() (sql string) {
	switch q.qtype {
	case SELECT:
		sql = q.getSQLForSelect()
	case DELETE:
		sql = q.getSQLForDelete()
	case UPDATE:
		sql = q.getSQLForUpdate()
	case INSERT:
		sql = q.getSQLForInsert()
	}

	return
}

func (q *Query) getSQLForSelect() (sql string) {
	parts := []string{
		"SELECT",
		strings.Join(q.cols, ","),
		"FROM",
		q.tableName,
	}

	if len(q.where) > 0 {
		parts = append(parts, "WHERE "+strings.Join(q.where, " "))
	}

	if len(q.orderBy) > 0 {
		parts = append(parts, "ORDER BY "+strings.Join(q.orderBy, ","))
	}

	sql = strings.Join(parts, " ")

	if q.maxResults > 0 {
		sql += fmt.Sprintf(" LIMIT %d", q.maxResults)
	}

	if q.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", q.offset)
	}

	return
}

func (q *Query) getSQLForDelete() (sql string) {
	parts := []string{
		"DELETE FROM",
		q.tableName,
	}

	if len(q.where) > 0 {
		parts = append(parts, "WHERE "+strings.Join(q.where, " "))
	}

	sql = strings.Join(parts, " ")

	return
}

func (q *Query) getSQLForUpdate() (sql string) {
	parts := []string{
		"UPDATE",
		q.tableName,
		"SET",
		strings.Join(q.set, ", "),
	}

	if len(q.where) > 0 {
		parts = append(parts, "WHERE "+strings.Join(q.where, " "))
	}

	sql = strings.Join(parts, " ")

	return
}

func (q *Query) getSQLForInsert() (sql string) {
	l := len(q.values)
	if l == 0 {
		panic("values should not be empty")
	}

	cols := make([]string, 0)
	values := make([]string, 0)
	for i := 0; i < l-1; i += 2 {
		cols = append(cols, q.values[i])
		values = append(values, q.values[i+1])
	}

	parts := []string{
		"INSERT INTO",
		q.tableName,
		"(" + strings.Join(cols, ", ") + ")",
		"VALUES",
		"(" + strings.Join(values, ", ") + ")",
	}

	sql = strings.Join(parts, " ")

	return
}

type Where struct {
	parts      []string
	logic      LogicType
	innerLogic LogicType
}

func NewWhere(parts []interface{}, logic LogicType) *Where {
	innerLogic := LogicAND

	if l := len(parts); l > 1 {
		last := parts[l-1]
		switch last {
		case LogicAND:
			parts = parts[:l-1]
		case LogicOR:
			innerLogic = LogicOR
			parts = parts[:l-1]
		}
	}

	var newParts []string
	for _, v := range parts {
		vv := reflect.ValueOf(v)
		switch vv.Kind() {
		case reflect.String:
			newParts = append(newParts, v.(string))
		case reflect.Slice:
			var p []interface{}
			for i := 0; i < vv.Len(); i++ {
				p = append(p, vv.Index(i).Interface())
			}
			newParts = append(newParts, NewWhere(p, LogicType(0)).String())
		}
	}

	return &Where{
		parts:      newParts,
		logic:      logic,
		innerLogic: innerLogic,
	}
}

func (w *Where) String() (sql string) {
	sql = ""
	if w.logic > 0 {
		sql += w.logic.String() + " "
	}

	sql += "(" + strings.Join(w.parts, " "+w.innerLogic.String()+" ") + ")"

	return
}
