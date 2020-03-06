package m_category

func (this Category) IsMenuText() string {
	return IsMenuText[this.IsMenu]
}

func CheckIsMenu(IsMenu int) bool {
	if IsMenu >= IsMenuYes && IsMenu <= IsMenuNo {
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
