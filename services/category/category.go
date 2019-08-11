package s_category

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	m_category "github.com/buexplain/go-blog/models/category"
)

func Destroy(ids[]int) (affected int64, err error) {
	childList := make(m_category.List, 0)
	if err := dao.Dao.In("Pid", ids).Find(&childList); err != nil {
		return 0, err
	}else if len(childList) > 0 {
		//ids中的id含有子级分类，判断这些子级分类是否也在ids中
		for _, child := range childList {
			has := false
			for _, id := range ids {
				if child.ID == id {
					has = true
					break
				}
			}
			if !has {
				return 0, fmt.Errorf("入参错误：ID【%d】必须与其父ID【%d】一并删除", child.ID, child.Pid)
			}
		}
	}
	affected, err = dao.Dao.In("ID", ids).Delete(new(m_category.Category))
	return
}