package qb

import (
	"time"
)

func InsertRow(db DbExecer, row DBTable) error {
	if r, ok := row.(Validatable); ok {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	if r, ok := row.(RowTimestamper); ok {
		r.RowCreating()
	}
	r, err := db.NamedExec(NewQB(row.TableName()).InsertSQL(row.Fields()...), row)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	(*row.PK()) = id
	return nil
}

func DeleteRow(db DbExecer, row DBTable) error {
	q := NewQB(row.TableName()).Where("id = ?", *row.PK())
	_, err := db.Exec(q.DeleteSQL(), q.Args...)
	return err
}

func UpdateRow(db DbExecer, row DBTable) error {
	if r, ok := row.(Validatable); ok {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	if r, ok := row.(RowTimestamper); ok {
		r.RowUpdating()
	}
	_, err := db.NamedExec(NewQB(row.TableName()).UpdateSQL([]string{"id"}, row.Fields()...), row)
	return err
}

func ListRows(db DbSelector, result interface{}, q *QB) error {
	err := db.Select(result, q.SelectSQL(), q.Args...)
	return err
}

func GetRowByPK(db DbGetter, result DBTable, pk int64) error {
	q := NewQB(result.TableName()).Where("id = ?", pk)
	err := db.Get(result, q.SelectSQL(), q.Args...)
	return err
}

type RowTimestamper interface {
	RowCreating()
	RowUpdating()
}
type RowTimestamps struct {
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (r *RowTimestamps) RowCreating() {
	r.CreatedAt, r.UpdatedAt = time.Now(), time.Now()
}
func (r *RowTimestamps) RowUpdating() {
	r.UpdatedAt = time.Now()
}
