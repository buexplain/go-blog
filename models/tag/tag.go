package m_tag

import (
	"github.com/buexplain/go-blog/models"
)

//文章标签表
type Tag struct {
	ID        int           `xorm:"not null pk autoincr INTEGER"`
	CreatedAt m_models.Time `xorm:"DateTime created"`
	UpdatedAt m_models.Time `xorm:"DateTime updated"`
	Name      string        `xorm:"unique TEXT"`
}

type List []Tag
