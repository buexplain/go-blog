package s_content

import (
	"fmt"
	"github.com/88250/lute"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/models/contentTag"
	"github.com/buexplain/go-fool/errors"
	"html/template"
	"time"
)

//保存内容
func Save(content *m_content.Content, tagsID []int, id int) error {
	//将内容转为html
	if html, err := Render(content.Body); err != nil {
		return err
	} else {
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
	} else {
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
func Destroy(ids []int) error {
	if affected, err := dao.Dao.
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
	has, err := dao.Dao.Where("ID=?", id).Get(&result)
	if err != nil {
		return "", err
	} else if !has {
		return "", errors.MarkClient(errors.New(code.Text(code.NOT_FOUND_DATA, id)))
	}
	return Render(result.Body)
}

func Render(markdown string) (string, error) {
	luteEngine := lute.New()
	//注销掉高亮部分，让js去渲染
	luteEngine.CodeSyntaxHighlight = false
	luteEngine.CodeSyntaxHighlightLineNum = false
	html, err := luteEngine.MarkdownStr("default", markdown)
	return html, errors.TryMarkClient(err)
}

type Place struct {
	Counter int
	CreatedAtYm string
}
type PlaceList []*Place

//获取归档信息
func GetPlace() (PlaceList, error) {
	mod := dao.Dao.Table("Content").
		Select("COUNT(*) as counter, strftime('%Y年%m月', CreatedAt) as CreatedAtYm").
		GroupBy("CreatedAtYm").OrderBy("CreatedAt DESC")
	result := make(PlaceList, 0)
	err := mod.Find(&result)
	return result, err
}

//获取列表
func GetList(page int, limit int, categoryID int, place string, keyword string) (result m_content.List, counter int64, err error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 15
	}
	mod := dao.Dao.Table("Content").Desc("ID")
	offset := (page - 1) * limit
	//设置分页
	mod.Limit(limit, offset)
	if categoryID > 0 {
		//查询分类
		mod.Where("categoryID=?", categoryID)
	}
	if place != "" {
		//查询归档时间内的列表，place的格式：2006-01
		if t, err := time.ParseInLocation("2006-01", place, time.Local); err == nil {
			mod.Where("CreatedAt>=?", t.String()).Where("CreatedAt<?", t.AddDate(0, 1, 1).String())
		}
	}
	if keyword != "" {
		//查询关键字
		mod.Where("Title LIKE ?", fmt.Sprintf("%s%s%s", "%", keyword, "%"))
	}
	counter, err = mod.FindAndCount(&result)
	return
}
