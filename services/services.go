package s_services

import (
	"github.com/buexplain/go-blog/dao"
	"strings"
	"time"
	"xorm.io/xorm"
	"xorm.io/xorm/schemas"
)

//检查字段值是否存在
func CheckUnique(tableName string, field string, value interface{}, selfID ...int) bool {
	d := dao.Dao.Table("`" + tableName + "`")
	d.Where("`"+field+"`"+"=?", value)
	if len(selfID) > 0 && selfID[0] > 0 {
		d.Where("`ID`!=?", selfID[0])
	}
	type Tmp struct {
		ID        int       `xorm:"INTEGER"`
		CreatedAt time.Time `xorm:"DATETIME created"`
	}
	has, err := d.Exist(new(Tmp))
	if err != nil {
		panic(err)
	}
	return !has
}

//获取表信息
func GetTableInfo(dao *xorm.Engine, tableName string) (*schemas.Table, error) {
	tables, err := dao.DBMetas()
	if err != nil {
		return nil, err
	}
	var table *schemas.Table
	for _, v := range tables {
		if strings.EqualFold(v.Name, tableName) {
			table = v
		}
	}
	return table, nil
}
