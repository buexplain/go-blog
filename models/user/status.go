package m_user

import m_models "github.com/buexplain/go-blog/models"

//用户状态，用于判断用户是否被ban
type Status m_models.Enum

const (
	StatusAllow Status = iota + 1
	StatusDeny
)

func (this Status) String() string {
	switch this {
	case StatusAllow:
		return "允许"
	case StatusDeny:
		return "禁止"
	default:
		return "UNKNOWN"
	}
}
