package helpers

// ConvertDBMSNullToBool convert DBMS representations of nullable columns
// the values of which are "YES" and "NO", to a boolean type
func ConvertDBMSNullToBool(s string) bool {
	if s == "YES" {
		return true
	}
	return false
}
