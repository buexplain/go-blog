package s_services

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"strings"
	"time"
	"xorm.io/core"
	"xorm.io/xorm"
)

//检查字段值是否存在
func CheckUnique(tableName string, field string, value interface{}, selfID ...int) bool {
	d := dao.Dao.Table("`" + tableName + "`")
	d.Where("`"+field+"`"+"=?", value)
	if len(selfID) > 0 && selfID[0] > 0 {
		d.Where("`ID`!=?", selfID[0])
	}
	type Tmp struct {
		ID                    int `xorm:"INTEGER"`
		CreatedAt time.Time `xorm:"DATETIME created"`
	}
	has, err := d.Exist(new(Tmp))
	if err != nil {
		panic(err)
	}
	return !has
}

//获取表信息
func GetTableInfo(dao *xorm.Engine, tableName string) (*core.Table, error) {
	tables, err := dao.Dialect().GetTables()
	if err != nil {
		return nil, err
	}
	var table *core.Table
	for _, v := range tables {
		if strings.EqualFold(v.Name, tableName) {
			table = v
		}
	}
	if table == nil {
		return nil, fmt.Errorf("Unknown table: %s", tableName)
	}
	colSeq, cols, err := dao.Dialect().GetColumns(table.Name)
	if err != nil {
		return nil, err
	}
	for _, name := range colSeq {
		table.AddColumn(cols[name])
	}
	indexes, err := dao.Dialect().GetIndexes(table.Name)
	if err != nil {
		return nil, err
	}
	table.Indexes = indexes
	for _, index := range indexes {
		for _, name := range index.Cols {
			if col := table.GetColumn(name); col != nil {
				col.Indexes[index.Name] = index.Type
			} else {
				return nil, fmt.Errorf("Unknown col %s in index %v of table %v, columns %v", name, index.Name, table.Name, table.ColumnsSeq())
			}
		}
	}
	return table, nil
}
