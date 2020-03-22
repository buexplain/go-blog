package m_contentTag

type ContentTag struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	ContentID      int `xorm:"index(ContentID_TagID) INTEGER"`
	TagID          int `xorm:"index(ContentID_TagID) INTEGER"`
}

type List []ContentTag

func (this List) HasTagID(TagID int) bool {
	for _, v := range this {
		if v.TagID == TagID {
			return true
		}
	}
	return false
}
