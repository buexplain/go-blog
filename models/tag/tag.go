package m_tag

import "github.com/buexplain/go-blog/models"

type Tag struct {
	models.Field `xorm:"extends"`
	Name         string `xorm:"index TEXT"`
}

type List []Tag 