package mysql

import (
	"bufio"
	"database/sql"
	"fmt"
	db2 "github.com/dbench/internal/app/dbench/db"
	"io"
	"os"
	"strings"
	"sync"
)

type Data db2.DataStruct

var (
	fOnce    sync.Once
	instance Data
)

func New() *Data {
	fOnce.Do(func() {
		instance = Data{}
	})

	return &instance
}

func (d *Data) SetHandle(database *sql.DB) {
	d.Handle = database
}

func (d *Data) FormConnect() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", d.User, d.Password, d.Host, d.Db)
}

func (d *Data) SetDataConnect(driver string, host string, user string, password string, db string) {
	d.Db = db
	d.Host = host
	d.User = user
	d.Driver = driver
	d.Password = password
}

func (d *Data) GetDataConnect() db2.DataStruct {
	info := New()

	return db2.DataStruct{
		Db:       info.Db,
		Host:     info.Host,
		User:     info.User,
		Driver:   info.Driver,
		Password: info.Password,
	}
}

//ParseDump Method for reading sql file for mysql dbms
func (d *Data) ParseDump(path string) (*db2.ParseStruct, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	parseStruct := &db2.ParseStruct{}

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
