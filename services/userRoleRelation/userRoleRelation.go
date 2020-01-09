package s_userRoleRelation

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	m_role "github.com/buexplain/go-blog/models/role"
	m_userRoleRelation "github.com/buexplain/go-blog/models/userRoleRelation"
)

type UserRole struct {
	m_role.Role `xorm:"extends"`
	Checked     bool `xorm:"-"`
}

type UserRoleList []*UserRole

func (this UserRoleList) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

//获取用户角色
func GetUserRole(userID int) (UserRoleList, error) {
	//获取所有的角色
	allRole := make(UserRoleList, 0)
	err := dao.Dao.Table("Role").Desc("SortID").Find(&allRole)
	if err != nil {
		return nil, err
	}

	if userID > 0 {
		//获取用户拥有的角色
		userRole := make(m_userRoleRelation.List, 0)
		err = dao.Dao.Table("UserRoleRelation").Where("UserID=?", userID).Find(&userRole)
		if err != nil {
			return nil, err
		}
		if len(userRole) > 0 {
			for _, role := range allRole {
				role.Checked = userRole.HasRoleID(role.ID)
			}
		}
	}

	return allRole, nil
}

//设置用户角色
func SetUserRole(userID int, roleID []int) error {
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
		result = append(result, m_userRoleRelation.UserRoleRelation{UserID: userID, RoleID: v})
	}
	_, err = session.Insert(&result)
	if err != nil {
		if err := session.Rollback(); err != nil {
			return err
		}
		return err
	} else {
		if err := session.Commit(); err != nil {
			return err
		}
		return nil
	}
}

//用户角色id
type UserRoleID int

type UserRoleIDList []UserRoleID

//获取用户的角色id列表
func GetRoleIDByUserID(userID int) (UserRoleIDList, error) {
	result := make(UserRoleIDList, 0)
	session := dao.Dao.Table("`User`").Where("`User`.`ID`=?", userID)
	session.Join("INNER", "`UserRoleRelation`", "`UserRoleRelation`.UserID = `User`.ID")
	session.Select("`UserRoleRelation`.`RoleID` as `UserRoleID`")
	err := session.Find(&result)
	return result, err
}
