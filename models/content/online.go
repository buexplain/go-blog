package m_content

import m_models "github.com/buexplain/go-blog/models"

// 用于判断前台是否显示
type Online m_models.Enum

func CheckOnline(online Online) bool {
	if online >= OnlineYes && online <= OnlineNo {
		return true
	}
	return false
}

const (
	OnlineYes Online = iota + 1
	OnlineNo
)

func (this Online) String() string {
	switch this {
	case OnlineYes:
		return "已上线"
	case OnlineNo:
		return "已下线"
	default:
		return m_models.EnumUNKNOWN
	}
}
