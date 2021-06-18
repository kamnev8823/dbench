package db

import (
	"database/sql"
	"errors"
)

type Connect interface {
	Handler
	Analyze
	FormConnect() string
	// todo change data getter
	GetDataConnect() DataStruct
	SetDataConnect(driver string, host string, user string, password string, db string)
}

//Analyze get information about the used database
type Analyze interface {
	Analyze() []Table
}

//Handler handle operation
type Handler interface {
	SetHandle(database *sql.DB)
}

//DataStruct info about connection data
type DataStruct struct {
	Db       string
	Host     string
	User     string
	Driver   string
	Password string
	Handle   *sql.DB
	Tables   []Table
}

//Table info about the table
type Table struct {
	Name        string
	Columns     []Column
	ForeignKeys []ForeignKey
}

//Column information about column in table
type Column struct {
	Name     string
	Type     string
	Nullable bool
	Default  string
}

// ForeignKey information about foreign keys in table
type ForeignKey struct {
	ColumnName      string
	ConstraintName  string
	ReferenceTable  string
	ReferenceColumn string
}

// PrimaryKey todo test
type PrimaryKey struct {
	ColumnName     string
	ConstraintName string
}

// IndexKey todo test
type IndexKey struct {
	Name string
	Def  string
}

// FindTable findTable Find table in database
func FindTable(name string, tables []Table) (Table, error) {
	for _, v := range tables {
		if v.Name == name {
			return v, nil
		}
	}
	return Table{}, errors.New("No such table exists ")
}
