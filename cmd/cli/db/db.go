package db

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Connect interface {
	ParseDump
	Handler
	Analyze
	FormConnect() string
	GetDataConnect() DataStruct
	SetDataConnect(driver string, host string, user string, password string, db string)
}

//ParseDump separation of logic for reading a dump for different dbms
type ParseDump interface {
	ParseDump(path string) (*ParseStruct, error)
}

//Analyze get information about the used database
type Analyze interface {
	Tables() []string
}

//Handler handle operation
type Handler interface {
	SetHandle(database *sql.DB)
}

type ParseStruct struct {
	Data []string
}

type DataStruct struct {
	Db       string
	Host     string
	User     string
	Driver   string
	Password string
	Handle   *sql.DB
}

//PrintInfoConnect print info about connection
func (d *DataStruct) PrintInfoConnect() {
	fmt.Printf("\n\tConnection info: \n\t\tDBMS: %v,\n\t\tHost: %v,\n\t\tDatabase: %v,\n\t\tUser: %v,\n\t\tPassword: %v\n",
		d.Driver,
		d.Host,
		d.Db,
		d.User,
		d.Password)
}

//PrintTables print info about existing tables
func PrintTables(tables []string) {

	if len(tables) < 0 {
		fmt.Println("Database is empty.")
		return
	}

	var lenNames []int

	for _, v := range tables {
		lenNames = append(lenNames, len(v))
	}

	sort.Ints(lenNames)
	maxLenName := lenNames[len(lenNames)-1] + 2
	line := fmt.Sprint("+" + strings.Repeat("-", maxLenName) + "+")

	fmt.Println()
	fmt.Println(line)
	for _, v := range tables {
		quantitySpaces := maxLenName - len(v) - 1
		fmt.Printf("| "+v+"%"+strconv.Itoa(quantitySpaces)+"s|\n", " ")
	}
	fmt.Println(line)
}
