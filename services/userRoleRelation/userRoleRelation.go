package s_userRoleRelation

import (
	"encoding/json"
	"github.com/buexplain/go-blog/dao"
	m_role "github.com/buexplain/go-blog/models/role"
	m_userRoleRelation "github.com/buexplain/go-blog/models/userRoleRelation"
	"strings"
)

type Relation struct {
	m_role.Role `xorm:"extends"`
	//是否拥有该角色
	Checked bool `xorm:"-"`
}

type RelationList []*Relation

func (this RelationList) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

//获取用户与所有角色的关系
func GetRelation(userID int) (RelationList, error) {
	//获取所有的角色
	allRole := make(RelationList, 0)
	err := dao.Dao.Table("Role").Desc("SortID").Find(&allRole)
	if err != nil {
		return nil, err
	}

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

	return allRole, nil
}

//设置用户角色
func SetRelation(userID int, roleID []int) error {
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
	session := dao.Dao.Table("`UserRoleRelation`").
		Where("`UserRoleRelation`.`UserID`=?", userID).
		Select("`UserRoleRelation`.`RoleID` as `UserRoleID`")
	err := session.Find(&result)
	return result, err
}

//用户角色id
type UserRoleName string

type UserRoleNameList []UserRoleName

func (this UserRoleNameList) String() string {
	b := strings.Builder{}
	for _, v := range this {
		b.WriteString(", ")
		b.WriteString(string(v))
	}
	s := b.String()
	if s != "" {
		s = s[len(", "):]
	}
	return s
}

//获取用户的角色名称列表
func GetRoleNameByUserID(userID int) (UserRoleNameList, error) {
	result := make(UserRoleNameList, 0)
	session := dao.Dao.Table("`UserRoleRelation`").Where("`UserRoleRelation`.`UserID`=?", userID)
	session.Join("INNER", "`Role`", "`UserRoleRelation`.RoleID = `Role`.ID")
	session.Select("`Role`.`Name` as `UserRoleName`")
	err := session.Find(&result)
	return result, err
}
