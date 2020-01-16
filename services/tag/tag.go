package s_tag

import (
	"github.com/buexplain/go-blog/dao"
	m_tag "github.com/buexplain/go-blog/models/tag"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-fool"
	"xorm.io/builder"
)

type Tag struct {
	m_tag.Tag `xorm:"extends"`
	ContentNum int
}

type List []*Tag

func GetList(ctx *fool.Ctx)(counter int64, result List, err error)  {
	query := s_services.NewQuery("Tag", ctx).Limit()
	query.Finder.Desc("ID")
	query.Where()
	//先获取分页所需的总条数
	counter = query.Count()
	//再连表查询得到每个标签的文章数量
	query.Finder.Join(
		"LEFT",
		"(SELECT count(*) as ContentNum, TagID FROM ContentTag GROUP BY TagID) as ContentTag",
		"Tag.ID = ContentTag.TagID",
		)
	query.Find(&result)
	err = query.Error
	return
}

func Destroy(ids []int) (int64, error) {
	notIn := builder.Select("TagID").From("ContentTag").Where(builder.In("TagID", ids)).GroupBy("TagID")
	sql, err := builder.Delete().From("Tag").
		Where(builder.In("ID", ids)).
		And(builder.NotIn("ID", notIn)).
		ToBoundSQL()
	if err != nil {
		return 0, err
	}
	if result, err := dao.Dao.Exec(sql); err != nil {
		return 0, err
	}else {
		return result.RowsAffected()
	}
}

func Store(name string) (int, error) {
	mod := &m_tag.Tag{}
	mod.Name = name
	has, err := dao.Dao.Get(mod)
	if err != nil {
		return 0, err
	}
	if has {
		return mod.ID, nil
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return 0, err
	}

	return mod.ID, nil
}
