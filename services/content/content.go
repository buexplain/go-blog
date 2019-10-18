package s_content

import (
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/models/contentTag"
	"github.com/buexplain/go-blog/services/tag"
)

//保存内容
func Save(content *m_content.Content, tagsID []int, id int, tagsName[]string) error {
	//先处理标签
	for _, v := range tagsName {
		if id, err := s_tag.Store(v); err != nil {
			return err
		}else {
			tagsID = append(tagsID, id)
		}
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