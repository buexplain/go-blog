package m_category

import (
	"encoding/json"
	"github.com/buexplain/go-blog/models"
)

//文章分类表
type Category struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	CreatedAt models.Time `xorm:"DateTime created"`
	UpdatedAt models.Time `xorm:"DateTime updated"`
	//父id
	Pid int `xorm:"INTEGER"`
	//分类名
	Name string `xorm:"TEXT"`
	//跳转地址
	Redirect string `xorm:"TEXT"`
	//是否为前台导航
	IsMenu int `xorm:"INTEGER"`
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
