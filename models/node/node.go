package m_node

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
)

//rbac 节点表
type Node struct {
	models.Field `xorm:"extends"`
	//父id
	Pid          int    `xorm:"INTEGER"`
	//节点名称
	Name         string `xorm:"TEXT"`
	//节点路径
	URL          string `xorm:"TEXT"`
	//是否为后台导航菜单
	IsMenu       int    `xorm:"INTEGER"`
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

func GetMenu() (List, error) {
	result := make(List, 0)
	err := dao.Dao.Table("Node").Where("IsMenu=?", IsMenuYes).Desc("SortID").Find(&result)
	return result, err
}