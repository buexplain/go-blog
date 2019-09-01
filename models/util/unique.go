package m_util

import (
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
)

//检查字段值是否存在
func CheckUnique(tableName string, field string, value interface{}, selfID ...int) bool {
	d := dao.Dao.Table("`" + tableName + "`")
	d.Where("`"+field+"`"+"=?", value)
	if len(selfID) > 0 && selfID[0] > 0 {
		d.Where("`ID`!=?", selfID[0])
	}
	type Tmp struct {
		ID int `xorm:"INTEGER"`
		models.DeletedAtField `xorm:"extends"`
	}
	has, err := d.Exist(new(Tmp))
	if err != nil {
		panic(err)
	}
	return !has
}