package m_user


const (
	IdentityOfficial = iota + 1
	IdentityCitizen
)

var IdentityText = map[int]string{
	IdentityOfficial:  "管理",
	IdentityCitizen: "用户",
}

func CheckIdentity(identity int) bool {
	if identity != IdentityOfficial && identity != IdentityCitizen {
		return false
	}
	return true
}

