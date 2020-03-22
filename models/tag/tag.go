package m_tag

import (
	"time"
)

//文章标签表
type Tag struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	CreatedAt time.Time `xorm:"DateTime created"`
	UpdatedAt time.Time `xorm:"DateTime updated"`
	Name         string `xorm:"unique TEXT"`
}

type List []Tag
