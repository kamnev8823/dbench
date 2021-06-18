package db

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

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
	fmt.Println()
}

//PrintColumns print columns with info (type, default, is null, foreign keys)
func PrintColumns(table Table) {
	fmt.Println()
	for _, t := range table.Columns {
		var foreign ForeignKey
		// todo take out a separate function for searching foreignKeys
		for _, f := range table.ForeignKeys {
			if f.ColumnName == t.Name {
				foreign = f
			}
		}

		fmt.Printf("-- %v\n\t|-->type: %v\n\t|-->nullable: %v\n\t|-->default: %v\n", t.Name, t.Type, t.Nullable, t.Default)
		if foreign != (ForeignKey{}) {
			fmt.Printf("\t\t|-->constrain_name: %v\n\t\t|-->reference_table: %v\n\t\t|-->reference_column: %v\n", foreign.ConstraintName, foreign.ReferenceTable, foreign.ReferenceColumn)
		}
		fmt.Println()
	}
}
