package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/dbench/internal/app/db"
	"sync"
)

type Data db.DataStruct

var (
	once sync.Once

	instance Data
)

func New() *Data {
	once.Do(func() {
		instance = Data{}
	})

	return &instance
}

func (d *Data) FormConnect() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", d.User, d.Password, d.Host, d.Db)
}

func (d *Data) SetHandle(database *sql.DB) {
	d.Handle = database
}

func (d *Data) SetDataConnect(driver string, host string, user string, password string, db string) {
	d.Db = db
	d.Host = host
	d.User = user
	d.Driver = driver
	d.Password = password
}

// todo change data getter
func (d *Data) GetDataConnect() db.DataStruct {
	info := New()

	return db.DataStruct{
		Db:       info.Db,
		Host:     info.Host,
		User:     info.User,
		Driver:   info.Driver,
		Password: info.Password,
	}
}
