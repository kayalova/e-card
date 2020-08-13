package helpers

import (
	"fmt"
	"net/url"
	"strings"
	"github.com/kayalova/e-card-catalog/postgres"
)

// PrepareDBfilters ...
func PrepareDBfilters(filters url.Values) map[string]interface{} {
	mapFilters := make(map[string]interface{})
	for k, v := range filters {
		mapFilters[k] = v[0]
	}

	return mapFilters
}

// BuildSQLStatement returns sql statement string
func BuildSQLStatement(m map[string]interface{}, tableName string) string {
	sqlStatement := fmt.Sprintf(`SELECT * FROM %v WHERE `, tableName)
	for k, v := range m {
		sqlStatement += fmt.Sprintf(`%v=%v AND `, k, v)
	}

	return strings.ReplaceAll(sqlStatement[:len(sqlStatement)-4], `"`, `'`)
}

func GetDBRecords(sqlStatement string, model interface{}, typeModel struct) ([]interface{}, error) {
	db:= postgres.CreateConnection()
	defer db.Close()

	var records []typeModel
	rows, err := db.Query(sqlStatement)

	if err != nil {
		return records, err
	}

	var record typeModel
	for rows.Next() {
		err = rows.Scan(model)
		if err != nil {
			return records, nil
		}
		records = append(records, card)
	}

	return records, nil
}
