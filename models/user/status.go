package m_user

const (
	StatusAllow = iota + 1
	StatusDeny
)

var StatusText = map[int]string{
	StatusAllow: "允许",
	StatusDeny:  "禁止",
}

func CheckStatus(status int) bool {
	if status != StatusAllow && status != StatusDeny {
		return false
	}
	return true
}

