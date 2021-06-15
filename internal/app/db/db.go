package db

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
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

// ForeignKey information about foreign keys in table
type ForeignKey struct {
	TableName       string
	ColumnName      string
	ConstraintName  string
	ReferenceTable  string
	ReferenceColumn string
}

//Column information about column in table
type Column struct {
	Name     string
	Type     string
	Nullable bool
	Default  string
}

//PrintInfoConnect print info about connection
func (d *DataStruct) PrintInfoConnect() {
	fmt.Printf("\n\tConnection info: \n\t\tDBMS: %v,\n\t\tHost: %v,\n\t\tDatabase: %v,\n\t\tUser: %v,\n\t\tPassword: %v\n",
		d.Driver,
		d.Host,
		d.Db,
		d.User,
		d.Password,
	)
}

//PrintTables print info about existing tables
func PrintTables(tables []Table) {
	if len(tables) == 0 {
		fmt.Println("Database is empty.")
		return
	}

	var lenNames []int

	for _, t := range tables {
		lenNames = append(lenNames, len(t.Name))
	}

	sort.Ints(lenNames)
	maxLenName := lenNames[len(lenNames)-1] + 2
	line := fmt.Sprint("+" + strings.Repeat("-", maxLenName) + "+")

	fmt.Println()
	fmt.Println(line)
	for _, t := range tables {
		quantitySpaces := maxLenName - len(t.Name) - 1
		fmt.Printf("| %v%"+strconv.Itoa(quantitySpaces)+"s|\n", t.Name, " ")
	}
	fmt.Println(line)
}

//todo change the printing
func PrintColumns(columns []Column) {
	fmt.Println()
	for _, t := range columns {
		fmt.Printf("--%v--\n", t.Name)
	}
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
