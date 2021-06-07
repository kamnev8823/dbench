package driver

import (
	"database/sql"
	"errors"
	"github.com/dbench/internal/app/dbench/db"
	"github.com/dbench/internal/app/dbench/driver/mysql"
	"github.com/dbench/internal/app/dbench/driver/pgsql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Supported Drivers
const (
	Mysql    = "mysql"
	Postgres = "postgres"
)

func Get(dr string) (db.Connect, error) {
	switch dr {
	case Mysql:
		return mysql.New(), nil
	case Postgres:
		return pgsql.New(), nil
	default:
		return nil, errors.New("Driver not found ")
	}
}

func Connect(connect db.Connect) (*sql.DB, error) {
	data := connect.GetDataConnect()
	database, err := sql.Open(data.Driver, connect.FormConnect())

	if err != nil {
		return nil, err
	}

	return database, nil
}
