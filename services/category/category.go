package s_category

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/category"
	"xorm.io/builder"
)

func Destroy(ids []int) (affected int64, err error) {
	//筛选没有文章数的分类
	notInContent := builder.Select("CategoryID").From("Content").Where(builder.In("CategoryID", ids)).GroupBy("CategoryID")
	//筛选有子分类的分类
	notIn := builder.Select("Pid").From("Category").Where(builder.In("Pid", ids))
	var sql string
	sql, err = builder.Delete().From("Category").
		Where(builder.In("ID", ids)).
		And(builder.NotIn("ID", notInContent)).
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

type Category struct {
	m_category.Category `xorm:"extends"`
	ContentNum int
}

type List []*Category

func (this List) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func GetALL() (List, error) {
	result := make(List, 0)
	err := dao.Dao.
		Table("Category").
		Select("Category.*, Content.ContentNum").
		Join("LEFT", "(SELECT count(*) as ContentNum, CategoryID FROM Content GROUP BY CategoryID) as Content", "Category.ID = Content.CategoryID").
		Desc("SortID").Find(&result)
	return result, err
}
