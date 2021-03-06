package m_userRoleRelation

//rbac 用户与角色的关系表
type UserRoleRelation struct {
	ID int `xorm:"not null pk autoincr INTEGER"`
	//用户id
	UserID int `xorm:"INTEGER"`
	//角色id
	RoleID int `xorm:"INTEGER"`
}

type List []UserRoleRelation

func (this List) HasRoleID(roleID int) bool {
	for _, v := range this {
		if v.RoleID == roleID {
			return true
		}
	}
	return false
}
