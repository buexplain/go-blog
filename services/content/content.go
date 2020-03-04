package s_content

import (
	"fmt"
	"github.com/88250/lute"
	h_boot "github.com/buexplain/go-blog/app/http/boot"
	"github.com/buexplain/go-blog/app/http/boot/code"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/content"
	"github.com/buexplain/go-blog/models/contentTag"
	s_category "github.com/buexplain/go-blog/services/category"
	"github.com/buexplain/go-fool/errors"
	"time"
)

//保存内容
func Save(content *m_content.Content, tagsID []int, id int) error {
	//将内容转为html
	if _, err := Render(content.Body); err != nil {
		return err
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
		if _, err := session.Where("ContentID=?", content.ID).Delete(new(m_contentTag.ContentTag)); err != nil {
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
	if affected, err := dao.Dao.In("ID", ids).
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
		err := dao.Dao.In("ID", ids).Find(&lists)
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
	Total int
	CreatedAtYm string
}
type PlaceList []*Place

//获取归档信息
func GetPlace() PlaceList {
	mod := dao.Dao.Table("Content").
		Select("COUNT(*) as total, strftime('%Y年%m月', CreatedAt) as CreatedAtYm").
		Where("Online=?", m_content.OnlineYes).
		GroupBy("CreatedAtYm").OrderBy("CreatedAt DESC")
	result := make(PlaceList, 0)
	err := mod.Find(&result)
	if err != nil {
		panic(err)
	}
	return result
}

//获取列表
func GetList(page int, limit int, categoryID int, keyword string, tagID int, place string) (result m_content.List) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 15
	}
	mod := dao.Dao.Table("Content").Desc("`Content`.`ID`")
	mod.Select("`Content`.`ID`, `Content`.`Title`, `Content`.`CreatedAt`, `Content`.`Hits`, `Content`.`Origin`")
	mod.Where("Content.`Online`=?", m_content.OnlineYes)
	//设置分页
	offset := (page - 1) * limit
	mod.Limit(limit, offset)
	//查询分类
	if categoryID > 0 {
		if tmp := s_category.GetSons(categoryID); tmp == nil {
			mod.Where("`Content`.`categoryID`=?", categoryID)
		}else {
			if len(tmp) == 1 {
				mod.Where("`Content`.`categoryID`=?", tmp[0].ID)
			}else {
				in := []int{}
				for _, v:= range tmp {
					in = append(in, v.ID)
				}
				mod.In("`Content`.`categoryID`", in)
			}
		}
	}
	//查询归档时间内的列表，place的格式：2006年01月
	if place != "" {
		if t, err := time.ParseInLocation("2006年01月", place, time.Local); err == nil {
			mod.Where("`Content`.`CreatedAt`>=?", t.String()).Where("`Content`.`CreatedAt`<?", t.AddDate(0, 1, 1).String())
		}else {
			h_boot.Logger.InfoF("parse time error: %s", err)
		}
	}
	//查询关键字
	if keyword != "" {
		mod.Where("`Content`.`Title` LIKE ?", fmt.Sprintf("%s%s%s", "%", keyword, "%"))
	}
	//标签查询
	if tagID > 0 {
		mod.Join("left", "ContentTag", "`Content`.`ID`=`ContentTag`.`ContentID`")
		mod.Where("`ContentTag`.`TagID`=?", tagID)
	}
	if err := mod.Find(&result); err != nil {
		panic(err)
	}
	return
}
