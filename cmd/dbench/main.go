package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dbench/internal/app/dbench/db"
	"github.com/dbench/internal/app/dbench/driver"
	"github.com/dbench/internal/app/dbench/helpers"
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
	MonitorHelp        = "\\h"
	MonitorExit        = "\\q"
	MonitorStart       = "\\s"
	MonitorHistory     = "\\hc"
	MonitorInfoConnect = "\\i"
	MonitorLoadFile    = "\\lf"
	MonitorTables      = "\\t"
	MonitorConfig      = "\\cf"
)

var (
	dbDriver   = flag.String(DBDriver, None, "Database driver, you can set pgsql/mysql.")
	dbHost     = flag.String(DBHost, "17.0.0.1", "Database host.")
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
		MonitorLoadFile + " - [file path].sql - load sql file for creating table(or insert or drop tables).\n\t\t" +
		MonitorStart + "  - start tests.\n\t\t" +
		MonitorTables + "  - check existing tables\n\t\t" +
		MonitorExit + "  - exit.")
}

func main() {
	flag.Parse()

	if *help {
		helpText()
		os.Exit(0)
	}

	dbDriver := *dbDriver
	//dbDriver := "mysql"
	if dbDriver == None || (dbDriver != driver.Mysql && dbDriver != driver.Postgres) {
		log.Fatalf("Error, invalid database driver!")
	}

	dbName := *dbName
	//dbName := "tests"
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
	fmt.Print("Welcome to DBench monitor. Successful database connection.\nType '\\h' for help. Type '\\q' for exit.\n\n")

	terminal := Terminal{
		[]string{},
		fmt.Sprintf("dbench[%v]> ", dbDriver),
	}

	existsTables := conn.Tables()

	fmt.Printf("Your database has %v tables. \n\n", len(existsTables))

	terminal.Cursor()

	for sc.Scan() {
		terminal.Cursor()
		command := sc.Text()
		terminal.SaveHistory(command)

		if strings.Contains(command, MonitorLoadFile) {
			file := strings.Replace(command, MonitorLoadFile, "", 1)
			file = strings.TrimSpace(file)

			if file == "" || !strings.Contains(command, ".sql") {
				log.Println("Error, invalid file!")
				terminal.Cursor()
				continue
			}

			dump, err := conn.ParseDump(file)

			if err != nil {
				log.Println(err)
				terminal.Cursor()
				continue
			}

			for _, v := range dump.Data {
				_, err = database.Exec(v)
			}

			fmt.Println("Ok")
			terminal.Cursor()
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
			terminal.PrintHistory()
		case MonitorInfoConnect:
			data := conn.GetDataConnect()
			data.PrintInfoConnect()
		case MonitorTables:
			tables := conn.Tables()
			db.PrintTables(tables)
		default:
			fmt.Println("Invalid command. \\h - for get helping information.")
		}

		terminal.Cursor()
	}
}
