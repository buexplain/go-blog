package s_content

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/category"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/models/tag"
)

type Details struct {
	Content    *m_content.Content    `xorm:"extends"`
	Category   m_category.List
	Tag        *m_tag.List  `xorm:"extends"`
}

func GetDetails(id int) (*Details, error) {
	details := new(Details)
	details.Content = new(m_content.Content)
	details.Category = m_category.List{}
	details.Tag = new(m_tag.List)
	if b, err := dao.Dao.Table("Content").ID(id).Get(details.Content); err != nil {
		return nil, err
	}else if(b == false) {
		return nil, fmt.Errorf("not found content %d", id)
	}

	err := dao.Dao.
		Table("Tag").
		Join("INNER", "ContentTag", "Tag.ID = ContentTag.TagID").
		Where("ContentTag.ContentID=?", details.Content.ID).
		Find(details.Tag)
	if err != nil {
		return nil, err
	}

	tmp := new(m_category.Category)
	if b, err := dao.Dao.ID(details.Content.Category).Get(tmp); err == nil && b {
		details.Category = append(details.Category, tmp)
		if p, err := tmp.Parents(); err == nil {
			details.Category = append(details.Category, p...)
		}
	}

	return details, nil
}