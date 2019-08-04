package s_menu

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	m_menu "github.com/buexplain/go-blog/models/menu"
)

func Destroy(ids[]int) (affected int64, err error) {
	childList := make(m_menu.List, 0)
	if err := dao.Dao.In("Pid", ids).Find(&childList); err != nil {
		return 0, err
	}else if len(childList) > 0 {
		//ids中的id含有子级菜单，判断这些子级菜单是否也在ids中
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
	affected, err = dao.Dao.In("ID", ids).Delete(new(m_menu.Menu))
	return
}