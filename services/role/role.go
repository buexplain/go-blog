package s_role

import (
	"fmt"
	"github.com/buexplain/go-blog/dao"
	m_role "github.com/buexplain/go-blog/models/role"
	"github.com/buexplain/go-blog/models/roleNodeRelation"
	"github.com/buexplain/go-blog/models/userRoleRelation"
)

func Destroy(ids[]int) (int64, error) {
	childList := make(m_role.List, 0)
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

	//开启事务
	session := dao.Dao.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, err
	}

	//删除角色
	affected, deleteErr := session.In("ID", ids).Delete(new(m_role.Role))
	if deleteErr != nil {
		if err := session.Rollback(); err != nil {
			return 0, err
		}
		return affected, deleteErr
	}

	//删除角色节点表
	if _, err := session.In("RoleID", ids).Delete(new(m_roleNodeRelation.RoleNodeRelation)); err != nil {
		if err := session.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	//删除用户角色表
	if _, err := session.In("RoleID", ids).Delete(new(m_userRoleRelation.UserRoleRelation)); err != nil {
		if err := session.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	if err := session.Commit(); err != nil {
		return 0, err
	}

	return affected, deleteErr
}