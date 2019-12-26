package s_content

import (
	"fmt"
	"github.com/88250/lute"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/models/contentTag"
	"html/template"
)

//保存内容
func Save(content *m_content.Content, tagsID []int, id int) error {
	//将内容转为html
	if html, err := Render(content.Body); err != nil {
		return err
	}else {
		content.HTML = template.HTML(html)
	}

	session := dao.Dao.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	//编辑
	if id > 0 {
		content.ID = id
		//删除原有标签
		if _, err := session.Unscoped().Where("ContentID=?", content.ID).Delete(new(m_contentTag.ContentTag)); err != nil {
			if err := session.Rollback(); err != nil {
				return err
			}
			return err
		}
		//更新内容
		if _, err := session.ID(content.ID).AllCols().Update(content); err != nil {
			if err := session.Rollback(); err != nil {
				return err
			}
			return err
		}
	}else {
		//插入内容
		if _, err := session.Insert(content); err != nil {
			if err := session.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	//插入标签
	if len(tagsID) > 0 {
		contentTag := make([]*m_contentTag.ContentTag, 0, len(tagsID))
		for _, v := range tagsID {
			contentTag = append(contentTag, &m_contentTag.ContentTag{ContentID: content.ID, TagID: v})
		}
		if _, err := session.Insert(contentTag); err != nil {
			if err := session.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	if err := session.Commit(); err != nil {
		return err
	}

	return nil
}

//删除
func DestroyBatch(ids []int) error {
	if affected, err:= dao.Dao.
		Unscoped().
		In("ID", ids).
		Select("ID").
		Where("Online=?", m_content.OnlineNo).
		Delete(new(m_content.Content)); err != nil {
		return err
	} else if affected > 0 {
		if len(ids) == 1 {
			//只有一条数据时，直接删除tag
			_, _ = dao.Dao.Where("ContentID", ids[0]).Delete(&m_contentTag.ContentTag{})
			return nil
		}
		//检查还有哪些id是存在的
		lists := make(m_content.List, 0)
		err := dao.Dao.Unscoped().In("ID", ids).Find(&lists)
		if err != nil {
			return err
		}
		//筛选出不存在的id，删除它们的tag
		var ok bool
		new_ids := []int{}
		for _, id := range ids {
			ok = false
			for _, l := range lists {
				if id == l.ID {
					ok = true
					break
				}
			}
			if !ok {
				new_ids = append(new_ids, id)
			}
		}
		if len(new_ids) > 0 {
			_, _ = dao.Dao.In("ContentID", new_ids).Delete(&m_contentTag.ContentTag{})
		}
	}
	return nil
}

func RenderByID(id int) (string, error) {
	result := m_content.Content{}
	ok, err := dao.Dao.Where("ID=?", id).Get(&result)
	if err != nil {
		return "", err
	}else if !ok {
		return "", fmt.Errorf("%s", code.Text(code.NOT_FOUND_DATA))
	}
	return Render(result.Body)
}

func Render(markdown string) (string, error) {
	luteEngine := lute.New()
	//注销掉高亮部分，让js去渲染
	luteEngine.CodeSyntaxHighlight = false
	luteEngine.CodeSyntaxHighlightLineNum = false
	html, err := luteEngine.MarkdownStr("default", markdown)
	return html, err
}