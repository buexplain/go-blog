package s_database

import "github.com/buexplain/go-blog/dao"

type Tables map[string][]string

func GetTables() Tables {
	tables, err := dao.Dao.DBMetas()
	if err != nil {
		panic(err)
	}
	result := make(Tables)
	for _, table := range tables {
		columns := table.Columns()
		tmp := make([]string, len(columns), len(columns))
		for k, column := range columns {
			tmp[k] = column.Name
		}
		result[table.Name] = tmp
	}
	return result
}
