package m_configGroup

import (
	"time"
)

//配置组表
type ConfigGroup struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	CreatedAt time.Time `xorm:"DateTime created"`
	UpdatedAt time.Time `xorm:"DateTime updated"`
	//配置项名称
	Name string `xorm:"TEXT"`
	//字段名
	Key string `xorm:"TEXT"`
	//备注
	Comment string `xorm:"TEXT"`
}

type List []*ConfigGroup
