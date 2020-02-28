package m_contentTag

import "github.com/buexplain/go-blog/models"

type ContentTag struct {
	models.Field `xorm:"extends"`
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
