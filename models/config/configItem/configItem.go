package m_configItem

import "github.com/buexplain/go-blog/models"

//配置组表
type ConfigItem struct {
	models.IDField `xorm:"extends"`
	//配置组的字段
	Group string `xorm:"TEXT"`
	//配置项名称
	Name string `xorm:"TEXT"`
	//字段名
	Key string `xorm:"TEXT"`
	//配置的值
	Value string `xorm:"TEXT"`
	//备注
	Comment string `xorm:"TEXT"`
}
