package m_category

func (this Category) IsMenuText() string {
	return IsMenuText[this.IsMenu]
}

func CheckIsMenu(isNav int) bool {
	if isNav >= IsMenuYes && isNav <= IsMenuNo {
		return true
	}
	return false
}

const (
	IsMenuYes = iota + 1
	IsMenuNo
)

var IsMenuText = map[int]string{
	IsMenuYes: "是",
	IsMenuNo:  "否",
}
