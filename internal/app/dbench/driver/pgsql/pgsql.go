package pgsql

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/dbench/internal/app/dbench/db"
	"io"
	"os"
	"strings"
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

//ParseDump Method for reading sql file for postgresql dbms
func (d *Data) ParseDump(path string) (*db.ParseStruct, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	parseStruct := &db.ParseStruct{}

	for {
		l, _, err := reader.ReadLine()
		line := string(l)

		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		if line == "" {
			continue
		} else if strings.Contains(line, "/*!") || strings.Contains(line, "*/") {
			continue
		} else if strings.Contains(line, "--") {
			continue
		}

		if len(parseStruct.Data) > 0 && !strings.HasSuffix(parseStruct.Data[len(parseStruct.Data)-1], ";") {
			parseStruct.Data[len(parseStruct.Data)-1] += line
			continue
		}

		parseStruct.Data = append(parseStruct.Data, line)
	}

	return parseStruct, err
}
