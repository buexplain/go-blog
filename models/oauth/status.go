package m_oauth

import m_models "github.com/buexplain/go-blog/models"

//第三方类型
type Status m_models.Enum

const (
	StatusGithub Status = iota + 1
)

func (this Status) String() string {
	switch this {
	case StatusGithub:
		return "Github"
	default:
		return m_models.EnumUNKNOWN
	}
}
