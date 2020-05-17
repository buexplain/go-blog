package m_user

import m_models "github.com/buexplain/go-blog/models"

//用户身份，用于区分是否可以登录后台
type Identity m_models.Enum

const (
	IdentityOfficial Identity = iota + 1
	IdentityCitizen
)

func (this Identity) String() string {
	switch this {
	case IdentityOfficial:
		return "管理人员"
	case IdentityCitizen:
		return "普通用户"
	default:
		return m_models.EnumUNKNOWN
	}
}
