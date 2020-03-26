package m_configItem

import (
	"github.com/buexplain/go-blog/models"
)

//配置组表
type ConfigItem struct {
	ID        int           `xorm:"not null pk autoincr INTEGER"`
	CreatedAt m_models.Time `xorm:"DateTime created"`
	UpdatedAt m_models.Time `xorm:"DateTime updated"`
	GroupID   int           `xorm:"INTEGER"`
	//配置项名称
	Name string `xorm:"TEXT"`
	//字段名
	Key string `xorm:"TEXT"`
	//配置的值
	Value string `xorm:"TEXT"`
	//备注
	Comment string `xorm:"TEXT"`
}

type List []*ConfigItem
