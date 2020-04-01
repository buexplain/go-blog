package s_content

import (
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/category"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/models/tag"
	s_category "github.com/buexplain/go-blog/services/category"
	"github.com/buexplain/go-fool/errors"
	"html/template"
)

type Details struct {
	Content  *m_content.Content `xorm:"extends"`
	Category m_category.List
	Tag      *m_tag.List `xorm:"extends"`
}

func GetDetails(id int, online m_content.Online) (*Details, error) {
	details := new(Details)
	details.Content = new(m_content.Content)
	details.Category = m_category.List{}
	details.Tag = new(m_tag.List)
	mod := dao.Dao.Table("Content").ID(id)
	if m_content.CheckOnline(online) {
		mod.Where("Online=?", int(online));
	}
	if has, err := mod.Get(details.Content); err != nil {
		return nil, err
	} else if !has {
		return nil, errors.MarkClient(errors.New(code.Text(code.NOT_FOUND_DATA, id)))
	}

	//渲染html成
	details.Content.HTML = template.HTML(LuteEngine.Md2HTML(details.Content.Body))

	err := dao.Dao.
		Table("Tag").
		Join("INNER", "ContentTag", "Tag.ID = ContentTag.TagID").
		Where("ContentTag.ContentID=?", details.Content.ID).
		Find(details.Tag)
	if err != nil {
		return nil, err
	}

	details.Category = s_category.GetParents(details.Content.CategoryID)

	return details, nil
}
