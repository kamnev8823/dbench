package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dbench/internal/app/db"
	"github.com/dbench/internal/app/driver"
	"github.com/dbench/internal/app/helpers"
	"github.com/dbench/internal/app/terminal"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
)

// Databases information flags
const (
	DBDriver   = "driver"
	DBHost     = "host"
	DBUser     = "user"
	DBPassword = "password"
	DBName     = "db"
	DBClear    = "clear"
	Help       = "help"
	None       = "none"
)

// Monitor commands
const (
	MonitorHelp          = "\\h"
	MonitorExit          = "\\q"
	MonitorConfig        = "\\cf"
	MonitorHistory       = "\\hc"
	MonitorInfoConnect   = "\\i"
	MonitorStart         = "\\s"
	MonitorLoadFile      = "\\lf"
	MonitorTables        = "\\t"
	MonitorTableColumns  = "\\tc"
	MonitorRelationships = "\\tr"
)

var (
	dbDriver   = flag.String(DBDriver, None, "Database driver, you can set postgres/mysql.")
	dbHost     = flag.String(DBHost, "127.0.0.1", "Database host.")
	dbUser     = flag.String(DBUser, "root", "Database user name.")
	dbPassword = flag.String(DBPassword, "root", "Database user password.")
	dbName     = flag.String(DBName, None, "Database name.")
	help       = flag.Bool(Help, false, "Get help info.")
	clear      = flag.Bool(DBClear, false, "Delete generated tables after testing.")
)

// helpText help information
func helpText() {
	fmt.Println("FLAGS:")
	flag.PrintDefaults()
}

// monitorHelpText monitor help information
func monitorHelpText() {
	fmt.Println("\n\tLists of all commands:\n\t\t" +
		MonitorHelp + "  - help information.\n\t\t" +
		MonitorConfig + " - check config.\n\t\t" +
		MonitorHistory + " - history commands.\n\t\t" +
		MonitorInfoConnect + "  - info about connection.\n\t\t" +
		MonitorLoadFile + " [file path].sql - load sql file for creating table(or insert or drop tables).\n\t\t" +
		MonitorStart + "  - start tests.\n\t\t" +
		MonitorTables + "  - check existing tables\n\t\t" +
		MonitorTableColumns + " [table] - check existing table column\n\t\t" +
		MonitorRelationships + " - check relationship tables." +
		"\n\t\t\tRelations are determined by the name of the column that matches the pattern 'table_id' or by the presence of a foreign key.\n\t\t\t" +
		"To manage relationships, use the d command.\n\t\t" +
		MonitorExit + "  - exit.")
}

func main() {
	flag.Parse()

	if *help {
		helpText()
		os.Exit(0)
	}

	dbDriver := *dbDriver
	if dbDriver == None || (dbDriver != driver.Mysql && dbDriver != driver.Postgres) {
		log.Fatalf("Error, invalid database driver!")
	}

	dbName := *dbName
	if dbName == None {
		log.Fatalf("Error, invalid database name!")
	}

	conn, err := driver.Get(dbDriver)
	if err != nil {
		log.Fatalf(err.Error())
	}

	dbHost := *dbHost
	dbUser := *dbUser
	dbPassword := *dbPassword
	conn.SetDataConnect(dbDriver, dbHost, dbUser, dbPassword, dbName)

	database, err := driver.Connect(conn)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err = database.Ping(); err != nil {
		log.Fatalf(err.Error())
	}
	conn.SetHandle(database)
	defer database.Close()

	sc := bufio.NewScanner(os.Stdin)

	t := terminal.Terminal{
		History: []string{},
		Cursor:  fmt.Sprintf("dbench[%v]> ", dbDriver),
	}
	existsTables := conn.Analyze()
	fmt.Print("Welcome to DBench monitor. Successful database connection.\nType '\\h' for help. Type '\\q' for exit.\n\n")
	fmt.Printf("Your database has %v tables. \n\n", len(existsTables))
	t.PrintCursor()

	for sc.Scan() {
		command := sc.Text()
		t.SaveHistory(command)

		if strings.Contains(command, MonitorLoadFile) {
			file := strings.Replace(command, MonitorLoadFile, "", 1)
			file = strings.TrimSpace(file)

			if file == "" || !strings.Contains(command, ".sql") {
				log.Println("Error, invalid file!")
				t.PrintCursor()
				continue
			}

			dump, err := conn.ParseDump(file)

			if err != nil {
				log.Println(err)
				t.PrintCursor()
				continue
			}

			for _, v := range dump.Data {
				_, err = database.Exec(v)
			}

			fmt.Println("Ok")
			t.PrintCursor()
			continue
		}

		if strings.Contains(command, MonitorTableColumns) {
			tableName := strings.Replace(command, MonitorTableColumns, "", 1)
			tableName = strings.TrimSpace(tableName)

			tables := conn.Analyze()
			table, err := db.FindTable(tableName, tables)

			if err != nil {
				log.Println(err)
				t.PrintCursor()
				continue
			}

			db.PrintColumns(table.Columns)
			t.PrintCursor()
			continue
		}

		switch command {
		case MonitorExit:
			fmt.Println("Bye")
			return
		case MonitorHelp:
			monitorHelpText()
		case MonitorConfig:
			helpers.PrintConfig()
		case MonitorHistory:
			t.PrintHistory()
		case MonitorInfoConnect:
			data := conn.GetDataConnect()
			data.PrintInfoConnect()
		case MonitorTables:
			tables := conn.Analyze()
			db.PrintTables(tables)
		default:
			fmt.Println("Invalid command. \\h - for get helping information.")
		}

		t.PrintCursor()
	}
}