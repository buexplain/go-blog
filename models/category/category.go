package m_category

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
)

type Category struct {
	models.Field `xorm:"extends"`
	//父id
	Pid          int    `xorm:"INTEGER"`
	//分类名
	Name         string `xorm:"TEXT"`
	//分类url
	URL          string `xorm:"TEXT"`
	//排序id
	SortID       int    `xorm:"INTEGER"`
}

type List []Category

func (this List) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func GetALL() (List, error) {
	result := make(List, 0)
	err := dao.Dao.Table("Category").Desc("SortID").Find(&result)
	return result, err
}

