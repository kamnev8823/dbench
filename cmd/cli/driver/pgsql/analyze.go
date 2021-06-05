package pgsql

func (d *Data) Tables() []string {
	handle := d.Handle
	rows, _ := handle.Query("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'public' ORDER BY TABLE_NAME")

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
