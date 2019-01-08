package qb

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DBTable interface {
	PK() *int64
	TableName() string
	Fields() []string
}

type DbGetter interface {
	Get(interface{}, string, ...interface{}) error
}

type DbSelector interface {
	Select(interface{}, string, ...interface{}) error
	NamedQuery(string, interface{}) (*sqlx.Rows, error)
}

type DbExecer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	NamedExec(string, interface{}) (sql.Result, error)
}

// Validatable is the interface indicating the type implementing it supports data validation.
type Validatable interface {
	// Validate validates the data and returns an error if validation fails.
	Validate() error
}
