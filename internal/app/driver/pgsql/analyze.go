package pgsql

import (
	"github.com/dbench/internal/app/db"
)

func (d *Data) Analyze() []db.Table {
	handle := d.Handle
	//todo get user to change schema
	rows, _ := handle.Query("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = 'public' ORDER BY TABLE_NAME")

	var (
		tables      []db.Table
		tableName   string
		columns     []db.Column
		foreignKeys []db.ForeignKey
	)

	for rows.Next() {
		rows.Scan(&tableName)
		tables = append(tables, db.Table{Name: tableName})
	}

	d.Tables = tables

	for k, t := range d.Tables {
		columns = d.getColumns(t.Name)
		d.Tables[k].Columns = columns
	}

	for k, t := range d.Tables {
		foreignKeys = d.getForeignKeys(t.Name)
		d.Tables[k].ForeignKeys = foreignKeys
	}

	return tables
}

func (d *Data) getForeignKeys(table string) []db.ForeignKey {
	handle := d.Handle
	query := `
	SELECT  tc.constraint_name,
		tc.table_name,
		kcu.column_name,
		ccu.table_name AS references_table,
		ccu.column_name AS references_field
	FROM information_schema.table_constraints tc

	LEFT JOIN information_schema.key_column_usage kcu
		ON tc.constraint_catalog = kcu.constraint_catalog
		AND tc.constraint_schema = kcu.constraint_schema
		AND tc.constraint_name = kcu.constraint_name

	LEFT JOIN information_schema.referential_constraints rc
		ON tc.constraint_catalog = rc.constraint_catalog
		AND tc.constraint_schema = rc.constraint_schema
		AND tc.constraint_name = rc.constraint_name

	LEFT JOIN information_schema.constraint_column_usage ccu
		ON rc.unique_constraint_catalog = ccu.constraint_catalog
		AND rc.unique_constraint_schema = ccu.constraint_schema
		AND rc.unique_constraint_name = ccu.constraint_name

	WHERE tc.constraint_type = 'FOREIGN KEY' AND tc.table_name = $1;`

	rows, _ := handle.Query(query, table)

	var (
		keys []db.ForeignKey
		key  db.ForeignKey
	)

	for rows.Next() {
		rows.Scan(&key.ConstraintName, &key.TableName, &key.ColumnName, &key.ReferenceTable, &key.ReferenceColumn)
		keys = append(keys, key)
	}

	return keys
}

//getColumns get existing columns in the table
func (d *Data) getColumns(table string) []db.Column {
	handle := d.Handle
	rows, _ := handle.Query("SELECT column_name, is_nullable, data_type, column_default FROM information_schema.columns WHERE table_schema = 'public' AND table_name = $1", table)

	var (
		columns    []db.Column
		name       string
		dataType   string
		defaultVal string
		nullable   string
		isNull     bool
	)

	for rows.Next() {
		rows.Scan(&name, &nullable, &dataType, &defaultVal)

		if nullable == "YES" {
			isNull = true
		} else {
			isNull = false
		}

		columns = append(columns, db.Column{
			Name:     name,
			Type:     dataType,
			Nullable: isNull,
			Default:  defaultVal,
		})
	}

	return columns
}
