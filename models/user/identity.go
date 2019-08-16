package m_user

//用户身份，用于区分是否可用登录后台
const (
	IdentityOfficial = iota + 1
	IdentityCitizen
)

var IdentityText = map[int]string{
	IdentityOfficial:  "管理人员",
	IdentityCitizen: "普通用户",
}

func CheckIdentity(identity int) bool {
	if identity != IdentityOfficial && identity != IdentityCitizen {
		return false
	}
	return true
}

