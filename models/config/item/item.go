package m_configItem

import (
	"time"
)

//配置组表
type ConfigItem struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	CreatedAt time.Time `xorm:"DateTime created"`
	UpdatedAt time.Time `xorm:"DateTime updated"`
	GroupID      int `xorm:"INTEGER"`
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
