package m_category

import m_models "github.com/buexplain/go-blog/models"

// 用于判断前台是否显示
type IsMenu m_models.Enum

func CheckIsMenu(isMenu IsMenu) bool {
	if isMenu >= IsMenuYes && isMenu <= IsMenuNo {
		return true
	}
	return false
}

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
		return "UNKNOWN"
	}
}
