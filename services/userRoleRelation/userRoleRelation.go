package s_userRoleRelation

import (
	"github.com/buexplain/go-blog/dao"
	m_userRoleRelation "github.com/buexplain/go-blog/models/userRoleRelation"
)

func Store(userID int, roleID []int) error {
	//开启事务
	session := dao.Dao.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return err
	}

	//先删除已有关系
	r := &m_userRoleRelation.UserRoleRelation{}
	r.UserID = userID
	_, err := session.Delete(r)
	if err != nil {
		if err := session.Rollback(); err != nil {
			return err
		}
		return err
	}

	//插入新关系
	result := make(m_userRoleRelation.List, 0, len(roleID))
	for _, v := range roleID {
		result = append(result, m_userRoleRelation.UserRoleRelation{UserID:userID, RoleID:v})
	}
	_, err = session.Insert(&result)
	if err != nil {
		if err := session.Rollback(); err != nil {
			return err
		}
		return err
	}else {
		if err := session.Commit(); err != nil {
			return err
		}
		return nil
	}
}
