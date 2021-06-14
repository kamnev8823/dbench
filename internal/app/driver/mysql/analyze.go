package mysql

import "github.com/dbench/internal/app/db"

func (d *Data) Analyze() []db.Table {
	return []db.Table{}
}

func (d *Data) GetTables() []string {
	handle := d.Handle
	rows, _ := handle.Query("SHOW TABLES")

	var (
		tables []string
		table  string
	)

	for rows.Next() {
		rows.Scan(&table)
		tables = append(tables, table)
	}
	return tables
}
