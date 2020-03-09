package m_configGroup

import "github.com/buexplain/go-blog/models"

//配置组表
type ConfigGroup struct {
	models.Field `xorm:"extends"`
	//配置项名称
	Name string `xorm:"TEXT"`
	//字段名
	Key string `xorm:"TEXT"`
	//备注
	Comment string `xorm:"TEXT"`
}

type List []*ConfigGroup
