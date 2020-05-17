package m_node

import m_models "github.com/buexplain/go-blog/models"

// 用于判断后台是否显示
type IsMenu m_models.Enum

const (
	IsMenuYes IsMenu = iota + 1
	IsMenuNo
)

func (this IsMenu) String() string {
	switch this {
	case IsMenuYes:
		return "是"
	case IsMenuNo:
		return "否"
	default:
		return m_models.EnumUNKNOWN
	}
}
