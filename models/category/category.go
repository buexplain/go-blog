package m_category

import (
	"encoding/json"
	"github.com/buexplain/go-blog/models"
)

//文章分类表
type Category struct {
	models.Field `xorm:"extends"`
	//父id
	Pid int `xorm:"INTEGER"`
	//分类名
	Name string `xorm:"TEXT"`
	//跳转地址
	Redirect string `xorm:"TEXT"`
	//排序id
	SortID int `xorm:"INTEGER"`
}

type List []*Category

func (this List) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}