package s_services


import (
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
	"time"
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

//根据id批量软删除数据
func DestroyBatch(tableName string, IDs []int, callback ...dao.CallBack) (int64, error) {
	deletedAt := map[string]interface{}{"DeletedAt":time.Now()}
	mod := dao.Dao.Table(tableName).In("ID", IDs).Where("`DeletedAt` IS NULL OR `DeletedAt`=?", "0001-01-01 00:00:00")
	if callback != nil && len(callback) > 0 {
		mod = callback[0](mod)
	}
	return mod.Update(deletedAt)
}