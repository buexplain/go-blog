package m_category

import (
	"encoding/json"
	"fmt"
	"github.com/buexplain/go-blog/dao"
	"github.com/buexplain/go-blog/models"
)

//文章分类表
type Category struct {
	models.Field `xorm:"extends"`
	//父id
	Pid int `xorm:"INTEGER"`
	//分类名
	Name string `xorm:"TEXT"`
	//排序id
	SortID int `xorm:"INTEGER"`
}

//返回当前分类的父级分类
func (this Category) Parent() (*Category, error) {
	if this.Pid == 0 {
		return nil, nil
	}
	tmp := &Category{}
	if b, err := dao.Dao.ID(this.Pid).Get(tmp); err != nil {
		return nil, err
	} else if !b {
		return nil, fmt.Errorf("not found parent category %d", this.Pid)
	}
	return tmp, nil
}

//返回当前分类的所有父级别分类
func (this Category) Parents() (List, error) {
	l := make(List, 0)
	tmp := &this
	for {
		if c, err := tmp.Parent(); err != nil {
			return nil, err
		} else if c == nil {
			break
		} else {
			l = append(l, c)
			tmp = c
		}
	}
	return l, nil
}

type List []*Category

func (this List) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}