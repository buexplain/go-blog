package s_category

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models/category"
	s_services "github.com/buexplain/go-blog/services"
	"strings"
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
	} else {
		return result.RowsAffected()
	}
}

type Category struct {
	m_category.Category `xorm:"extends"`
	Total               int
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
func GetALL() List {
	result := make(List, 0)
	err := dao.Dao.
		Table("Category").
		Select("Category.*, Content.Total").
		Join("LEFT", "(SELECT count(*) as Total, CategoryID FROM Content GROUP BY CategoryID) as Content", "Category.ID = Content.CategoryID").
		Asc("SortID").Find(&result)
	if err != nil {
		panic(err)
	}
	return result
}

type TreeItem struct {
	m_category.Category `xorm:"extends"`
	Children            []*TreeItem `xorm:"-"`
}

type TreeList []*TreeItem

//获取所有分类，返回树状结构
func GetTree() TreeList {
	lists := make(TreeList, 0)
	err := dao.Dao.Table("Category").Where("IsMenu=?", m_category.IsMenuYes).Asc("SortID").Find(&lists)
	if err != nil {
		panic(err)
	}
	m := make(map[int]*TreeItem)
	for _, v := range lists {
		m[v.ID] = v
	}
	result := make(TreeList, 0)
	for _, v := range lists {
		if v.Pid == 0 {
			result = append(result, v)
		}
		i, ok := m[v.Pid]
		if !ok {
			continue
		}
		i.Children = append(i.Children, v)
	}
	return result
}

//获取一个分类的父分类
func GetParents(categoryID int) m_category.List {
	if categoryID <= 0 {
		return nil
	}
	list := make(m_category.List, 0)
	err := s_services.GetParents("category", categoryID, &list)
	if err != nil {
		return nil
	}
	return list
}

//获取一个分类的子分类
func GetSons(categoryID int, isMenu int) m_category.List {
	if categoryID <= 0 {
		return nil
	}
	list := make(m_category.List, 0)
	var err error
	if m_category.CheckIsMenu(isMenu) {
		err = s_services.GetSons("category", categoryID, &list, "SortID", builder.Eq{"IsMenu": isMenu})
	} else {
		err = s_services.GetSons("category", categoryID, &list, "SortID", nil)
	}
	if err != nil {
		return nil
	}
	return list
}

func Store(mod *m_category.Category) (affected int64, err error) {
	if mod.IsMenu == 0 {
		mod.IsMenu = m_category.IsMenuNo
	} else {
		mod.IsMenu = m_category.IsMenuYes
	}
	mod.Redirect = strings.Trim(mod.Redirect, " ")
	if mod.ID == 0 {
		return dao.Dao.MustCols("Pid", "Redirect").Insert(mod)
	}
	return dao.Dao.ID(mod.ID).MustCols("Pid", "Redirect").Update(mod)
}
