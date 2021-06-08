package mysql

func (d *Data) Tables() []string {
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
