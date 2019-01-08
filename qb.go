package qb

import (
	"fmt"
	"strings"
)

type QB struct {
	Table       string
	fields      string
	OrderBy     string
	Start, Size int
	conditions  []string
	open        int
	op          bool
	Args        []interface{}
}

func NewQB(table string) *QB {
	return &QB{
		Table:  table,
		fields: "*",
	}
}

func (wb *QB) Open(op string) *QB {
	if wb.op {
		wb.conditions = append(wb.conditions, op)
	}
	wb.open++
	wb.op = false
	wb.conditions = append(wb.conditions, "(")
	return wb
}

func (wb *QB) Close() *QB {
	if wb.open > 0 {
		wb.open--
		wb.conditions = append(wb.conditions, ")")
		wb.op = true
	}
	return wb
}

func (wb *QB) where(op, cond string, args ...interface{}) *QB {
	if !wb.op {
		wb.conditions = append(wb.conditions, cond)
	} else {
		wb.conditions = append(wb.conditions, op, cond)
	}
	wb.op = true
	wb.Args = append(wb.Args, args...)
	return wb
}

func (wb *QB) Where(cond string, args ...interface{}) *QB {
	return wb.where("and", cond, args...)
}

func (wb *QB) Or(cond string, args ...interface{}) *QB {
	return wb.where("or", cond, args...)
}

func (wb *QB) Build() string {
	var sql string
	if len(wb.conditions) > 1 && wb.conditions[len(wb.conditions)-1] == "(" {
		wb.conditions = wb.conditions[:len(wb.conditions)-2]
		wb.open--
	}
	sql = strings.Join(wb.conditions, " ")

	if wb.open > 0 {
		sql = sql + strings.Repeat(")", wb.open)
	}
	return sql
}

func (qb *QB) Order(order string) *QB {
	qb.OrderBy = order
	return qb
}

func (qb *QB) Fields(f string) *QB {
	qb.fields = f
	return qb
}

func (qb *QB) Limit(start, size int) *QB {
	if start < 0 {
		start = 0
	}
	if size < 0 {
		size = 0
	}
	qb.Start, qb.Size = start, size
	return qb
}

func (qb *QB) SelectSQL() string {
	sql := "select " + qb.fields + " from " + qb.Table
	if len(qb.conditions) > 0 {
		sql = sql + " where " + qb.Build()
	}
	if qb.OrderBy != "" {
		sql = sql + " order by " + qb.OrderBy
	}
	if qb.Start != 0 || qb.Size != 0 {
		sql = sql + fmt.Sprintf(" limit %d,%d", qb.Start, qb.Size)
	}
	return sql
}

func (qb *QB) CountSQL() string {
	sql := "select count(0) as cnt from " + qb.Table
	if len(qb.conditions) > 0 {
		sql = sql + " where " + qb.Build()
	}

	return sql
}

func (qb *QB) InsertSQL(fields ...string) string {
	sql := "insert into " + qb.Table + " ("
	sql = sql + strings.Join(fields, ",") + ") values ("
	for _, f := range fields {
		sql = sql + ":" + f + ","
	}
	sql = sql[:len(sql)-1] + ")"
	return sql
}

func (qb *QB) UpdateSQL(whereFields []string, fields ...string) string {
	sql := "update " + qb.Table + " set "
	for _, f := range fields {
		sql = sql + f + "=:" + f + ","
	}
	sql = sql[:len(sql)-1] + " where "
	for _, w := range whereFields {
		sql = sql + w + "=:" + w + " and "
	}
	return sql[:len(sql)-4]
}

func (qb *QB) DeleteSQL() string {
	sql := "delete from " + qb.Table
	if len(qb.conditions) > 0 {
		sql = sql + " where " + qb.Build()
	}
	return sql
}

func CountRows(db DbGetter, q *QB) (count int64, err error) {
	err = db.Get(&count, q.CountSQL(), q.Args...)
	return
}
