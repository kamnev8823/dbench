package mysql

import (
	"github.com/dbench/internal/app/db"
	"github.com/dbench/internal/app/helpers"
)

func (d *Data) Analyze() []db.Table {
	var (
		columns     []db.Column
		foreignKeys []db.ForeignKey
	)

	d.Tables = d.getTables()

	for k, t := range d.Tables {
		columns = d.getColumns(t.Name)
		d.Tables[k].Columns = columns
	}

	for k, t := range d.Tables {
		foreignKeys = d.getForeignKeys(t.Name)
		d.Tables[k].ForeignKeys = foreignKeys
	}

	return d.Tables
}

func (d *Data) getTables() []db.Table {
	handle := d.Handle
	rows, _ := handle.Query("SHOW TABLES")
	defer rows.Close()

	var tables []db.Table

	for rows.Next() {
		var table db.Table
		rows.Scan(&table.Name)
		tables = append(tables, table)
	}
	return tables
}

func (d *Data) getColumns(table string) []db.Column {
	handle := d.Handle
	query := `SELECT COLUMN_NAME, 
       				 DATA_TYPE, 
                     IS_NULLABLE, 
                     COLUMN_DEFAULT 
			  FROM information_schema.COLUMNS 
			  WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`

	rows, _ := handle.Query(query, d.Db, table)
	defer rows.Close()

	var columns []db.Column

	for rows.Next() {
		var (
			column   db.Column
			nullable string
		)
		rows.Scan(&column.Name, &column.Type, &nullable, &column.Default)

		column.Nullable = helpers.ConvertDBMSNullToBool(nullable)
		columns = append(columns, column)
	}
	return columns
}

func (d *Data) getForeignKeys(table string) []db.ForeignKey {
	handle := d.Handle
	query := `SELECT COLUMN_NAME,
					 CONSTRAINT_NAME, 
					 REFERENCED_TABLE_NAME,
                     REFERENCED_COLUMN_NAME 
			  FROM information_schema.KEY_COLUMN_USAGE
			  WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND REFERENCED_COLUMN_NAME IS NOT NULL`
	rows, _ := handle.Query(query, d.Db, table)
	defer rows.Close()

	var foreignKeys []db.ForeignKey

	for rows.Next() {
		var foreignKey db.ForeignKey
		rows.Scan(&foreignKey.ColumnName, &foreignKey.ConstraintName, &foreignKey.ReferenceTable, &foreignKey.ReferenceColumn)
		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys
}
