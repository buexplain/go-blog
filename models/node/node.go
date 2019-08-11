package m_node

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
)

type Node struct {
	models.Field `xorm:"extends"`
	//父id
	Pid          int    `xorm:"INTEGER"`
	//菜单名
	Name         string `xorm:"TEXT"`
	//菜单url
	URL          string `xorm:"TEXT"`
	//排序id
	SortID       int    `xorm:"INTEGER"`
}

type List []Node

func (this List) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func GetALL() (List, error) {
	result := make(List, 0)
	err := dao.Dao.Table("Node").Desc("SortID").Find(&result)
	return result, err
}

