package m_role

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
)

//rbac 角色表
type Role struct {
	models.Field `xorm:"extends"`
	//父id
	Pid          int    `xorm:"INTEGER"`
	//角色名
	Name         string `xorm:"TEXT"`
	//排序id
	SortID       int    `xorm:"INTEGER"`
}

type List []Role

func (this List) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func GetALL() (List, error) {
	result := make(List, 0)
	err := dao.Dao.Table("Role").Desc("SortID").Find(&result)
	return result, err
}
