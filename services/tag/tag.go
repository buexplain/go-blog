package s_tag

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	m_content "github.com/buexplain/go-blog/models/content"
	m_tag "github.com/buexplain/go-blog/models/tag"
	s_services "github.com/buexplain/go-blog/services"
	"github.com/buexplain/go-slim"
	"strings"
	"time"
	"xorm.io/builder"
)

type Tag struct {
	m_tag.Tag `xorm:"extends"`
	Total     int
}

type List []*Tag

func GetList(ctx *slim.Ctx) (counter int64, result List, err error) {
	query := s_services.NewQuery("Tag", ctx)
	//先获取分页所需的总条数
	counter = query.Where().Count()
	//再连表查询得到每个标签的文章数量
	query.Where().Limit().Finder.Desc("Tag.ID")
	query.Finder.Join(
		"LEFT",
		"(SELECT count(*) as Total, TagID FROM ContentTag GROUP BY TagID) as ContentTag",
		"Tag.ID = ContentTag.TagID",
	)
	query.Find(&result)
	err = query.Error
	return
}

func GetALL() (result List) {
	countSql := fmt.Sprintf("(SELECT count(*) as Total, TagID FROM ContentTag INNER JOIN Content ON `ContentTag`.`ContentID`=`Content`.`ID` WHERE `Content`.`Online`=%d GROUP BY TagID) as ContentTag", m_content.OnlineYes)
	mod := dao.Dao.Table("Tag").Desc("`Total`").Join(
		"LEFT",
		countSql,
		"Tag.ID = ContentTag.TagID",
	)
	err := mod.Find(&result)
	if err != nil {
		panic(err)
	}
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
	} else {
		return result.RowsAffected()
	}
}

func Store(name string) (*m_tag.Tag, error) {
	mod := &m_tag.Tag{}
	mod.Name = name
	has, err := dao.Dao.Get(mod)
	if err != nil {
		return nil, err
	}
	if has {
		return mod, nil
	}

	if _, err := dao.Dao.Insert(mod); err != nil {
		return nil, err
	}

	return mod, nil
}

func Stores(names []string) (int64, error) {
	t := time.Now().UTC().Format("2006-01-02 15:04:05")
	args := make([]interface{}, 0, len(names)*3)
	for _, name := range names {
		args = append(args, t, t, name)
	}
	values, err := builder.ConvertToBoundSQL(strings.Repeat("(?,?,?),", len(names)), args)
	if err != nil {
		return 0, err
	}
	sql := "INSERT OR IGNORE INTO `Tag` (`CreatedAt`,`UpdatedAt`,`Name`) VALUES " + values[0:len(values)-1]
	result, err := dao.Dao.Exec(sql)
	if err != nil {
		return 0, err
	}
	if affected, err := result.RowsAffected(); err != nil {
		return 0, err
	} else {
		return affected, nil
	}
}
