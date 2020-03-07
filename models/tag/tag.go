package m_tag

import "github.com/buexplain/go-blog/models"

//文章标签表
type Tag struct {
	models.Field `xorm:"extends"`
	Name         string `xorm:"unique TEXT"`
}

type List []Tag
