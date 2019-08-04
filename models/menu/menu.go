package m_menu

import (
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
)

type Menu struct {
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

type List []Menu

func GetALL() ([]Menu, error) {
	result := make([]Menu, 0)
	err := dao.Dao.Table("Menu").Desc("SortID").Find(&result);
	return result, err
}

