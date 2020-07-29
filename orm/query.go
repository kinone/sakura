package orm

import (
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
	Select(...string) QueryInterface
	From(string) QueryInterface
	Where(...interface{}) QueryInterface
	AndWhere(...interface{}) QueryInterface
	OrWhere(...interface{}) QueryInterface
	OrderBy(string) QueryInterface
	Offset(int) QueryInterface
	MaxResults(int) QueryInterface
	Bind(...interface{}) QueryInterface

	String() string
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

type Query struct {
	qtype      QueryType
	cols       []string
	from       string
	set        map[string]string
	where      []*Where
	orderBy    []string
	values     []string
	binds      []interface{}
	maxResults int
	offset     int
}

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) Select(cols ...string) QueryInterface {
	q.qtype = SELECT
	q.cols = append(q.cols, cols...)
	return q
}

func (q *Query) From(table string) QueryInterface {
	q.from = table
	return q
}

func (q *Query) Where(parts ...interface{}) QueryInterface {
	q.where = append(q.where, NewWhere(parts, LogicType(0)))
	return q
}

func (q *Query) AndWhere(parts ...interface{}) QueryInterface {
	q.where = append(q.where, NewWhere(parts, LogicAND))

	return q
}
func (q *Query) OrWhere(parts ...interface{}) QueryInterface {
	q.where = append(q.where, NewWhere(parts, LogicOR))

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
		sql = q.getSQLForUpdate()
	}

	return
}

func (q *Query) getSQLForSelect() (sql string) {
	var part []string
	part = append(part, "SELECT "+strings.Join(q.cols, ","))
	part = append(part, "FROM "+q.from)

	if len(q.where) > 0 {
		var where []string
		for _, v := range q.where {
			where = append(where, v.String())
		}
		part = append(part, "WHERE "+strings.Join(where, " "))
	}

	if len(q.orderBy) > 0 {
		part = append(part, "ORDER BY "+strings.Join(q.orderBy, ","))
	}

	sql = strings.Join(part, " ")

	if q.maxResults > 0 {
		sql += fmt.Sprintf(" LIMIT %d", q.maxResults)
	}

	if q.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", q.offset)
	}

	return
}

func (q *Query) getSQLForDelete() (sql string) {
	return
}

func (q *Query) getSQLForUpdate() (sql string) {
	return
}

func (q *Query) getSQLForInsert() (sql string) {
	return
}
