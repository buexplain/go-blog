package s_category

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/category"
	s_services "github.com/buexplain/go-blog/services"
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

//返回所有的分类，返回平面结构
func GetALL() (List, error) {
	result := make(List, 0)
	err := dao.Dao.
		Table("Category").
		Select("Category.*, Content.ContentNum").
		Join("LEFT", "(SELECT count(*) as ContentNum, CategoryID FROM Content GROUP BY CategoryID) as Content", "Category.ID = Content.CategoryID").
		Asc("SortID").Find(&result)
	return result, err
}

type TreeItem struct {
	m_category.Category `xorm:"extends"`
	Children []*TreeItem `xorm:"-"`
}

type TreeList []*TreeItem

//获取所有分类，返回树状结构
func GetTree() (TreeList, error) {
	lists := make(TreeList, 0)
	err := dao.Dao.Table("Category").Desc("SortID").Find(&lists)
	if err != nil {
		return nil, err
	}
	m := make(map[int]*TreeItem)
	for _,v := range lists {
		m[v.ID] = v
	}
	result := make(TreeList, 0)
	for _,v := range lists {
		if v.Pid == 0 {
			result = append(result, v)
		}
		i, ok := m[v.Pid]
		if !ok {
			continue
		}
		i.Children = append(i.Children, v)
	}
	return result, nil
}

//获取一个分类的父分类
func GetParents(categoryID int) (m_category.List) {
	if categoryID <= 0 {
		return nil
	}
	list := make(m_category.List, 0)
	err := s_services.GetRecursion("category", categoryID, &list, nil)
	if err != nil {
		return nil
	}
	return list
}